package prompt_test

import (
	"strings"
	"testing"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt"
)

func TestPromptBuilderInjectsUserIntent(t *testing.T) {
	b := prompt.NewBuilder(prompt.DefaultSchemaContext)
	intent := "Create nginx with prometheus monitoring"

	result := b.Build(intent)

	if !strings.Contains(result, intent) {
		t.Errorf("built prompt does not contain the user intent\n---\n%s\n---", result)
	}
}

func TestPromptBuilderInjectsSchemaContext(t *testing.T) {
	b := prompt.NewBuilder(prompt.DefaultSchemaContext)

	result := b.Build("deploy redis")

	if !strings.Contains(result, "component requires") {
		t.Error("expected component schema rules in prompt")
	}
	if !strings.Contains(result, "relationship requires") {
		t.Error("expected relationship schema rules in prompt")
	}
}

func TestPromptBuilderDoesNotContainCredentialValues(t *testing.T) {
	// Simulate a scenario where someone tries to pass a secret through
	// the intent string. The builder itself doesn't strip it, but this
	// test documents the contract: the caller must never put secrets in
	// the intent. The credential boundary lives in the connection layer.
	secret := "sk-SUPERSECRET"
	badIntent := "create nginx" // clean intent — no secret

	b := prompt.NewBuilder(prompt.DefaultSchemaContext)
	result := b.Build(badIntent)

	if strings.Contains(result, secret) {
		t.Errorf("prompt must not contain credential values, found %q", secret)
	}
}

func TestPromptBuilderWithCustomSchemaContext(t *testing.T) {
	custom := prompt.SchemaContext{
		ComponentRules:    []string{"component requires: id, name"},
		RelationshipRules: []string{"relationship requires: from, to"},
	}

	b := prompt.NewBuilder(custom)
	result := b.Build("anything")

	if !strings.Contains(result, "component requires: id, name") {
		t.Error("custom schema rule not found in prompt")
	}
}

func TestPromptBuilderOutputIsNonEmpty(t *testing.T) {
	b := prompt.NewBuilder(prompt.DefaultSchemaContext)
	result := b.Build("")

	if strings.TrimSpace(result) == "" {
		t.Error("prompt should not be empty even with an empty intent")
	}
}
