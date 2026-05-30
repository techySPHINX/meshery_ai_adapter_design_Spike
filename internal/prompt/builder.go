package prompt

import (
	"fmt"
	"strings"
)

// SchemaContext holds the minimal schema rules injected into every prompt.
// Keeping this as a typed struct (rather than free text) makes it easy to
// extend later and prevents callers from accidentally embedding secrets.
type SchemaContext struct {
	ComponentRules    []string
	RelationshipRules []string
}

// DefaultSchemaContext is a sensible baseline that mirrors the simplified
// Meshery design model used in this spike.
var DefaultSchemaContext = SchemaContext{
	ComponentRules: []string{
		"component requires: id (string), name (string), type (string)",
		"component type must match a Meshery model schema entry",
		"encode protocol detail in component type or name when needed",
	},
	RelationshipRules: []string{
		"relationship requires: source (component id), target (component id), kind (string)",
		"relationship kind should align with the model schema guidance",
	},
}

// Builder constructs the prompt that gets sent to the LLM provider.
// It knows about intent and schema — it does not know about credentials,
// provider URLs, or API keys.
type Builder struct {
	schema       SchemaContext
	catalog      ModelCatalog
	window       ContextWindow
	systemPrompt string
}

// Option customizes how the prompt is assembled.
type Option func(*Builder)

// WithModelCatalog overrides the default Meshery model catalog.
func WithModelCatalog(catalog ModelCatalog) Option {
	return func(b *Builder) {
		b.catalog = catalog
	}
}

// WithContextWindow overrides the default context window settings.
func WithContextWindow(window ContextWindow) Option {
	return func(b *Builder) {
		b.window = window
	}
}

// WithSystemPrompt overrides the default system prompt.
func WithSystemPrompt(prompt string) Option {
	return func(b *Builder) {
		b.systemPrompt = prompt
	}
}

// NewBuilder assembles a prompt builder with defaults, plus optional overrides.
func NewBuilder(schema SchemaContext, opts ...Option) *Builder {
	b := &Builder{
		schema:       schema,
		catalog:      DefaultModelCatalog(),
		window:       DefaultContextWindow,
		systemPrompt: DefaultSystemPrompt,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Build assembles the final prompt string.
// The format is intentionally plain so any LLM can follow it without
// model-specific templating.
func (b *Builder) Build(intent string) string {
	var sb strings.Builder

	sb.WriteString("System prompt:\n")
	sb.WriteString(b.systemPrompt)

	sb.WriteString("\n\nUser intent:\n")
	fmt.Fprintf(&sb, "%s\n", strings.TrimSpace(intent))

	sb.WriteString("\nRelevant Meshery model schemas:\n")
	sb.WriteString(renderModelContext(b.catalog, b.window, intent))

	sb.WriteString("\nSchema rules (your output MUST follow these):\n")
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

func renderModelContext(catalog ModelCatalog, window ContextWindow, intent string) string {
	models := catalog.Select(intent, window)
	if len(models) == 0 {
		return "  - (no model definitions available)\n"
	}

	var sb strings.Builder
	for _, m := range models {
		sb.WriteString(m.Render())
	}

	return sb.String()
}
