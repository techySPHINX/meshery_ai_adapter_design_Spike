package adapter

import (
	"context"
	"fmt"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/design"
	domerr "github.com/techysphinx/meshery-ai-adapter-design-spike/internal/errors"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider"
)

// Adapter is the orchestrator. It owns the pipeline:
//
//	intent → build prompt → resolve provider → generate → parse → validate → result
//
// It deliberately knows nothing about provider internals, credential values,
// or JSON parsing details. Each step delegates to its own package.
type Adapter struct {
	resolver      *provider.Resolver
	promptBuilder *prompt.Builder
	credStore     CredentialLookup
}

// CredentialLookup is the minimal interface the adapter needs from the
// credential store — just enough to confirm a credential exists before
// handing control to the provider. The actual secret is never returned here.
type CredentialLookup interface {
	Lookup(id string) (connection.CredentialRef, bool)
}

func New(resolver *provider.Resolver, pb *prompt.Builder, creds CredentialLookup) *Adapter {
	return &Adapter{
		resolver:      resolver,
		promptBuilder: pb,
		credStore:     creds,
	}
}

// GenerateDesign runs the full pipeline for a user intent against a connection.
func (a *Adapter) GenerateDesign(ctx context.Context, intent string, conn connection.Connection) (*design.Design, error) {
	// 1. Verify the credential reference exists (without touching the secret).
	if conn.CredentialID != "" {
		if _, ok := a.credStore.Lookup(conn.CredentialID); !ok {
			return nil, &domerr.CredentialError{CredentialID: conn.CredentialID}
		}
	}

	// 2. Resolve provider.
	p, err := a.resolver.Resolve(conn)
	if err != nil {
		return nil, fmt.Errorf("resolving provider: %w", err)
	}

	// 3. Health-check the chosen provider.
	if err := p.HealthCheck(ctx); err != nil {
		return nil, &domerr.ProviderError{Provider: p.Name(), Cause: err}
	}

	// 4. Build prompt.
	builtPrompt := a.promptBuilder.Build(intent)

	// 5. Call provider.
	resp, err := p.Generate(ctx, provider.PromptRequest{
		Prompt: builtPrompt,
		Model:  conn.Model,
	})
	if err != nil {
		return nil, &domerr.ProviderError{Provider: p.Name(), Cause: err}
	}

	// 6. Parse raw output (untrusted).
	parsed, err := design.ParseDesignResponse(resp.Raw)
	if err != nil {
		return nil, err
	}

	// 7. Validate parsed design.
	if err := design.ValidateDesign(parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}
