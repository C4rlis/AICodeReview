package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlr/codereviewtool/internal/analyzer"
	"github.com/carlr/codereviewtool/internal/config"
	"github.com/carlr/codereviewtool/internal/queue"
	"github.com/carlr/codereviewtool/internal/scm"
	"github.com/carlr/codereviewtool/internal/webhook"
	"github.com/carlr/codereviewtool/pkg/llm"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Code Review AI - Worker")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize LLM provider
	llmFactory := llm.NewFactory()
	llmConfig := map[string]string{}

	switch cfg.LLMProvider {
	case "openai":
		llmConfig["api_key"] = cfg.OpenAIAPIKey
		llmConfig["model"] = cfg.OpenAIModel
	case "anthropic":
		llmConfig["api_key"] = cfg.AnthropicAPIKey
		llmConfig["model"] = cfg.AnthropicModel
	case "ollama":
		llmConfig["url"] = cfg.OllamaURL
		llmConfig["model"] = cfg.OllamaModel
	}

	llmProvider, err := llmFactory.CreateProvider(cfg.LLMProvider, llmConfig)
	if err != nil {
		log.Fatalf("Failed to create LLM provider: %v", err)
	}

	log.Printf("Using LLM provider: %s", llmProvider.Name())

	// Initialize GitHub client
	githubClient := scm.NewGitHubClient(cfg.GitHubToken)

	// Initialize analyzer
	codeAnalyzer := analyzer.NewAnalyzer(llmProvider, githubClient)

	// Initialize RabbitMQ
	rabbitMQ, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	log.Println("Connected to RabbitMQ")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages
	go func() {
		err := rabbitMQ.Consume(func(body []byte) error {
			return processEvent(body, codeAnalyzer, githubClient)
		})
		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	log.Println("Worker is ready and waiting for pull requests...")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down worker...")
}

func processEvent(body []byte, codeAnalyzer *analyzer.Analyzer, githubClient *scm.GitHubClient) error {
	// Parse the webhook event
	var event webhook.GitHubPullRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	ctx := context.Background()

	owner := event.Repository.Owner.Login
	repo := event.Repository.Name
	prNumber := event.PullRequest.Number
	title := event.PullRequest.Title
	description := event.PullRequest.Body
	author := event.PullRequest.User.Login
	commitID := event.PullRequest.Head.Sha

	log.Printf("Processing PR #%d: %s from %s/%s (author: %s)", prNumber, title, owner, repo, author)

	// Analyze the pull request
	review, err := codeAnalyzer.AnalyzePullRequest(ctx, owner, repo, prNumber, title, description, author)
	if err != nil {
		log.Printf("Analysis failed: %v", err)
		return err
	}

	log.Printf("Analysis complete. Summary: %s", review.Summary)
	log.Printf("Found %d comments", len(review.Comments))

	// Post review to GitHub
	if len(review.Comments) > 0 {
		// Create a review with all comments
		githubReview := &scm.Review{
			Summary:  review.Summary,
			Comments: make([]scm.ReviewComment, 0, len(review.Comments)),
			CommitID: commitID,
		}

		for _, comment := range review.Comments {
			githubReview.Comments = append(githubReview.Comments, scm.ReviewComment{
				Filename: comment.Filename,
				Line:     comment.Line,
				Body:     formatComment(comment),
				CommitID: commitID,
			})
		}

		if err := githubClient.CreateReview(ctx, owner, repo, prNumber, githubReview); err != nil {
			log.Printf("Failed to create review: %v", err)
			// Fallback: post as summary comment
			if err := githubClient.PostReviewSummary(ctx, owner, repo, prNumber, formatReviewSummary(review)); err != nil {
				log.Printf("Failed to post review summary: %v", err)
				return err
			}
		}
	} else {
		// Just post the summary
		if err := githubClient.PostReviewSummary(ctx, owner, repo, prNumber, review.Summary); err != nil {
			log.Printf("Failed to post review summary: %v", err)
			return err
		}
	}

	log.Printf("Successfully posted review for PR #%d", prNumber)
	return nil
}

func formatComment(comment llm.ReviewComment) string {
	prefix := "[INFO]"
	switch comment.Severity {
	case "warning":
		prefix = "[WARNING]"
	case "error":
		prefix = "[ERROR]"
	}
	return prefix + " " + comment.Body
}

func formatReviewSummary(review *llm.CodeReviewResponse) string {
	summary := "## AI Code Review\n\n"
	summary += review.Summary + "\n\n"

	if len(review.Comments) > 0 {
		summary += "### Comments\n\n"
		for _, comment := range review.Comments {
			prefix := "[INFO]"
			switch comment.Severity {
			case "warning":
				prefix = "[WARNING]"
			case "error":
				prefix = "[ERROR]"
			}
			summary += fmt.Sprintf("- %s **%s**:%d - %s\n", prefix, comment.Filename, comment.Line, comment.Body)
		}
	}

	return summary
}
