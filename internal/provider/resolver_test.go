package provider_test

import (
	"testing"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider"
)

func TestResolverSelectsOllamaWhenLocalOnlyTrue(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	ollama := provider.NewMockOllamaProvider()
	r := provider.NewResolver(cloud, ollama)

	conn := connection.Connection{
		ProviderType: "cloud", // even if type says cloud...
		LocalOnly:    true,    // ...LocalOnly must win
	}

	got, err := r.Resolve(conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name() != "mock-ollama" {
		t.Errorf("expected mock-ollama, got %q", got.Name())
	}
}

func TestResolverSelectsCloudProviderWhenLocalOnlyFalse(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	ollama := provider.NewMockOllamaProvider()
	r := provider.NewResolver(cloud, ollama)

	conn := connection.Connection{
		ProviderType: "cloud",
		LocalOnly:    false,
	}

	got, err := r.Resolve(conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name() != "mock-cloud" {
		t.Errorf("expected mock-cloud, got %q", got.Name())
	}
}

func TestResolverSelectsOllamaByProviderType(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	ollama := provider.NewMockOllamaProvider()
	r := provider.NewResolver(cloud, ollama)

	conn := connection.Connection{ProviderType: "ollama"}

	got, err := r.Resolve(conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name() != "mock-ollama" {
		t.Errorf("expected mock-ollama, got %q", got.Name())
	}
}

func TestResolverReturnsErrorForUnknownProviderType(t *testing.T) {
	r := provider.NewResolver(
		provider.NewMockCloudProvider("cred-1"),
		provider.NewMockOllamaProvider(),
	)

	conn := connection.Connection{ProviderType: "gpt-banana"}

	_, err := r.Resolve(conn)
	if err == nil {
		t.Fatal("expected an error for unknown provider type, got nil")
	}
}

func TestResolverReturnsErrorWhenLocalOnlyButNoOllamaRegistered(t *testing.T) {
	r := provider.NewResolver(provider.NewMockCloudProvider("cred-1"), nil)

	conn := connection.Connection{LocalOnly: true}

	_, err := r.Resolve(conn)
	if err == nil {
		t.Fatal("expected error when ollama is nil and LocalOnly is true")
	}
}
