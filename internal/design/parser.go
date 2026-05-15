package design

import (
	"encoding/json"
	"strings"

	domerr "github.com/techysphinx/meshery-ai-adapter-design-spike/internal/errors"
)

// ParseDesignResponse converts raw LLM output into a Design struct.
// It strips markdown code fences defensively (some models add them
// even when told not to) and fails fast on malformed JSON.
// The returned Design is not yet validated — call ValidateDesign next.
func ParseDesignResponse(raw string) (*Design, error) {
	cleaned := stripFences(raw)

	var d Design
	if err := json.Unmarshal([]byte(cleaned), &d); err != nil {
		return nil, &domerr.ParseError{Cause: err}
	}
	return &d, nil
}

// stripFences removes ```json ... ``` wrappers that some LLMs add.
func stripFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
