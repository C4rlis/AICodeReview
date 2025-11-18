package llm

import (
	"context"
	"fmt"
)

// Provider defines the interface for LLM providers
type Provider interface {
	// Analyze sends code for analysis and returns feedback
	Analyze(ctx context.Context, prompt string) (string, error)
	// Name returns the provider name
	Name() string
}

// CodeReviewRequest represents a code review request
type CodeReviewRequest struct {
	RepositoryName string
	PullRequestID  int
	Diff           string
	FileChanges    []FileChange
	Author         string
	Title          string
	Description    string
}

// FileChange represents a changed file in a PR
type FileChange struct {
	Filename  string
	Status    string // "added", "modified", "removed"
	Additions int
	Deletions int
	Changes   int
	Patch     string
}

// CodeReviewResponse represents the AI's review feedback
type CodeReviewResponse struct {
	Summary  string
	Comments []ReviewComment
}

// ReviewComment represents a single review comment
type ReviewComment struct {
	Filename string
	Line     int
	Body     string
	Severity string // "info", "warning", "error"
}

// BuildPrompt creates a comprehensive prompt for code review
func BuildPrompt(req CodeReviewRequest) string {
	prompt := fmt.Sprintf(`You are an expert code reviewer. Analyze the following pull request and provide constructive feedback.

Repository: %s
Pull Request #%d: %s
Author: %s

Description:
%s

Changed Files:
`, req.RepositoryName, req.PullRequestID, req.Title, req.Author, req.Description)

	for _, file := range req.FileChanges {
		prompt += fmt.Sprintf("\n- %s (%s): +%d -%d", file.Filename, file.Status, file.Additions, file.Deletions)
	}

	prompt += "\n\nCode Changes:\n```diff\n" + req.Diff + "\n```\n\n"

	prompt += `Please review this code and provide:
1. A brief summary of the changes
2. Potential bugs or issues
3. Security concerns
4. Performance considerations
5. Code style and best practices
6. Suggestions for improvement

Format your response as JSON with the following structure:
{
  "summary": "Brief overview of the changes",
  "comments": [
    {
      "filename": "path/to/file.go",
      "line": 42,
      "body": "Detailed comment about this line",
      "severity": "info|warning|error"
    }
  ]
}

Focus on being constructive and helpful. Only mention issues if they are significant.`

	return prompt
}
