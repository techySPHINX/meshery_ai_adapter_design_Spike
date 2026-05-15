package connection

// Connection holds everything the provider resolver needs to talk to an LLM.
// The actual secret lives elsewhere — we only store a reference ID here so it
// never leaks into logs, traces, or error messages.
type Connection struct {
	ID           string
	Name         string
	ProviderType string // "cloud" | "ollama"
	BaseURL      string
	Model        string
	CredentialID string // reference only, not the secret
	LocalOnly    bool   // when true, resolver must pick a local provider
}

// CredentialRef is a lightweight pointer to a stored secret.
// Think of it the same way Meshery handles Credentials: the credential
// authenticates the connection, but is never surfaced in API responses.
type CredentialRef struct {
	ID   string
	Type string // "api_key" | "bearer_token" | "none"
}

// CredentialStore is the simplest possible in-memory store for the spike.
// In real Meshery this would call the credentials API.
type CredentialStore struct {
	entries map[string]CredentialRef
}

func NewCredentialStore() *CredentialStore {
	return &CredentialStore{entries: make(map[string]CredentialRef)}
}

func (s *CredentialStore) Register(ref CredentialRef) {
	s.entries[ref.ID] = ref
}

// Lookup returns a CredentialRef without exposing its underlying secret.
// Callers only need to know the reference exists; the actual key material
// is injected at request time by the provider itself.
func (s *CredentialStore) Lookup(id string) (CredentialRef, bool) {
	ref, ok := s.entries[id]
	return ref, ok
}
