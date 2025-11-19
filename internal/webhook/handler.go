package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Handler handles GitHub webhook events
type Handler struct {
	secret string
	queue  QueuePublisher
}

// QueuePublisher defines the interface for publishing events to a queue
type QueuePublisher interface {
	Publish(event interface{}) error
}

// NewHandler creates a new webhook handler
func NewHandler(secret string, queue QueuePublisher) *Handler {
	return &Handler{
		secret: secret,
		queue:  queue,
	}
}

// HandleGitHub processes GitHub webhook events
func (h *Handler) HandleGitHub(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify signature
	signature := r.Header.Get("X-Hub-Signature-256")
	if !h.verifySignature(body, signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Get event type
	eventType := r.Header.Get("X-GitHub-Event")

	// We only care about pull request events
	if eventType != "pull_request" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Event type %s ignored", eventType)
		return
	}

	// Parse the event
	var event GitHubPullRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		// Log the error for debugging
		fmt.Printf("Failed to parse GitHub event: %v\n", err)
		fmt.Printf("Event body: %s\n", string(body))
		http.Error(w, "Failed to parse event", http.StatusBadRequest)
		return
	}

	// Only process opened, synchronize, and reopened events
	action := event.Action
	if action != "opened" && action != "synchronize" && action != "reopened" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "PR action %s ignored", action)
		return
	}

	// Log the event for debugging
	fmt.Printf("[INFO] Processing PR #%d action: %s from %s/%s\n",
		event.Number, action, event.Repository.Owner.Login, event.Repository.Name)

	// Publish to queue for processing
	if err := h.queue.Publish(event); err != nil {
		fmt.Printf("[ERROR] Failed to queue event: %v\n", err)
		http.Error(w, "Failed to queue event", http.StatusInternalServerError)
		return
	}

	// Quick response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Event queued for processing")
}

// verifySignature verifies the HMAC signature from GitHub
func (h *Handler) verifySignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}

	// GitHub sends signature as "sha256=<hash>"
	if len(signature) < 7 || signature[:7] != "sha256=" {
		return false
	}

	expectedHash := signature[7:]

	// Compute HMAC
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	actualHash := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedHash), []byte(actualHash))
}

// GitHubPullRequestEvent represents a GitHub pull request webhook event
type GitHubPullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		ID     int64   `json:"id"` // int64 for large GitHub IDs
		Number int     `json:"number"`
		Title  string  `json:"title"`
		Body   *string `json:"body"` // Pointer to handle null values
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
		Head struct {
			Sha string `json:"sha"`
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Sha  string `json:"sha"`
			Ref  string `json:"ref"`
			Repo struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				Owner    struct {
					Login string `json:"login"`
				} `json:"owner"`
			} `json:"repo"`
		} `json:"base"`
	} `json:"pull_request"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}
