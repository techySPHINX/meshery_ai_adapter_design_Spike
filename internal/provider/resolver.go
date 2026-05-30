package provider

import (
	"fmt"
	"strings"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
)

// Resolver picks the right LLMProvider for a given Connection.
// It knows about the registered providers; it does not know about
// prompt building, parsing, or validation — those belong elsewhere.
type Resolver struct {
	providers     map[string]LLMProvider
	localOnlyType string
}

const defaultLocalProviderType = "ollama"

func NewResolver(cloud, ollama LLMProvider) *Resolver {
	providers := make(map[string]LLMProvider, 2)
	if cloud != nil {
		providers["cloud"] = cloud
	}
	if ollama != nil {
		providers["ollama"] = ollama
	}

	return NewResolverWithProviders(providers, defaultLocalProviderType)
}

// NewResolverWithProviders registers a flexible provider map for future backends.
func NewResolverWithProviders(providers map[string]LLMProvider, localOnlyType string) *Resolver {
	registry := make(map[string]LLMProvider, len(providers))
	for key, p := range providers {
		normalized := strings.ToLower(strings.TrimSpace(key))
		if normalized == "" || p == nil {
			continue
		}
		registry[normalized] = p
	}

	localOnlyType = strings.ToLower(strings.TrimSpace(localOnlyType))
	if localOnlyType == "" {
		localOnlyType = defaultLocalProviderType
	}

	return &Resolver{providers: registry, localOnlyType: localOnlyType}
}

// Resolve returns the provider that matches the connection's intent.
// LocalOnly connections always get the Ollama provider regardless of
// ProviderType, so user privacy preferences are respected.
func (r *Resolver) Resolve(conn connection.Connection) (LLMProvider, error) {
	if conn.LocalOnly {
		p, ok := r.providers[r.localOnlyType]
		if !ok || p == nil {
			return nil, fmt.Errorf("no local provider registered but connection requires local-only mode")
		}
		return p, nil
	}

	providerType := strings.ToLower(strings.TrimSpace(conn.ProviderType))
	if providerType == "" {
		return nil, fmt.Errorf("provider type is required")
	}

	p, ok := r.providers[providerType]
	if !ok || p == nil {
		return nil, fmt.Errorf("unknown provider type %q", conn.ProviderType)
	}

	return p, nil
}
