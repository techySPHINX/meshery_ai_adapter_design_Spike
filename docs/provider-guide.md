# Provider configuration guide

This spike models adapter connections with the Connection type and a resolver that chooses a provider implementation by name. The examples below show how to configure local and cloud-style providers.

## Connection fields

- ProviderType: provider key registered in the resolver (examples: cloud, ollama).
- LocalOnly: when true, the resolver forces a local provider (ollama) even if ProviderType says otherwise.
- BaseURL: base URL for providers that require it (Ollama uses http://localhost:11434 by default).
- Model: model identifier passed to the provider (for example: llama3).
- CredentialID: reference to a stored credential (never the secret value itself).

## Local Ollama (privacy-first)

1) Install and run Ollama locally.
2) Pull a model (example: ollama pull llama3).
3) Configure the connection:

```go
conn := connection.Connection{
	ProviderType: "ollama",
	LocalOnly:    true,
	BaseURL:      "http://localhost:11434",
	Model:        "llama3",
}
```

To run the demo against a real Ollama instance:

```bash
set OLLAMA_HTTP=1
set OLLAMA_BASE_URL=http://localhost:11434
set OLLAMA_MODEL=llama3
```

## Cloud provider (mock in this spike)

The cloud provider in this repository is a mock. It demonstrates the credential boundary and provider selection, but it does not call a real API.

```go
conn := connection.Connection{
	ProviderType: "cloud",
	CredentialID: "cred-cloud-001",
	Model:        "gpt-4o",
}
```

In a production integration, ProviderType would map to real implementations (OpenAI, Vertex AI, etc.) registered in the resolver.

## Swapping providers

The resolver supports a registry so you can replace providers without changing the adapter.

```go
providers := map[string]provider.LLMProvider{
	"cloud":  provider.NewMockCloudProvider("cred-cloud-001"),
	"ollama": provider.NewOllamaHTTPProvider("http://localhost:11434"),
}
resolver := provider.NewResolverWithProviders(providers, "ollama")
```
