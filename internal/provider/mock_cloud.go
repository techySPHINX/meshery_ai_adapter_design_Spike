package provider

import (
	"context"
	"fmt"
)

// MockCloudProvider simulates a cloud-hosted LLM (OpenAI, etc.).
// Like the Ollama mock it never makes a real HTTP call.
// The CredentialID field stores only the reference — the raw key
// is intentionally absent to match the connection/credential boundary.
type MockCloudProvider struct {
	Healthy      bool
	CredentialID string // reference only, never the secret itself

	FixedResponse string
}

func NewMockCloudProvider(credentialID string) *MockCloudProvider {
	return &MockCloudProvider{
		Healthy:      true,
		CredentialID: credentialID,
	}
}

func (p *MockCloudProvider) Name() string { return "mock-cloud" }

func (p *MockCloudProvider) HealthCheck(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !p.Healthy {
		return fmt.Errorf("cloud provider unreachable")
	}
	return nil
}

func (p *MockCloudProvider) Generate(ctx context.Context, req PromptRequest) (*ProviderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !p.Healthy {
		return nil, fmt.Errorf("cloud provider unreachable")
	}

	body := p.FixedResponse
	if body == "" {
		body = `{
  "name": "cloud-generated-design",
  "components": [
    {"id": "comp-1", "name": "nginx",      "type": "Deployment"},
    {"id": "comp-2", "name": "prometheus", "type": "Service"}
  ],
  "relationships": [
    {"source": "comp-1", "target": "comp-2", "kind": "observes"}
  ]
}`
	}

	return &ProviderResponse{Raw: body}, nil
}
