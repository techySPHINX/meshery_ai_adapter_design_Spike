package design

// Design is the validated output of the AI adapter pipeline.
// It mirrors a simplified Meshery design: a named unit containing
// components and the relationships between them.
type Design struct {
	Name          string         `json:"name"`
	Components    []Component    `json:"components"`
	Relationships []Relationship `json:"relationships"`
}

// Component maps to a Meshery component — a deployable unit like a
// Kubernetes Deployment, Service, or ConfigMap.
type Component struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Relationship captures a directed edge between two components.
// Kind describes the nature of the connection (observes, depends-on, etc.).
type Relationship struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
}
