package jetspotter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newJSONAPIServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ac":[],"total":0}`))
	}))
}

func newForbiddenAPIServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
}

func TestSelectBestAPIUsesFirstReachableURL(t *testing.T) {
	original := baseURL
	defer func() { baseURL = original }()

	blocked := newForbiddenAPIServer()
	defer blocked.Close()
	working := newJSONAPIServer()
	defer working.Close()

	SelectBestAPI([]string{blocked.URL, working.URL})

	if baseURL != working.URL {
		t.Fatalf("expected baseURL to be %q, got %q", working.URL, baseURL)
	}
}

func TestSelectBestAPIRespectsConfiguredOrder(t *testing.T) {
	original := baseURL
	defer func() { baseURL = original }()

	first := newJSONAPIServer()
	defer first.Close()
	second := newJSONAPIServer()
	defer second.Close()

	SelectBestAPI([]string{first.URL, second.URL})

	if baseURL != first.URL {
		t.Fatalf("expected baseURL to be the first reachable URL %q, got %q", first.URL, baseURL)
	}
}

func TestSelectBestAPIFallsBackToFirstURLWhenNoneReachable(t *testing.T) {
	original := baseURL
	defer func() { baseURL = original }()

	blocked := newForbiddenAPIServer()
	defer blocked.Close()

	SelectBestAPI([]string{blocked.URL, "http://127.0.0.1:0"})

	if baseURL != blocked.URL {
		t.Fatalf("expected baseURL to fall back to %q, got %q", blocked.URL, baseURL)
	}
}

func TestSelectBestAPIUsesDefaultsWhenListEmpty(t *testing.T) {
	originalBaseURL := baseURL
	originalDefaults := defaultAPIURLs
	defer func() {
		baseURL = originalBaseURL
		defaultAPIURLs = originalDefaults
	}()

	working := newJSONAPIServer()
	defer working.Close()
	defaultAPIURLs = []string{working.URL}

	SelectBestAPI(nil)

	if baseURL != working.URL {
		t.Fatalf("expected baseURL to be taken from defaultAPIURLs (%q), got %q", working.URL, baseURL)
	}
}
