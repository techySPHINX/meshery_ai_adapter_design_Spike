package prompt

import (
	"fmt"
	"strings"
)

// SchemaContext holds the minimal schema rules injected into every prompt.
// Keeping this as a typed struct (rather than free text) makes it easy to
// extend later and prevents callers from accidentally embedding secrets.
type SchemaContext struct {
	ComponentRules   []string
	RelationshipRules []string
}

// DefaultSchemaContext is a sensible baseline that mirrors the simplified
// Meshery design model used in this spike.
var DefaultSchemaContext = SchemaContext{
	ComponentRules: []string{
		"component requires: id (string), name (string), type (string)",
	},
	RelationshipRules: []string{
		"relationship requires: source (component id), target (component id), kind (string)",
	},
}

// Builder constructs the prompt that gets sent to the LLM provider.
// It knows about intent and schema — it does not know about credentials,
// provider URLs, or API keys.
type Builder struct {
	schema SchemaContext
}

func NewBuilder(schema SchemaContext) *Builder {
	return &Builder{schema: schema}
}

// Build assembles the final prompt string.
// The format is intentionally plain so any LLM can follow it without
// model-specific templating.
func (b *Builder) Build(intent string) string {
	var sb strings.Builder

	sb.WriteString("You are a Meshery design assistant.\n\n")
	sb.WriteString("User intent:\n")
	fmt.Fprintf(&sb, "%s\n\n", strings.TrimSpace(intent))

	sb.WriteString("Schema rules (your output MUST follow these):\n")
	for _, r := range b.schema.ComponentRules {
		fmt.Fprintf(&sb, "  - %s\n", r)
	}
	for _, r := range b.schema.RelationshipRules {
		fmt.Fprintf(&sb, "  - %s\n", r)
	}

	sb.WriteString("\nRespond with a single JSON object matching this structure:\n")
	sb.WriteString(`{
  "name": "<design name>",
  "components": [{"id": "...", "name": "...", "type": "..."}],
  "relationships": [{"source": "...", "target": "...", "kind": "..."}]
}`)
	sb.WriteString("\n\nReturn only the JSON. No explanation, no markdown fences.")

	return sb.String()
}
