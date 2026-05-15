package adapter_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/adapter"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider"
)

const (
	ollamaIntegrationEnv = "OLLAMA_INTEGRATION"
	ollamaBaseURLEnv     = "OLLAMA_BASE_URL"
	ollamaModelEnv       = "OLLAMA_MODEL"
)

type realOllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
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

func newRealOllamaProvider(baseURL, model string) *realOllamaProvider {
	return &realOllamaProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *realOllamaProvider) Name() string { return "ollama-http" }

func (p *realOllamaProvider) HealthCheck(ctx context.Context) error {
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

	if !modelPresent(tags.Models, p.model) {
		return fmt.Errorf("ollama model %q not found; run `ollama pull %s`", p.model, p.model)
	}

	return nil
}

func (p *realOllamaProvider) Generate(ctx context.Context, req provider.PromptRequest) (*provider.ProviderResponse, error) {
	prompt := req.Prompt + "\n\nReturn the following minimal valid JSON exactly, with no extra text:\n" +
		`{"name":"ollama-integration","components":[{"id":"comp-1","name":"nginx","type":"Deployment"}],"relationships":[]}`

	payload := ollamaGenerateRequest{
		Model:  p.model,
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

	return &provider.ProviderResponse{Raw: out.Response}, nil
}

func modelPresent(models []ollamaTag, want string) bool {
	for _, m := range models {
		if m.Name == want {
			return true
		}
		if strings.HasPrefix(m.Name, want+":") {
			return true
		}
	}
	return false
}

func TestAdapterWithRealOllamaIntegration(t *testing.T) {
	if os.Getenv(ollamaIntegrationEnv) != "1" {
		t.Skip("set OLLAMA_INTEGRATION=1 to enable real Ollama integration testing")
	}

	model := strings.TrimSpace(os.Getenv(ollamaModelEnv))
	if model == "" {
		t.Skip("set OLLAMA_MODEL to a local model name, for example: llama3")
	}

	baseURL := strings.TrimSpace(os.Getenv(ollamaBaseURLEnv))
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	ollama := newRealOllamaProvider(baseURL, model)
	resolver := provider.NewResolver(provider.NewMockCloudProvider("cred-1"), ollama)
	pb := prompt.NewBuilder(prompt.DefaultSchemaContext)
	store := connection.NewCredentialStore()

	a := adapter.New(resolver, pb, store)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	conn := connection.Connection{
		ProviderType: "ollama",
		LocalOnly:    true,
		Model:        model,
	}

	d, err := a.GenerateDesign(ctx, "Create a single nginx deployment", conn)
	if err != nil {
		t.Fatalf("unexpected error from real Ollama provider: %v", err)
	}
	if strings.TrimSpace(d.Name) == "" {
		t.Fatal("expected a non-empty design name from Ollama")
	}
	if len(d.Components) == 0 {
		t.Fatal("expected at least one component in the design")
	}
}
