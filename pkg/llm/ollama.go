package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OllamaProvider implements the Provider interface for Ollama (local LLMs)
type OllamaProvider struct {
	url    string
	model  string
	client *http.Client
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(url, model string) *OllamaProvider {
	return &OllamaProvider{
		url:    url,
		model:  model,
		client: &http.Client{},
	}
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// Analyze sends code for analysis to Ollama
func (p *OllamaProvider) Analyze(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  p.model,
		"prompt": prompt,
		"stream": false,
		"system": "You are an expert code reviewer providing constructive feedback on pull requests.",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/generate", p.url)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}
