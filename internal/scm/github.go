package scm

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
)

// GitHubClient wraps the GitHub API client
type GitHubClient struct {
	client *github.Client
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string) *GitHubClient {
	client := github.NewClient(nil).WithAuthToken(token)
	return &GitHubClient{
		client: client,
	}
}

// GetPullRequestDiff retrieves the diff for a pull request
func (g *GitHubClient) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	diff, _, err := g.client.PullRequests.GetRaw(
		ctx,
		owner,
		repo,
		prNumber,
		github.RawOptions{Type: github.Diff},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get PR diff: %w", err)
	}
	return diff, nil
}

// GetPullRequestFiles retrieves the list of files changed in a PR
func (g *GitHubClient) GetPullRequestFiles(ctx context.Context, owner, repo string, prNumber int) ([]*github.CommitFile, error) {
	files, _, err := g.client.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list PR files: %w", err)
	}
	return files, nil
}

// PostReviewComment posts a review comment on a specific line of a file
func (g *GitHubClient) PostReviewComment(ctx context.Context, owner, repo string, prNumber int, comment *ReviewComment) error {
	githubComment := &github.PullRequestComment{
		Body:     github.String(comment.Body),
		Path:     github.String(comment.Filename),
		Line:     github.Int(comment.Line),
		CommitID: github.String(comment.CommitID),
	}

	_, _, err := g.client.PullRequests.CreateComment(ctx, owner, repo, prNumber, githubComment)
	if err != nil {
		return fmt.Errorf("failed to post review comment: %w", err)
	}
	return nil
}

// PostReviewSummary posts a general review comment on the PR
func (g *GitHubClient) PostReviewSummary(ctx context.Context, owner, repo string, prNumber int, body string) error {
	comment := &github.IssueComment{
		Body: github.String(body),
	}

	_, _, err := g.client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to post review summary: %w", err)
	}
	return nil
}

// CreateReview creates a review with multiple comments
func (g *GitHubClient) CreateReview(ctx context.Context, owner, repo string, prNumber int, review *Review) error {
	comments := make([]*github.DraftReviewComment, 0, len(review.Comments))
	for _, c := range review.Comments {
		comments = append(comments, &github.DraftReviewComment{
			Path: github.String(c.Filename),
			Line: github.Int(c.Line),
			Body: github.String(c.Body),
		})
	}

	githubReview := &github.PullRequestReviewRequest{
		Body:     github.String(review.Summary),
		Event:    github.String("COMMENT"),
		Comments: comments,
		CommitID: github.String(review.CommitID),
	}

	_, _, err := g.client.PullRequests.CreateReview(ctx, owner, repo, prNumber, githubReview)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

// ReviewComment represents a single review comment
type ReviewComment struct {
	Filename string
	Line     int
	Body     string
	CommitID string
}

// Review represents a complete review with multiple comments
type Review struct {
	Summary  string
	Comments []ReviewComment
	CommitID string
}
