package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// DefaultOllamaBaseURL is the standard local endpoint for Ollama.
const DefaultOllamaBaseURL = "http://localhost:11434"

// OllamaHTTPProvider talks to a locally running Ollama server over HTTP.
type OllamaHTTPProvider struct {
	baseURL string
	client  *http.Client
}

// NewOllamaHTTPProvider builds a provider with a sane timeout and base URL.
func NewOllamaHTTPProvider(baseURL string) *OllamaHTTPProvider {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		trimmed = DefaultOllamaBaseURL
	}

	return &OllamaHTTPProvider{
		baseURL: strings.TrimRight(trimmed, "/"),
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *OllamaHTTPProvider) Name() string { return "ollama-http" }

func (p *OllamaHTTPProvider) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/api/tags", nil)
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama health check failed: %s", resp.Status)
	}

	var tags ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return err
	}

	return nil
}

func (p *OllamaHTTPProvider) Generate(ctx context.Context, req PromptRequest) (*ProviderResponse, error) {
	model := strings.TrimSpace(req.Model)
	if model == "" {
		return nil, fmt.Errorf("ollama model is required")
	}
	prompt := strings.TrimSpace(req.Prompt)
	if prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	payload := ollamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama generate failed: %s", resp.Status)
	}

	var out ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if strings.TrimSpace(out.Response) == "" {
		return nil, fmt.Errorf("ollama returned an empty response")
	}

	return &ProviderResponse{Raw: out.Response}, nil
}

type ollamaTagsResponse struct {
	Models []ollamaTag `json:"models"`
}

type ollamaTag struct {
	Name string `json:"name"`
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
}
