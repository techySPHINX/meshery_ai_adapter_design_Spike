package design

import (
	"fmt"
	"strings"

	domerr "github.com/techysphinx/meshery-ai-adapter-design-spike/internal/errors"
)

// ValidateDesign checks that a parsed Design satisfies the field-level
// contract expected by the rest of the system.
// It returns the first error it finds; callers shouldn't assume partial
// validity when an error is returned.
func ValidateDesign(d *Design) error {
	if strings.TrimSpace(d.Name) == "" {
		return &domerr.ValidationError{Field: "name", Message: "design name is required"}
	}

	// build an index of valid component IDs for relationship checks below
	componentIDs := make(map[string]struct{}, len(d.Components))

	for i, c := range d.Components {
		loc := fmt.Sprintf("components[%d]", i)

		if strings.TrimSpace(c.ID) == "" {
			return &domerr.ValidationError{Field: loc + ".id", Message: "component id is required"}
		}
		if strings.TrimSpace(c.Name) == "" {
			return &domerr.ValidationError{Field: loc + ".name", Message: "component name is required"}
		}
		if strings.TrimSpace(c.Type) == "" {
			return &domerr.ValidationError{Field: loc + ".type", Message: "component type is required"}
		}

		componentIDs[c.ID] = struct{}{}
	}

	for i, rel := range d.Relationships {
		loc := fmt.Sprintf("relationships[%d]", i)

		if _, ok := componentIDs[rel.Source]; !ok {
			return &domerr.ValidationError{
				Field:   loc + ".source",
				Message: fmt.Sprintf("source %q does not reference a known component", rel.Source),
			}
		}
		if _, ok := componentIDs[rel.Target]; !ok {
			return &domerr.ValidationError{
				Field:   loc + ".target",
				Message: fmt.Sprintf("target %q does not reference a known component", rel.Target),
			}
		}
		if strings.TrimSpace(rel.Kind) == "" {
			return &domerr.ValidationError{Field: loc + ".kind", Message: "relationship kind is required"}
		}
	}

	return nil
}
