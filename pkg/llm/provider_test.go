package llm

import (
	"context"
	"testing"
)

type mockProvider struct {
	name     string
	response string
	err      error
}

func (m *mockProvider) Analyze(ctx context.Context, prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func (m *mockProvider) Name() string {
	return m.name
}

func TestBuildPrompt(t *testing.T) {
	req := CodeReviewRequest{
		RepositoryName: "user/repo",
		PullRequestID:  42,
		Title:          "Fix bug",
		Author:         "testuser",
		Description:    "This fixes a bug",
		Diff:           "diff content here",
		FileChanges: []FileChange{
			{
				Filename:  "main.go",
				Status:    "modified",
				Additions: 5,
				Deletions: 2,
			},
		},
	}

	prompt := BuildPrompt(req)

	// Check that prompt contains key information
	if len(prompt) == 0 {
		t.Error("Expected non-empty prompt")
	}

	tests := []string{
		"user/repo",
		"Pull Request #42",
		"Fix bug",
		"testuser",
		"main.go",
		"diff content here",
	}

	for _, expected := range tests {
		if !contains(prompt, expected) {
			t.Errorf("Expected prompt to contain '%s'", expected)
		}
	}
}

func TestFactory_CreateOpenAIProvider(t *testing.T) {
	factory := NewFactory()
	config := map[string]string{
		"api_key": "sk-test123",
		"model":   "gpt-4",
	}

	provider, err := factory.CreateProvider("openai", config)
	if err != nil {
		t.Fatalf("Failed to create OpenAI provider: %v", err)
	}

	if provider.Name() != "openai" {
		t.Errorf("Expected provider name 'openai', got '%s'", provider.Name())
	}
}

func TestFactory_CreateAnthropicProvider(t *testing.T) {
	factory := NewFactory()
	config := map[string]string{
		"api_key": "sk-ant-test123",
		"model":   "claude-3-opus-20240229",
	}

	provider, err := factory.CreateProvider("anthropic", config)
	if err != nil {
		t.Fatalf("Failed to create Anthropic provider: %v", err)
	}

	if provider.Name() != "anthropic" {
		t.Errorf("Expected provider name 'anthropic', got '%s'", provider.Name())
	}
}

func TestFactory_CreateOllamaProvider(t *testing.T) {
	factory := NewFactory()
	config := map[string]string{
		"url":   "http://localhost:11434",
		"model": "codellama",
	}

	provider, err := factory.CreateProvider("ollama", config)
	if err != nil {
		t.Fatalf("Failed to create Ollama provider: %v", err)
	}

	if provider.Name() != "ollama" {
		t.Errorf("Expected provider name 'ollama', got '%s'", provider.Name())
	}
}

func TestFactory_UnsupportedProvider(t *testing.T) {
	factory := NewFactory()
	config := map[string]string{}

	_, err := factory.CreateProvider("unsupported", config)
	if err == nil {
		t.Error("Expected error for unsupported provider")
	}
}

func TestFactory_MissingAPIKey(t *testing.T) {
	factory := NewFactory()
	config := map[string]string{
		"model": "gpt-4",
	}

	_, err := factory.CreateProvider("openai", config)
	if err == nil {
		t.Error("Expected error for missing API key")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
