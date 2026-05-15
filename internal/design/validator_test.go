package design_test

import (
	"strings"
	"testing"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/design"
	domerr "github.com/techysphinx/meshery-ai-adapter-design-spike/internal/errors"
)

// ---- parser tests ----

func TestParserAcceptsValidJSON(t *testing.T) {
	raw := `{"name":"my-design","components":[{"id":"c1","name":"nginx","type":"Deployment"}],"relationships":[]}`

	d, err := design.ParseDesignResponse(raw)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if d.Name != "my-design" {
		t.Errorf("got name %q, want %q", d.Name, "my-design")
	}
}

func TestParserRejectsMalformedJSON(t *testing.T) {
	raw := `{"name": "broken", "components": [`

	_, err := design.ParseDesignResponse(raw)
	if err == nil {
		t.Fatal("expected a parse error for malformed JSON, got nil")
	}

	var pe *domerr.ParseError
	if !errorAs(err, &pe) {
		t.Errorf("expected *ParseError, got %T: %v", err, err)
	}
}

func TestParserStripsMarkdownFences(t *testing.T) {
	raw := "```json\n{\"name\":\"fenced\",\"components\":[{\"id\":\"c1\",\"name\":\"redis\",\"type\":\"StatefulSet\"}],\"relationships\":[]}\n```"

	d, err := design.ParseDesignResponse(raw)
	if err != nil {
		t.Fatalf("should have stripped fences and parsed: %v", err)
	}
	if d.Name != "fenced" {
		t.Errorf("got name %q, want %q", d.Name, "fenced")
	}
}

// ---- validator tests ----

func TestValidatorAcceptsWellFormedDesign(t *testing.T) {
	d := &design.Design{
		Name: "ok",
		Components: []design.Component{
			{ID: "c1", Name: "nginx", Type: "Deployment"},
		},
		Relationships: []design.Relationship{},
	}

	if err := design.ValidateDesign(d); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidatorRejectsMissingDesignName(t *testing.T) {
	d := &design.Design{
		Components: []design.Component{
			{ID: "c1", Name: "nginx", Type: "Deployment"},
		},
	}

	err := design.ValidateDesign(d)
	assertValidationError(t, err, "name")
}

func TestValidatorRejectsMissingComponentID(t *testing.T) {
	d := &design.Design{
		Name: "test",
		Components: []design.Component{
			{Name: "nginx", Type: "Deployment"}, // missing ID
		},
	}

	err := design.ValidateDesign(d)
	assertValidationError(t, err, "id")
}

func TestValidatorRejectsMissingComponentType(t *testing.T) {
	d := &design.Design{
		Name: "test",
		Components: []design.Component{
			{ID: "c1", Name: "nginx"}, // missing Type
		},
	}

	err := design.ValidateDesign(d)
	assertValidationError(t, err, "type")
}

func TestValidatorRejectsUnknownRelationshipSource(t *testing.T) {
	d := &design.Design{
		Name: "test",
		Components: []design.Component{
			{ID: "c1", Name: "nginx", Type: "Deployment"},
		},
		Relationships: []design.Relationship{
			{Source: "c-does-not-exist", Target: "c1", Kind: "depends-on"},
		},
	}

	err := design.ValidateDesign(d)
	assertValidationError(t, err, "source")
}

func TestValidatorRejectsUnknownRelationshipTarget(t *testing.T) {
	d := &design.Design{
		Name: "test",
		Components: []design.Component{
			{ID: "c1", Name: "nginx", Type: "Deployment"},
		},
		Relationships: []design.Relationship{
			{Source: "c1", Target: "c-does-not-exist", Kind: "depends-on"},
		},
	}

	err := design.ValidateDesign(d)
	assertValidationError(t, err, "target")
}

func TestValidatorAcceptsValidRelationships(t *testing.T) {
	d := &design.Design{
		Name: "full-design",
		Components: []design.Component{
			{ID: "c1", Name: "nginx", Type: "Deployment"},
			{ID: "c2", Name: "prometheus", Type: "Service"},
		},
		Relationships: []design.Relationship{
			{Source: "c1", Target: "c2", Kind: "observes"},
		},
	}

	if err := design.ValidateDesign(d); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
}

// ---- helpers ----

func assertValidationError(t *testing.T, err error, fieldSubstring string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	var ve *domerr.ValidationError
	if !errorAs(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	if !strings.Contains(ve.Field, fieldSubstring) {
		t.Errorf("expected field to contain %q, got %q", fieldSubstring, ve.Field)
	}
}

// errorAs is a minimal stand-in so we don't need errors.As from stdlib here.
// We do use it, just aliasing to make the helper readable.
func errorAs(err error, target interface{}) bool {
	switch t := target.(type) {
	case **domerr.ValidationError:
		if ve, ok := err.(*domerr.ValidationError); ok {
			*t = ve
			return true
		}
	case **domerr.ParseError:
		if pe, ok := err.(*domerr.ParseError); ok {
			*t = pe
			return true
		}
	}
	return false
}
