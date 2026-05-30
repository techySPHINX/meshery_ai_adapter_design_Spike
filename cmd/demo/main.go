package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/adapter"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt"
	"github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider"
)

func main() {
	// Wire up dependencies manually — no DI framework needed at this scale.
	store := connection.NewCredentialStore()
	store.Register(connection.CredentialRef{ID: "cred-cloud-001", Type: "api_key"})

	cloud := provider.NewMockCloudProvider("cred-cloud-001")

	useRealOllama := strings.TrimSpace(os.Getenv("OLLAMA_HTTP")) == "1"
	ollamaModel := strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
	if ollamaModel == "" {
		ollamaModel = "llama3"
	}
	ollamaBaseURL := strings.TrimSpace(os.Getenv("OLLAMA_BASE_URL"))
	if ollamaBaseURL == "" {
		ollamaBaseURL = provider.DefaultOllamaBaseURL
	}

	var ollama provider.LLMProvider
	if useRealOllama {
		ollama = provider.NewOllamaHTTPProvider(ollamaBaseURL)
	} else {
		ollama = provider.NewMockOllamaProvider()
	}

	resolver := provider.NewResolver(cloud, ollama)
	pb := prompt.NewBuilder(prompt.DefaultSchemaContext)

	a := adapter.New(resolver, pb, store)

	// Run two scenarios back-to-back to show both providers work.
	scenarios := []struct {
		label  string
		intent string
		conn   connection.Connection
	}{
		{
			label:  "cloud provider — nginx + prometheus",
			intent: "Create nginx with prometheus monitoring",
			conn: connection.Connection{
				ProviderType: "cloud",
				CredentialID: "cred-cloud-001",
				Model:        "gpt-4o",
			},
		},
		{
			label:  "local Ollama — privacy-first, no credentials needed",
			intent: "Deploy a redis instance",
			conn: connection.Connection{
				ProviderType: "ollama",
				LocalOnly:    true,
				BaseURL:      ollamaBaseURL,
				Model:        ollamaModel,
			},
		},
	}

	for _, s := range scenarios {
		fmt.Printf("\n=== %s ===\n", s.label)
		fmt.Printf("intent: %q\n", s.intent)

		d, err := a.GenerateDesign(context.Background(), s.intent, s.conn)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(d, "", "  ")
		fmt.Printf("validated design:\n%s\n", out)
	}
}
