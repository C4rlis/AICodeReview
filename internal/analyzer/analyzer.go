package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/carlr/codereviewtool/internal/scm"
	"github.com/carlr/codereviewtool/pkg/llm"
)

// Analyzer performs AI-powered code analysis
type Analyzer struct {
	llmProvider  llm.Provider
	githubClient *scm.GitHubClient
}

// NewAnalyzer creates a new code analyzer
func NewAnalyzer(llmProvider llm.Provider, githubClient *scm.GitHubClient) *Analyzer {
	return &Analyzer{
		llmProvider:  llmProvider,
		githubClient: githubClient,
	}
}

// AnalyzePullRequest analyzes a pull request and returns review feedback
func (a *Analyzer) AnalyzePullRequest(ctx context.Context, owner, repo string, prNumber int, title, description, author string) (*llm.CodeReviewResponse, error) {
	log.Printf("Analyzing PR #%d in %s/%s", prNumber, owner, repo)

	// Get PR diff
	diff, err := a.githubClient.GetPullRequestDiff(ctx, owner, repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR diff: %w", err)
	}

	// Get changed files
	files, err := a.githubClient.GetPullRequestFiles(ctx, owner, repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR files: %w", err)
	}

	// Build file changes list
	fileChanges := make([]llm.FileChange, 0, len(files))
	for _, file := range files {
		fc := llm.FileChange{
			Filename:  file.GetFilename(),
			Status:    file.GetStatus(),
			Additions: file.GetAdditions(),
			Deletions: file.GetDeletions(),
			Changes:   file.GetChanges(),
			Patch:     file.GetPatch(),
		}
		fileChanges = append(fileChanges, fc)
	}

	// Build review request
	request := llm.CodeReviewRequest{
		RepositoryName: fmt.Sprintf("%s/%s", owner, repo),
		PullRequestID:  prNumber,
		Diff:           diff,
		FileChanges:    fileChanges,
		Author:         author,
		Title:          title,
		Description:    description,
	}

	// Build prompt
	prompt := llm.BuildPrompt(request)

	log.Printf("Sending code to %s for analysis...", a.llmProvider.Name())

	// Get AI analysis
	response, err := a.llmProvider.Analyze(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	log.Printf("Received response from %s", a.llmProvider.Name())

	// Parse response
	var reviewResponse llm.CodeReviewResponse
	if err := json.Unmarshal([]byte(response), &reviewResponse); err != nil {
		// If parsing fails, try to extract JSON from markdown code block
		cleanedResponse := extractJSONFromMarkdown(response)
		if err := json.Unmarshal([]byte(cleanedResponse), &reviewResponse); err != nil {
			// If still fails, return raw response as summary
			log.Printf("Failed to parse LLM response as JSON, using raw response")
			reviewResponse = llm.CodeReviewResponse{
				Summary:  response,
				Comments: []llm.ReviewComment{},
			}
		}
	}

	return &reviewResponse, nil
}

// extractJSONFromMarkdown attempts to extract JSON from markdown code blocks
func extractJSONFromMarkdown(text string) string {
	// Simple extraction of JSON from ```json ... ``` blocks
	start := -1
	for i := 0; i < len(text)-6; i++ {
		if text[i:i+7] == "```json" {
			start = i + 7
			break
		}
		if text[i:i+3] == "```" && (text[i+3] == '\n' || text[i+3] == '{') {
			start = i + 3
			break
		}
	}

	if start == -1 {
		return text
	}

	end := -1
	for i := start; i < len(text)-2; i++ {
		if text[i:i+3] == "```" {
			end = i
			break
		}
	}

	if end == -1 {
		return text[start:]
	}

	return text[start:end]
}
