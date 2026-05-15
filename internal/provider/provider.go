package provider

import "context"

// PromptRequest is what the adapter hands to a provider.
// It contains only the rendered prompt text — no raw credentials,
// no connection secrets. Providers get auth injected at their own level.
type PromptRequest struct {
	Prompt string
	Model  string
}

// ProviderResponse is raw, untrusted output from an LLM.
// The caller (adapter) is responsible for parsing and validating it.
type ProviderResponse struct {
	Raw string // JSON or free text — not trusted yet
}

// LLMProvider is the single interface every provider must satisfy.
// Adding a new LLM backend means implementing these three methods;
// nothing else in the codebase changes.
type LLMProvider interface {
	Name() string
	Generate(ctx context.Context, req PromptRequest) (*ProviderResponse, error)
	HealthCheck(ctx context.Context) error
}
