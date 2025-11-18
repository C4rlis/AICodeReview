package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set required environment variables
	os.Setenv("GITHUB_WEBHOOK_SECRET", "test-secret")
	os.Setenv("GITHUB_TOKEN", "ghp_test123456")
	os.Setenv("LLM_PROVIDER", "openai")
	os.Setenv("OPENAI_API_KEY", "sk-test123")
	os.Setenv("RABBITMQ_URL", "amqp://localhost:5672")
	os.Setenv("POSTGRES_URL", "postgres://localhost:5432/test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.GitHubWebhookSecret != "test-secret" {
		t.Errorf("Expected GitHubWebhookSecret to be 'test-secret', got '%s'", cfg.GitHubWebhookSecret)
	}

	if cfg.LLMProvider != "openai" {
		t.Errorf("Expected LLMProvider to be 'openai', got '%s'", cfg.LLMProvider)
	}
}

func TestValidate_MissingGitHubSecret(t *testing.T) {
	cfg := &Config{GitHubToken: "test"}
	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing GitHub secret")
	}
}

func TestValidate_MissingGitHubToken(t *testing.T) {
	cfg := &Config{GitHubWebhookSecret: "test"}
	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing GitHub token")
	}
}

func TestValidate_InvalidLLMProvider(t *testing.T) {
	cfg := &Config{
		GitHubWebhookSecret: "test",
		GitHubToken:         "test",
		LLMProvider:         "invalid",
		RabbitMQURL:         "test",
		PostgresURL:         "test",
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid LLM provider")
	}
}

func TestValidate_OpenAINoAPIKey(t *testing.T) {
	cfg := &Config{
		GitHubWebhookSecret: "test",
		GitHubToken:         "test",
		LLMProvider:         "openai",
		RabbitMQURL:         "test",
		PostgresURL:         "test",
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing OpenAI API key")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &Config{
		GitHubWebhookSecret: "test",
		GitHubToken:         "test",
		LLMProvider:         "openai",
		OpenAIAPIKey:        "sk-test",
		RabbitMQURL:         "test",
		PostgresURL:         "test",
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}
