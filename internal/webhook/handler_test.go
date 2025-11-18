package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockQueue struct {
	published []interface{}
}

func (m *mockQueue) Publish(event interface{}) error {
	m.published = append(m.published, event)
	return nil
}

func TestHandleGitHub_ValidSignature(t *testing.T) {
	secret := "test-secret"
	queue := &mockQueue{}
	handler := NewHandler(secret, queue)

	// Create test payload
	event := GitHubPullRequestEvent{
		Action: "opened",
		Number: 1,
	}
	event.PullRequest.Number = 1
	event.PullRequest.Title = "Test PR"

	body, _ := json.Marshal(event)

	// Create signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// Create request
	req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "pull_request")
	req.Header.Set("X-Hub-Signature-256", signature)

	w := httptest.NewRecorder()
	handler.HandleGitHub(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if len(queue.published) != 1 {
		t.Errorf("Expected 1 event published, got %d", len(queue.published))
	}
}

func TestHandleGitHub_InvalidSignature(t *testing.T) {
	secret := "test-secret"
	queue := &mockQueue{}
	handler := NewHandler(secret, queue)

	event := GitHubPullRequestEvent{Action: "opened"}
	body, _ := json.Marshal(event)

	req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "pull_request")
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid")

	w := httptest.NewRecorder()
	handler.HandleGitHub(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	if len(queue.published) != 0 {
		t.Errorf("Expected 0 events published, got %d", len(queue.published))
	}
}

func TestHandleGitHub_IgnoredEventType(t *testing.T) {
	secret := "test-secret"
	queue := &mockQueue{}
	handler := NewHandler(secret, queue)

	event := map[string]string{"action": "test"}
	body, _ := json.Marshal(event)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature-256", signature)

	w := httptest.NewRecorder()
	handler.HandleGitHub(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if len(queue.published) != 0 {
		t.Errorf("Expected 0 events published for ignored event type, got %d", len(queue.published))
	}
}

func TestHandleGitHub_IgnoredAction(t *testing.T) {
	secret := "test-secret"
	queue := &mockQueue{}
	handler := NewHandler(secret, queue)

	event := GitHubPullRequestEvent{Action: "closed"}
	body, _ := json.Marshal(event)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "pull_request")
	req.Header.Set("X-Hub-Signature-256", signature)

	w := httptest.NewRecorder()
	handler.HandleGitHub(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if len(queue.published) != 0 {
		t.Errorf("Expected 0 events published for ignored action, got %d", len(queue.published))
	}
}
