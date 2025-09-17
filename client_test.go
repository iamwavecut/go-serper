package serper

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockLogger struct{}

func (m *mockLogger) Debug(string, ...any) {}
func (m *mockLogger) Info(string, ...any)  {}
func (m *mockLogger) Warn(string, ...any)  {}
func (m *mockLogger) Error(string, ...any) {}

func createMockServer(t *testing.T, expectedEndpoint string, response any, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expectedEndpoint != "" {
			if r.URL.Path != expectedEndpoint {
				t.Fatalf("unexpected endpoint: %s", r.URL.Path)
			}
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-API-KEY") == "" {
			t.Fatal("missing X-API-KEY header")
		}
		var req SearchRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Query == "" {
			t.Fatal("empty query")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestClient_Search_Success(t *testing.T) {
	server := createMockServer(t, "/search", SearchResponse{Results: []SearchResult{{Title: "ok", URL: "https://x"}}}, http.StatusOK)
	defer server.Close()
	client := NewClient("k", WithBaseURL(server.URL), WithLogger(&mockLogger{}))
	resp, err := client.Search(context.Background(), SearchRequest{Query: "q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Title != "ok" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestClient_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("fail"))
	}))
	defer server.Close()
	client := NewClient("k", WithBaseURL(server.URL), WithLogger(&mockLogger{}))
	_, err := client.Search(context.Background(), SearchRequest{Query: "q"})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "api error") {
		t.Fatalf("expected api error, got: %v", err)
	}
}

func TestClient_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()
	client := NewClient("k", WithBaseURL(server.URL), WithLogger(&mockLogger{}))
	_, err := client.Search(context.Background(), SearchRequest{Query: "q"})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "unmarshal") {
		t.Fatalf("expected unmarshal error, got: %v", err)
	}
}

func TestClient_Retry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(SearchResponse{Results: []SearchResult{{Title: "ok"}}})
	}))
	defer server.Close()
	client := NewClient("k", WithBaseURL(server.URL), WithRetryConfig(3, 5*time.Millisecond), WithLogger(&mockLogger{}))
	resp, err := client.Search(context.Background(), SearchRequest{Query: "q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 || attempts != 3 {
		t.Fatalf("unexpected resp/attempts: %+v, %d", resp, attempts)
	}
}

func TestClient_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(150 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client := NewClient("k", WithBaseURL(server.URL), WithTimeouts(50*time.Millisecond, 50*time.Millisecond))
	_, err := client.Search(context.Background(), SearchRequest{Query: "q"})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "timeout") && !strings.Contains(msg, "deadline") && !strings.Contains(msg, "context") {
		t.Fatalf("unexpected error: %v", err)
	}
}
