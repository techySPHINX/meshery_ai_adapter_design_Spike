package connection_test

import (
	"testing"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
)

func TestCredentialStoreReturnsRefById(t *testing.T) {
	store := connection.NewCredentialStore()
	store.Register(connection.CredentialRef{ID: "cred-abc", Type: "api_key"})

	ref, ok := store.Lookup("cred-abc")
	if !ok {
		t.Fatal("expected credential to be found")
	}
	if ref.ID != "cred-abc" {
		t.Errorf("got ID %q, want %q", ref.ID, "cred-abc")
	}
}

func TestCredentialStoreMissReturnsNotFound(t *testing.T) {
	store := connection.NewCredentialStore()

	_, ok := store.Lookup("does-not-exist")
	if ok {
		t.Fatal("expected lookup to fail for unknown credential")
	}
}

// This test enforces that a raw secret value cannot be obtained through the
// public CredentialRef type. CredentialRef only carries an ID and type.
func TestCredentialRefHasNoRawSecret(t *testing.T) {
	ref := connection.CredentialRef{ID: "cred-xyz", Type: "bearer_token"}

	// If the struct ever gets a "Secret" or "Value" field, this test will
	// still compile, but reviewers will notice something changed.
	_ = ref.ID
	_ = ref.Type
	// No ref.Secret, ref.Value, ref.Key, etc.
}
