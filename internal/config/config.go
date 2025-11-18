package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// GitHub
	GitHubWebhookSecret string
	GitHubToken         string

	// LLM Provider
	LLMProvider     string // "openai", "anthropic", "ollama"
	OpenAIAPIKey    string
	OpenAIModel     string
	AnthropicAPIKey string
	AnthropicModel  string
	OllamaURL       string
	OllamaModel     string

	// Infrastructure
	RabbitMQURL string
	PostgresURL string

	// Server
	WebhookPort      string
	WorkerConcurrency int

	// Logging
	LogLevel string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		// GitHub
		GitHubWebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),
		GitHubToken:         getEnv("GITHUB_TOKEN", ""),

		// LLM Provider
		LLMProvider:     getEnv("LLM_PROVIDER", "openai"),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:     getEnv("OPENAI_MODEL", "gpt-4-turbo-preview"),
		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		AnthropicModel:  getEnv("ANTHROPIC_MODEL", "claude-3-opus-20240229"),
		OllamaURL:       getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:     getEnv("OLLAMA_MODEL", "codellama"),

		// Infrastructure
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/codereview?sslmode=disable"),

		// Server
		WebhookPort:       getEnv("WEBHOOK_PORT", "8080"),
		WorkerConcurrency: getEnvInt("WORKER_CONCURRENCY", 5),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration values are present
func (c *Config) Validate() error {
	if c.GitHubWebhookSecret == "" {
		return fmt.Errorf("GITHUB_WEBHOOK_SECRET is required")
	}
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}

	// Validate LLM provider configuration
	switch c.LLMProvider {
	case "openai":
		if c.OpenAIAPIKey == "" {
			return fmt.Errorf("OPENAI_API_KEY is required when using OpenAI provider")
		}
	case "anthropic":
		if c.AnthropicAPIKey == "" {
			return fmt.Errorf("ANTHROPIC_API_KEY is required when using Anthropic provider")
		}
	case "ollama":
		if c.OllamaURL == "" {
			return fmt.Errorf("OLLAMA_URL is required when using Ollama provider")
		}
	default:
		return fmt.Errorf("invalid LLM_PROVIDER: %s (must be openai, anthropic, or ollama)", c.LLMProvider)
	}

	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.PostgresURL == "" {
		return fmt.Errorf("POSTGRES_URL is required")
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an environment variable as an integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
