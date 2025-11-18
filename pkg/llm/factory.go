package llm

import (
	"fmt"
)

// Factory creates an LLM provider based on configuration
type Factory struct{}

// NewFactory creates a new LLM factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateProvider creates an LLM provider based on the provider type
func (f *Factory) CreateProvider(providerType string, config map[string]string) (Provider, error) {
	switch providerType {
	case "openai":
		apiKey := config["api_key"]
		model := config["model"]
		if apiKey == "" {
			return nil, fmt.Errorf("OpenAI API key is required")
		}
		if model == "" {
			model = "gpt-4-turbo-preview"
		}
		return NewOpenAIProvider(apiKey, model), nil

	case "anthropic":
		apiKey := config["api_key"]
		model := config["model"]
		if apiKey == "" {
			return nil, fmt.Errorf("Anthropic API key is required")
		}
		if model == "" {
			model = "claude-3-opus-20240229"
		}
		return NewAnthropicProvider(apiKey, model), nil

	case "ollama":
		url := config["url"]
		model := config["model"]
		if url == "" {
			url = "http://localhost:11434"
		}
		if model == "" {
			model = "codellama"
		}
		return NewOllamaProvider(url, model), nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", providerType)
	}
}
