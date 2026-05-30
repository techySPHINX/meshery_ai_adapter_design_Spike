package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllamaHTTPProviderHealthCheckOK(t *testing.T) {
	var gotPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"models":[{"name":"llama3"}]}`))
	}))
	defer srv.Close()

	p := NewOllamaHTTPProvider(srv.URL)
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Fatalf("unexpected health check error: %v", err)
	}
	if gotPath != "/api/tags" {
		t.Fatalf("expected /api/tags, got %q", gotPath)
	}
}

func TestOllamaHTTPProviderHealthCheckNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := NewOllamaHTTPProvider(srv.URL)
	if err := p.HealthCheck(context.Background()); err == nil {
		t.Fatal("expected health check error for non-200 response")
	}
}

func TestOllamaHTTPProviderGenerateRequiresModel(t *testing.T) {
	p := NewOllamaHTTPProvider("http://example")
	_, err := p.Generate(context.Background(), PromptRequest{Prompt: "hello"})
	if err == nil {
		t.Fatal("expected error for missing model")
	}
}

func TestOllamaHTTPProviderGenerateHappyPath(t *testing.T) {
	var gotReq ollamaGenerateRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ollamaGenerateResponse{
			Response: `{"name":"ok","components":[],"relationships":[]}`,
		})
	}))
	defer srv.Close()

	p := NewOllamaHTTPProvider(srv.URL)
	resp, err := p.Generate(context.Background(), PromptRequest{Prompt: "build", Model: "llama3"})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}
	if strings.TrimSpace(resp.Raw) == "" {
		t.Fatal("expected a non-empty response body")
	}
	if gotReq.Model != "llama3" {
		t.Fatalf("expected model llama3, got %q", gotReq.Model)
	}
	if gotReq.Format != "json" {
		t.Fatalf("expected format json, got %q", gotReq.Format)
	}
}
