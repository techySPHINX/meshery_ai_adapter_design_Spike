package adapter_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/adapter"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
	domerr "github.com/techysphinx/meshery-ai-adapter-design-spike/internal/errors"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider"
)

func newTestAdapter(cloud *provider.MockCloudProvider, ollama *provider.MockOllamaProvider) *adapter.Adapter {
	store := connection.NewCredentialStore()
	store.Register(connection.CredentialRef{ID: "cred-1", Type: "api_key"})

	resolver := provider.NewResolver(cloud, ollama)
	pb := prompt.NewBuilder(prompt.DefaultSchemaContext)
	return adapter.New(resolver, pb, store)
}

func TestAdapterReturnsValidatedDesignFromMockCloudProvider(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	a := newTestAdapter(cloud, provider.NewMockOllamaProvider())

	conn := connection.Connection{
		ProviderType: "cloud",
		CredentialID: "cred-1",
	}

	d, err := a.GenerateDesign(context.Background(), "Create nginx with prometheus", conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name == "" {
		t.Error("expected a non-empty design name")
	}
	if len(d.Components) == 0 {
		t.Error("expected at least one component")
	}
}

func TestAdapterReturnsValidatedDesignFromMockOllamaProvider(t *testing.T) {
	ollama := provider.NewMockOllamaProvider()
	a := newTestAdapter(provider.NewMockCloudProvider("cred-1"), ollama)

	conn := connection.Connection{
		ProviderType: "ollama",
		LocalOnly:    true,
	}

	d, err := a.GenerateDesign(context.Background(), "Deploy redis", conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name == "" {
		t.Error("expected a non-empty design name")
	}
}

func TestAdapterFailsWhenCredentialNotFound(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	a := newTestAdapter(cloud, provider.NewMockOllamaProvider())

	conn := connection.Connection{
		ProviderType: "cloud",
		CredentialID: "cred-that-does-not-exist",
	}

	_, err := a.GenerateDesign(context.Background(), "something", conn)
	if err == nil {
		t.Fatal("expected credential error, got nil")
	}

	var ce *domerr.CredentialError
	if !asError(err, &ce) {
		t.Errorf("expected *CredentialError, got %T: %v", err, err)
	}
}

func TestAdapterFailsWhenProviderUnhealthy(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	cloud.Healthy = false

	a := newTestAdapter(cloud, provider.NewMockOllamaProvider())

	conn := connection.Connection{ProviderType: "cloud", CredentialID: "cred-1"}

	_, err := a.GenerateDesign(context.Background(), "anything", conn)
	if err == nil {
		t.Fatal("expected provider error, got nil")
	}

	var pe *domerr.ProviderError
	if !asError(err, &pe) {
		t.Errorf("expected *ProviderError, got %T: %v", err, err)
	}
}

func TestAdapterFailsWhenLocalProviderUnhealthy(t *testing.T) {
	ollama := provider.NewMockOllamaProvider()
	ollama.Healthy = false

	a := newTestAdapter(provider.NewMockCloudProvider("cred-1"), ollama)

	conn := connection.Connection{
		ProviderType: "ollama",
		LocalOnly:    true,
	}

	_, err := a.GenerateDesign(context.Background(), "local provider down", conn)
	if err == nil {
		t.Fatal("expected provider error, got nil")
	}

	var pe *domerr.ProviderError
	if !asError(err, &pe) {
		t.Errorf("expected *ProviderError, got %T: %v", err, err)
	}
}

func TestAdapterFailsOnProviderTimeout(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	a := newTestAdapter(cloud, provider.NewMockOllamaProvider())

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	conn := connection.Connection{
		ProviderType: "cloud",
		CredentialID: "cred-1",
	}

	_, err := a.GenerateDesign(ctx, "timeout path", conn)
	if err == nil {
		t.Fatal("expected provider error due to context deadline, got nil")
	}

	var pe *domerr.ProviderError
	if !asError(err, &pe) {
		t.Errorf("expected *ProviderError, got %T: %v", err, err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context deadline exceeded, got %v", err)
	}
}

func TestAdapterFailsOnInvalidProviderOutput(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	cloud.FixedResponse = `{"name": "", "components": [{"id":"","name":"","type":""}], "relationships": []}`

	a := newTestAdapter(cloud, nil)

	conn := connection.Connection{ProviderType: "cloud", CredentialID: "cred-1"}

	_, err := a.GenerateDesign(context.Background(), "bad output path", conn)
	if err == nil {
		t.Fatal("expected validation error from empty fields, got nil")
	}
}

func TestAdapterFailsOnMalformedJSON(t *testing.T) {
	cloud := provider.NewMockCloudProvider("cred-1")
	cloud.FixedResponse = `this is not json {`

	a := newTestAdapter(cloud, nil)

	conn := connection.Connection{ProviderType: "cloud", CredentialID: "cred-1"}

	_, err := a.GenerateDesign(context.Background(), "malformed response", conn)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}

	var pe *domerr.ParseError
	if !asError(err, &pe) {
		t.Errorf("expected *ParseError, got %T: %v", err, err)
	}
}

func TestAdapterLocalOnlyRoutesToOllama(t *testing.T) {
	ollama := provider.NewMockOllamaProvider()
	// Cloud provider is healthy, but LocalOnly should bypass it.
	cloud := provider.NewMockCloudProvider("cred-1")

	a := newTestAdapter(cloud, ollama)

	conn := connection.Connection{
		ProviderType: "cloud", // deliberately says cloud
		LocalOnly:    true,    // but LocalOnly must win
	}

	d, err := a.GenerateDesign(context.Background(), "local only test", conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The ollama mock returns "ollama-generated-design"
	if d.Name != "ollama-generated-design" {
		t.Errorf("expected ollama response, got design name %q", d.Name)
	}
}

// asError is a tiny type-assert helper so the test file has no imports beyond
// the project's own packages and stdlib.
func asError(err error, target interface{}) bool {
	switch t := target.(type) {
	case **domerr.CredentialError:
		if v, ok := err.(*domerr.CredentialError); ok {
			*t = v
			return true
		}
	case **domerr.ProviderError:
		if v, ok := err.(*domerr.ProviderError); ok {
			*v = *err.(*domerr.ProviderError)
			*t = v
			return true
		}
	case **domerr.ParseError:
		if v, ok := err.(*domerr.ParseError); ok {
			*t = v
			return true
		}
	}
	return false
}
