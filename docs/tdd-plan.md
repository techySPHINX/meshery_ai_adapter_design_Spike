# TDD plan

Tests were written before implementation in this order. Each entry shows the test name, the failure reason before implementation, and what it drove.

## 1. Credential store

`TestCredentialStoreReturnsRefById` — failed: store didn't exist  
`TestCredentialStoreMissReturnsNotFound` — failed: no lookup  
`TestCredentialRefHasNoRawSecret` — contract test: struct must not gain a secret field  

Drove: `CredentialStore`, `CredentialRef`, `Lookup` method.

## 2. Provider resolver

`TestResolverSelectsOllamaWhenLocalOnlyTrue` — drove the LocalOnly override logic  
`TestResolverSelectsCloudProviderWhenLocalOnlyFalse` — drove cloud path  
`TestResolverSelectsOllamaByProviderType` — drove ollama-by-type path  
`TestResolverReturnsErrorForUnknownProviderType` — drove default error case  
`TestResolverReturnsErrorWhenLocalOnlyButNoOllamaRegistered` — drove nil-provider guard  

Drove: `Resolver`, `Resolve` method, `NewResolver`.

## 3. Prompt builder

`TestPromptBuilderInjectsUserIntent` — simplest possible assertion  
`TestPromptBuilderInjectsSchemaContext` — drove SchemaContext inclusion  
`TestPromptBuilderDoesNotContainCredentialValues` — documents the contract  
`TestPromptBuilderWithCustomSchemaContext` — drove the SchemaContext parameter  
`TestPromptBuilderOutputIsNonEmpty` — edge case: empty intent  

Drove: `Builder`, `Build`, `SchemaContext`, `DefaultSchemaContext`.

## 4. Design parser

`TestParserAcceptsValidJSON` — baseline  
`TestParserRejectsMalformedJSON` — drove *ParseError type  
`TestParserStripsMarkdownFences` — drove `stripFences` helper  

Drove: `ParseDesignResponse`, `ParseError`.

## 5. Design validator

`TestValidatorAcceptsWellFormedDesign` — baseline  
`TestValidatorRejectsMissingDesignName` — drove name check  
`TestValidatorRejectsMissingComponentID` — drove component loop  
`TestValidatorRejectsMissingComponentType` — drove type check  
`TestValidatorRejectsUnknownRelationshipSource` — drove component ID index  
`TestValidatorRejectsUnknownRelationshipTarget` — drove target check  
`TestValidatorAcceptsValidRelationships` — regression: valid relationships must pass  

Drove: `ValidateDesign`, `ValidationError`.

## 6. Adapter end-to-end

`TestAdapterReturnsValidatedDesignFromMockCloudProvider` — happy path cloud  
`TestAdapterReturnsValidatedDesignFromMockOllamaProvider` — happy path ollama  
`TestAdapterFailsWhenCredentialNotFound` — drove credential check in adapter  
`TestAdapterFailsWhenProviderUnhealthy` — drove health check call  
`TestAdapterFailsWhenLocalProviderUnhealthy` — drove local provider health failure path  
`TestAdapterFailsOnInvalidProviderOutput` — drove validation integration  
`TestAdapterFailsOnMalformedJSON` — drove parse integration  
`TestAdapterFailsOnProviderTimeout` — drove context cancellation handling in providers  
`TestAdapterLocalOnlyRoutesToOllama` — integration of LocalOnly override  

Drove: `Adapter`, `GenerateDesign`, `CredentialLookup` interface.

---

## 7. Optional Ollama integration

`TestAdapterWithRealOllamaIntegration` — exercises a real Ollama HTTP call (env-gated)

Drove: a safe, optional integration path without coupling core tests to local model state.

---

## Test count: 33 tests, 1 optional integration (env-gated), 0 external dependencies by default
