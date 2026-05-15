# Test Results 

This file captures the terminal outputs from the test runs as proof of execution.

## Environment

Go version:

```text
go version go1.26.3 windows/amd64
```

Ollama integration settings used for the full suite:

```text
OLLAMA_INTEGRATION=1
OLLAMA_MODEL=llama3.2:latest
OLLAMA_BASE_URL=http://localhost:11434
```

## Package Tests (one by one)

Command:

```text
go test -v ./internal/connection
```

Output:

```text
=== RUN   TestCredentialStoreReturnsRefById
--- PASS: TestCredentialStoreReturnsRefById (0.00s)
=== RUN   TestCredentialStoreMissReturnsNotFound
--- PASS: TestCredentialStoreMissReturnsNotFound (0.00s)
=== RUN   TestCredentialRefHasNoRawSecret
--- PASS: TestCredentialRefHasNoRawSecret (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection      1.210s
```

Command:

```text
go test -v ./internal/provider
```

Output:

```text
=== RUN   TestResolverSelectsOllamaWhenLocalOnlyTrue
--- PASS: TestResolverSelectsOllamaWhenLocalOnlyTrue (0.00s)
=== RUN   TestResolverSelectsCloudProviderWhenLocalOnlyFalse
--- PASS: TestResolverSelectsCloudProviderWhenLocalOnlyFalse (0.00s)
=== RUN   TestResolverSelectsOllamaByProviderType
--- PASS: TestResolverSelectsOllamaByProviderType (0.00s)
=== RUN   TestResolverReturnsErrorForUnknownProviderType
--- PASS: TestResolverReturnsErrorForUnknownProviderType (0.00s)
=== RUN   TestResolverReturnsErrorWhenLocalOnlyButNoOllamaRegistered
--- PASS: TestResolverReturnsErrorWhenLocalOnlyButNoOllamaRegistered (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider        1.439s
```

Command:

```text
go test -v ./internal/prompt
```

Output:

```text
=== RUN   TestPromptBuilderInjectsUserIntent
--- PASS: TestPromptBuilderInjectsUserIntent (0.00s)
=== RUN   TestPromptBuilderInjectsSchemaContext
--- PASS: TestPromptBuilderInjectsSchemaContext (0.00s)
=== RUN   TestPromptBuilderDoesNotContainCredentialValues
--- PASS: TestPromptBuilderDoesNotContainCredentialValues (0.00s)
=== RUN   TestPromptBuilderWithCustomSchemaContext
--- PASS: TestPromptBuilderWithCustomSchemaContext (0.00s)
=== RUN   TestPromptBuilderOutputIsNonEmpty
--- PASS: TestPromptBuilderOutputIsNonEmpty (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt  1.130s
```

Command:

```text
go test -v ./internal/design
```

Output:

```text
=== RUN   TestParserAcceptsValidJSON
--- PASS: TestParserAcceptsValidJSON (0.00s)
=== RUN   TestParserRejectsMalformedJSON
--- PASS: TestParserRejectsMalformedJSON (0.00s)
=== RUN   TestParserStripsMarkdownFences
--- PASS: TestParserStripsMarkdownFences (0.00s)
=== RUN   TestValidatorAcceptsWellFormedDesign
--- PASS: TestValidatorAcceptsWellFormedDesign (0.00s)
=== RUN   TestValidatorRejectsMissingDesignName
--- PASS: TestValidatorRejectsMissingDesignName (0.00s)
=== RUN   TestValidatorRejectsMissingComponentID
--- PASS: TestValidatorRejectsMissingComponentID (0.00s)
=== RUN   TestValidatorRejectsMissingComponentType
--- PASS: TestValidatorRejectsMissingComponentType (0.00s)
=== RUN   TestValidatorRejectsUnknownRelationshipSource
--- PASS: TestValidatorRejectsUnknownRelationshipSource (0.00s)
=== RUN   TestValidatorRejectsUnknownRelationshipTarget
--- PASS: TestValidatorRejectsUnknownRelationshipTarget (0.00s)
=== RUN   TestValidatorAcceptsValidRelationships
--- PASS: TestValidatorAcceptsValidRelationships (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/design  1.405s
```

Command:

```text
go test -v ./internal/adapter
```

Output:

```text
=== RUN   TestAdapterReturnsValidatedDesignFromMockCloudProvider
--- PASS: TestAdapterReturnsValidatedDesignFromMockCloudProvider (0.00s)
=== RUN   TestAdapterReturnsValidatedDesignFromMockOllamaProvider
--- PASS: TestAdapterReturnsValidatedDesignFromMockOllamaProvider (0.00s)
=== RUN   TestAdapterFailsWhenCredentialNotFound
--- PASS: TestAdapterFailsWhenCredentialNotFound (0.00s)
=== RUN   TestAdapterFailsWhenProviderUnhealthy
--- PASS: TestAdapterFailsWhenProviderUnhealthy (0.00s)
=== RUN   TestAdapterFailsWhenLocalProviderUnhealthy
--- PASS: TestAdapterFailsWhenLocalProviderUnhealthy (0.00s)
=== RUN   TestAdapterFailsOnProviderTimeout
--- PASS: TestAdapterFailsOnProviderTimeout (0.00s)
=== RUN   TestAdapterFailsOnInvalidProviderOutput
--- PASS: TestAdapterFailsOnInvalidProviderOutput (0.00s)
=== RUN   TestAdapterFailsOnMalformedJSON
--- PASS: TestAdapterFailsOnMalformedJSON (0.00s)
=== RUN   TestAdapterLocalOnlyRoutesToOllama
--- PASS: TestAdapterLocalOnlyRoutesToOllama (0.00s)
=== RUN   TestAdapterWithRealOllamaIntegration
    ollama_integration_test.go:146: set OLLAMA_INTEGRATION=1 to enable real Ollama integration testing
--- SKIP: TestAdapterWithRealOllamaIntegration (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/adapter 1.975s
```

## Full Suite (with real Ollama integration)

Command:

```text
go test -v ./...
```

Output:

```text
?       github.com/techysphinx/meshery-ai-adapter-design-spike/cmd/demo [no test files]
=== RUN   TestAdapterReturnsValidatedDesignFromMockCloudProvider
--- PASS: TestAdapterReturnsValidatedDesignFromMockCloudProvider (0.00s)
=== RUN   TestAdapterReturnsValidatedDesignFromMockOllamaProvider
--- PASS: TestAdapterReturnsValidatedDesignFromMockOllamaProvider (0.00s)
=== RUN   TestAdapterFailsWhenCredentialNotFound
--- PASS: TestAdapterFailsWhenCredentialNotFound (0.00s)
=== RUN   TestAdapterFailsWhenProviderUnhealthy
--- PASS: TestAdapterFailsWhenProviderUnhealthy (0.00s)
=== RUN   TestAdapterFailsWhenLocalProviderUnhealthy
--- PASS: TestAdapterFailsWhenLocalProviderUnhealthy (0.00s)
=== RUN   TestAdapterFailsOnProviderTimeout
--- PASS: TestAdapterFailsOnProviderTimeout (0.00s)
=== RUN   TestAdapterFailsOnInvalidProviderOutput
--- PASS: TestAdapterFailsOnInvalidProviderOutput (0.00s)
=== RUN   TestAdapterFailsOnMalformedJSON
--- PASS: TestAdapterFailsOnMalformedJSON (0.00s)
=== RUN   TestAdapterLocalOnlyRoutesToOllama
--- PASS: TestAdapterLocalOnlyRoutesToOllama (0.00s)
=== RUN   TestAdapterWithRealOllamaIntegration
--- PASS: TestAdapterWithRealOllamaIntegration (5.31s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/adapter 6.124s
=== RUN   TestCredentialStoreReturnsRefById
--- PASS: TestCredentialStoreReturnsRefById (0.00s)
=== RUN   TestCredentialStoreMissReturnsNotFound
--- PASS: TestCredentialStoreMissReturnsNotFound (0.00s)
=== RUN   TestCredentialRefHasNoRawSecret
--- PASS: TestCredentialRefHasNoRawSecret (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/connection      (cached)
=== RUN   TestParserAcceptsValidJSON
--- PASS: TestParserAcceptsValidJSON (0.00s)
=== RUN   TestParserRejectsMalformedJSON
--- PASS: TestParserRejectsMalformedJSON (0.00s)
=== RUN   TestParserStripsMarkdownFences
--- PASS: TestParserStripsMarkdownFences (0.00s)
=== RUN   TestValidatorAcceptsWellFormedDesign
--- PASS: TestValidatorAcceptsWellFormedDesign (0.00s)
=== RUN   TestValidatorRejectsMissingDesignName
--- PASS: TestValidatorRejectsMissingDesignName (0.00s)
=== RUN   TestValidatorRejectsMissingComponentID
--- PASS: TestValidatorRejectsMissingComponentID (0.00s)
=== RUN   TestValidatorRejectsMissingComponentType
--- PASS: TestValidatorRejectsMissingComponentType (0.00s)
=== RUN   TestValidatorRejectsUnknownRelationshipSource
--- PASS: TestValidatorRejectsUnknownRelationshipSource (0.00s)
=== RUN   TestValidatorRejectsUnknownRelationshipTarget
--- PASS: TestValidatorRejectsUnknownRelationshipTarget (0.00s)
=== RUN   TestValidatorAcceptsValidRelationships
--- PASS: TestValidatorAcceptsValidRelationships (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/design  (cached)
?       github.com/techysphinx/meshery-ai-adapter-design-spike/internal/errors  [no test files]
=== RUN   TestPromptBuilderInjectsUserIntent
--- PASS: TestPromptBuilderInjectsUserIntent (0.00s)
=== RUN   TestPromptBuilderInjectsSchemaContext
--- PASS: TestPromptBuilderInjectsSchemaContext (0.00s)
=== RUN   TestPromptBuilderDoesNotContainCredentialValues
--- PASS: TestPromptBuilderDoesNotContainCredentialValues (0.00s)
=== RUN   TestPromptBuilderWithCustomSchemaContext
--- PASS: TestPromptBuilderWithCustomSchemaContext (0.00s)
=== RUN   TestPromptBuilderOutputIsNonEmpty
--- PASS: TestPromptBuilderOutputIsNonEmpty (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/prompt  (cached)
=== RUN   TestResolverSelectsOllamaWhenLocalOnlyTrue
--- PASS: TestResolverSelectsOllamaWhenLocalOnlyTrue (0.00s)
=== RUN   TestResolverSelectsCloudProviderWhenLocalOnlyFalse
--- PASS: TestResolverSelectsCloudProviderWhenLocalOnlyFalse (0.00s)
=== RUN   TestResolverSelectsOllamaByProviderType
--- PASS: TestResolverSelectsOllamaByProviderType (0.00s)
=== RUN   TestResolverReturnsErrorForUnknownProviderType
--- PASS: TestResolverReturnsErrorForUnknownProviderType (0.00s)
=== RUN   TestResolverReturnsErrorWhenLocalOnlyButNoOllamaRegistered
--- PASS: TestResolverReturnsErrorWhenLocalOnlyButNoOllamaRegistered (0.00s)
PASS
ok      github.com/techysphinx/meshery-ai-adapter-design-spike/internal/provider        (cached)
```
