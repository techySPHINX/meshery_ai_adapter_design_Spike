# Architecture

## System flow

```mermaid
flowchart LR
	classDef intent fill:#FFF4CE,stroke:#B08800,color:#3B2E00;
	classDef core fill:#E7F3FF,stroke:#2B6CB0,color:#1A365D;
	classDef provider fill:#E6FFFA,stroke:#2C7A7B,color:#234E52;
	classDef validate fill:#FDEBC8,stroke:#DD6B20,color:#7B341E;

	A[User Intent]:::intent --> B[Adapter Orchestrator]:::core
	B --> C[Prompt Builder]:::core
	C --> D[Provider Resolver]:::core
	D --> E[Connection Config]:::core
	E --> F[Credential Ref]:::core
	D --> G{Provider Type}:::provider
	G --> H[Mock Cloud Provider]:::provider
	G --> I[Mock Ollama Provider]:::provider
	H --> J[Untrusted LLM Output]:::validate
	I --> J
	J --> K[Parser]:::validate
	K --> L[Schema Validator]:::validate
	L --> M[Validated Design Result]:::core
	L --> N[Safe Error]:::validate
```

## Component responsibilities

### Adapter (internal/adapter)

Orchestrates the pipeline. Knows nothing about:
- Provider HTTP details
- Credential secret values
- JSON parsing internals
- Validation rules

If you add logging, tracing, or retry logic, this is the only place it goes.

### Provider (internal/provider)

Defines the `LLMProvider` interface and the two mocks. The interface is three methods:

```go
Name()        string
Generate()    (*ProviderResponse, error)
HealthCheck() error
```

`ProviderResponse.Raw` is named to remind callers it is untrusted and must be parsed before use. The mocks honor context cancellation so timeout paths are testable.

### Connection and Credential (internal/connection)

`Connection` holds everything needed to route a request: provider type, base URL, model, and a `CredentialID` reference. It never holds a secret value.

`CredentialRef` is a pointer (ID + type). `CredentialStore.Lookup` confirms existence without surfacing the actual key. In a real Meshery implementation this would call the credentials API.

### Prompt (internal/prompt)

`Builder.Build(intent)` assembles the final prompt string from:
1. User intent (free text)
2. Schema rules (typed struct, not free text)

The builder has no knowledge of credentials, provider types, or network addresses. It cannot accidentally leak secrets because they are never passed to it.

### Design (internal/design)

Two steps, always in order:

1. `ParseDesignResponse(raw)` - JSON unmarshal with markdown fence stripping. Returns `*ParseError` on failure.
2. `ValidateDesign(d)` - field-level checks plus referential integrity on relationships. Returns `*ValidationError` on failure.

A `Design` that passes both is considered safe to return to the caller.

### Errors (internal/errors)

Four typed errors:

| Type | When |
|---|---|
| `ProviderError` | provider health check or generate call fails |
| `ValidationError` | a required field is missing or a relationship reference is invalid |
| `ParseError` | LLM output is not valid JSON |
| `CredentialError` | credential ID not found in the store |

All errors wrap their cause so callers can use `errors.As` to distinguish them.

## Provider resolver logic

```
LocalOnly == true  ->  always Ollama
ProviderType == "ollama"  ->  Ollama
ProviderType == "cloud"   ->  Cloud
anything else             ->  error
```

`LocalOnly` takes precedence over `ProviderType` so that user privacy preferences cannot be overridden by an incorrect config field.

## Sequence (single request)

```mermaid
sequenceDiagram
	participant User
	participant Adapter
	participant Prompt
	participant Resolver
	participant Provider
	participant Parser
	participant Validator

	User->>Adapter: intent + connection
	Adapter->>Prompt: Build(intent)
	Adapter->>Resolver: Resolve(connection)
	Resolver-->>Adapter: provider
	Adapter->>Provider: HealthCheck(ctx)
	Adapter->>Provider: Generate(prompt, model)
	Provider-->>Adapter: raw output
	Adapter->>Parser: Parse(raw)
	Parser-->>Adapter: design
	Adapter->>Validator: Validate(design)
	Validator-->>Adapter: ok or error
	Adapter-->>User: validated design or safe error
```
