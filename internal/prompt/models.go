package prompt

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultSystemPrompt sets the guardrails for the LLM.
const DefaultSystemPrompt = "You are a Meshery design assistant. " +
	"Use the Meshery model schemas in this prompt to select valid component types and relationships. " +
	"You have hands-on experience with REST, GraphQL, and gRPC APIs. " +
	"When the intent implies API exposure or service-to-service communication, " +
	"choose protocol-appropriate component types (RESTService, GraphQLService, gRPCService) and relationship kinds. " +
	"Do not invent fields beyond the schema rules; if protocol detail is needed, encode it in component type or name. " +
	"The model schemas provided are the relevant subset for this request."

// ModelDefinition is a compact Meshery model schema for prompt injection.
// It is intentionally terse to fit in the context window.
type ModelDefinition struct {
	Name          string
	ComponentType string
	Summary       string
	Relationships []string
	Protocols     []string
	Keywords      []string
	Core          bool
}

// Render returns a concise, single-line schema description.
func (m ModelDefinition) Render() string {
	line := fmt.Sprintf("- %s [type: %s]", m.Name, m.ComponentType)

	var details []string
	if strings.TrimSpace(m.Summary) != "" {
		details = append(details, m.Summary)
	}
	if len(m.Protocols) > 0 {
		details = append(details, "protocols: "+strings.Join(m.Protocols, ", "))
	}
	if len(m.Relationships) > 0 {
		details = append(details, "relationships: "+strings.Join(m.Relationships, ", "))
	}
	if len(details) > 0 {
		line += ": " + strings.Join(details, "; ")
	}

	return line + "\n"
}

// ModelCatalog contains all Meshery model definitions known to the prompt layer.
type ModelCatalog struct {
	Models []ModelDefinition
}

// ContextWindow bounds how many model definitions to include.
type ContextWindow struct {
	MaxChars  int
	MaxModels int
}

const (
	defaultMaxModelChars = 1800
	defaultMaxModels     = 6
)

// DefaultContextWindow keeps the prompt compact while still providing signal.
var DefaultContextWindow = ContextWindow{MaxChars: defaultMaxModelChars, MaxModels: defaultMaxModels}

// DefaultModelCatalog provides a concise Meshery model index for the spike.
func DefaultModelCatalog() ModelCatalog {
	return ModelCatalog{
		Models: []ModelDefinition{
			{
				Name:          "Kubernetes Deployment",
				ComponentType: "Deployment",
				Summary:       "Runs stateless workloads as pods",
				Relationships: []string{"exposes", "depends-on", "observes"},
				Keywords:      []string{"deploy", "deployment", "pod", "workload", "nginx", "app"},
				Core:          true,
			},
			{
				Name:          "Kubernetes Service",
				ComponentType: "Service",
				Summary:       "Stable network endpoint for workloads",
				Relationships: []string{"exposes", "routes-to"},
				Keywords:      []string{"service", "endpoint", "expose", "load"},
				Core:          true,
			},
			{
				Name:          "Ingress",
				ComponentType: "Ingress",
				Summary:       "HTTP routing into services",
				Relationships: []string{"routes-to"},
				Keywords:      []string{"ingress", "http", "https", "gateway"},
			},
			{
				Name:          "REST API Service",
				ComponentType: "RESTService",
				Summary:       "HTTP/JSON API endpoint",
				Relationships: []string{"exposes", "depends-on"},
				Protocols:     []string{"REST"},
				Keywords:      []string{"rest", "http", "api", "openapi"},
			},
			{
				Name:          "GraphQL API Service",
				ComponentType: "GraphQLService",
				Summary:       "Schema-driven GraphQL API endpoint",
				Relationships: []string{"exposes", "depends-on"},
				Protocols:     []string{"GraphQL"},
				Keywords:      []string{"graphql", "apollo", "schema"},
			},
			{
				Name:          "gRPC Service",
				ComponentType: "gRPCService",
				Summary:       "Protobuf-based RPC endpoint",
				Relationships: []string{"exposes", "depends-on"},
				Protocols:     []string{"gRPC"},
				Keywords:      []string{"grpc", "rpc", "protobuf"},
			},
			{
				Name:          "Prometheus",
				ComponentType: "Prometheus",
				Summary:       "Scrapes and stores metrics",
				Relationships: []string{"observes"},
				Keywords:      []string{"prometheus", "metrics", "monitor"},
			},
			{
				Name:          "PostgreSQL",
				ComponentType: "PostgreSQL",
				Summary:       "Relational database for stateful workloads",
				Relationships: []string{"depends-on"},
				Keywords:      []string{"postgres", "postgresql", "database", "sql"},
			},
		},
	}
}

type scoredModel struct {
	Model ModelDefinition
	Score int
	Hits  int
}

// Select returns the most relevant models for the intent within the context window.
func (c ModelCatalog) Select(intent string, window ContextWindow) []ModelDefinition {
	if len(c.Models) == 0 {
		return nil
	}

	window = normalizeWindow(window)
	intentLower := strings.ToLower(intent)
	keywordHit := false

	scored := make([]scoredModel, 0, len(c.Models))
	for _, m := range c.Models {
		hits := scoreKeywords(intentLower, m.Keywords)
		if hits > 0 {
			keywordHit = true
		}
		score := hits * 10
		if m.Core {
			score++
		}

		scored = append(scored, scoredModel{Model: m, Score: score, Hits: hits})
	}

	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].Model.Name < scored[j].Model.Name
		}
		return scored[i].Score > scored[j].Score
	})

	var selected []ModelDefinition
	usedChars := 0
	for _, entry := range scored {
		if len(selected) >= window.MaxModels {
			break
		}
		if keywordHit {
			if entry.Hits == 0 && !entry.Model.Core {
				continue
			}
		} else if !entry.Model.Core {
			continue
		}

		line := entry.Model.Render()
		if window.MaxChars > 0 && usedChars+len(line) > window.MaxChars {
			if len(selected) == 0 {
				selected = append(selected, entry.Model)
			}
			break
		}

		selected = append(selected, entry.Model)
		usedChars += len(line)
	}

	return selected
}

func normalizeWindow(window ContextWindow) ContextWindow {
	if window.MaxChars <= 0 {
		window.MaxChars = DefaultContextWindow.MaxChars
	}
	if window.MaxModels <= 0 {
		window.MaxModels = DefaultContextWindow.MaxModels
	}
	return window
}

func scoreKeywords(intent string, keywords []string) int {
	if strings.TrimSpace(intent) == "" {
		return 0
	}

	score := 0
	for _, keyword := range keywords {
		kw := strings.ToLower(strings.TrimSpace(keyword))
		if kw == "" {
			continue
		}
		if strings.Contains(intent, kw) {
			score++
		}
	}
	return score
}
