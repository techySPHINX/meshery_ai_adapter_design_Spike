package provider

import (
	"fmt"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
)

// Resolver picks the right LLMProvider for a given Connection.
// It knows about the registered providers; it does not know about
// prompt building, parsing, or validation — those belong elsewhere.
type Resolver struct {
	cloud  LLMProvider
	ollama LLMProvider
}

func NewResolver(cloud, ollama LLMProvider) *Resolver {
	return &Resolver{cloud: cloud, ollama: ollama}
}

// Resolve returns the provider that matches the connection's intent.
// LocalOnly connections always get the Ollama provider regardless of
// ProviderType, so user privacy preferences are respected.
func (r *Resolver) Resolve(conn connection.Connection) (LLMProvider, error) {
	if conn.LocalOnly {
		if r.ollama == nil {
			return nil, fmt.Errorf("no local provider registered but connection requires local-only mode")
		}
		return r.ollama, nil
	}

	switch conn.ProviderType {
	case "ollama":
		if r.ollama == nil {
			return nil, fmt.Errorf("ollama provider not registered")
		}
		return r.ollama, nil
	case "cloud":
		if r.cloud == nil {
			return nil, fmt.Errorf("cloud provider not registered")
		}
		return r.cloud, nil
	default:
		return nil, fmt.Errorf("unknown provider type %q", conn.ProviderType)
	}
}
