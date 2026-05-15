package errors

import "fmt"

// These are the error types that can come out of the adapter pipeline.
// Each wraps a cause so callers can distinguish them without parsing strings.

type ProviderError struct {
	Provider string
	Cause    error
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider %q: %v", e.Provider, e.Cause)
}

func (e *ProviderError) Unwrap() error { return e.Cause }

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Message)
}

type ParseError struct {
	Cause error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("could not parse provider output: %v", e.Cause)
}

func (e *ParseError) Unwrap() error { return e.Cause }

type CredentialError struct {
	CredentialID string
}

func (e *CredentialError) Error() string {
	// intentionally omit the secret value — only log the ID
	return fmt.Sprintf("credential %q not found or unavailable", e.CredentialID)
}
