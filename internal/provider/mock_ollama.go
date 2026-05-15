package provider

import (
	"context"
	"fmt"
)

// MockOllamaProvider simulates a locally running Ollama instance.
// It never calls a real HTTP endpoint; it returns a canned response so
// tests stay fast and deterministic.
type MockOllamaProvider struct {
	// Healthy controls whether HealthCheck passes. Flip it in tests
	// to exercise the "local provider unavailable" failure path.
	Healthy bool

	// FixedResponse lets individual tests control what the "model" returns.
	FixedResponse string
}

func NewMockOllamaProvider() *MockOllamaProvider {
	return &MockOllamaProvider{Healthy: true}
}

func (p *MockOllamaProvider) Name() string { return "mock-ollama" }

func (p *MockOllamaProvider) HealthCheck(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !p.Healthy {
		return fmt.Errorf("ollama process is not running")
	}
	return nil
}

func (p *MockOllamaProvider) Generate(ctx context.Context, req PromptRequest) (*ProviderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !p.Healthy {
		return nil, fmt.Errorf("ollama process is not running")
	}

	body := p.FixedResponse
	if body == "" {
		// default happy-path fixture — a minimal valid design JSON
		body = `{
  "name": "ollama-generated-design",
  "components": [
    {"id": "comp-1", "name": "nginx", "type": "Deployment"}
  ],
  "relationships": []
}`
	}

	return &ProviderResponse{Raw: body}, nil
}
