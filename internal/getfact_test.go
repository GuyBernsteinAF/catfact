package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ============================================
// Helper Functions
// ============================================

// helper to spin up a mock upstream server with status + body
func makeMockUpstream(status int, body string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
	return httptest.NewServer(h)
}

// helper to create a mock server that returns different facts
func makeMockUpstreamMultiple(facts []string) *httptest.Server {
	counter := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if counter < len(facts) {
			fact := facts[counter]
			payload := fmt.Sprintf(`{"fact":"%s", "length": %d}`, fact, len(fact))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
			counter++
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"no more facts"}`))
		}
	})
	return httptest.NewServer(h)
}

// ============================================
// GetFact Tests - Edge Cases
// ============================================

func TestGetFact_ServerError_ReturnsEmpty(t *testing.T) {
	srv := makeMockUpstream(http.StatusInternalServerError, `{"error":"boom"}`)
	defer srv.Close()

	got := GetFact(srv.URL + "/fact")
	want := ""
	if got != want {
		t.Fatalf("server 500: expected %q, got %q", want, got)
	}
}

func TestGetFact_InvalidJSON_ReturnsEmpty(t *testing.T) {
	// broken JSON: missing closing brace + wrong type for "fact"
	srv := makeMockUpstream(http.StatusOK, `{"fact": 123, "length": 3`)
	defer srv.Close()

	got := GetFact(srv.URL + "/fact")
	want := ""
	if got != want {
		t.Fatalf("invalid JSON: expected %q, got %q", want, got)
	}
}

func TestGetFact_EmptyFact_ReturnsEmpty(t *testing.T) {
	srv := makeMockUpstream(http.StatusOK, `{"fact":"", "length":0}`)
	defer srv.Close()

	got := GetFact(srv.URL + "/fact")
	want := ""
	if got != want {
		t.Fatalf("empty fact: expected %q, got %q", want, got)
	}
}

func TestGetFact_Success_AlternativeFact(t *testing.T) {
	const fact = "Cats can rotate their ears 180 degrees."
	payload := `{"fact":"` + strings.ReplaceAll(fact, `"`, `\"`) + `", "length": 43}`
	srv := makeMockUpstream(http.StatusOK, payload)
	defer srv.Close()

	got := GetFact(srv.URL + "/fact")
	if got != fact {
		t.Fatalf("success alt: expected %q, got %q", fact, got)
	}
}

func TestGetFact_MissingFactField_ReturnsEmpty(t *testing.T) {
	srv := makeMockUpstream(http.StatusOK, `{"notfact":"something", "length":9}`)
	defer srv.Close()

	got := GetFact(srv.URL + "/fact")
	want := ""
	if got != want {
		t.Fatalf("missing fact field: expected %q, got %q", want, got)
	}
}

func TestGetFact_FactNotString_ReturnsEmpty(t *testing.T) {
	srv := makeMockUpstream(http.StatusOK, `{"fact":123, "length":3}`)
	defer srv.Close()

	got := GetFact(srv.URL + "/fact")
	want := ""
	if got != want {
		t.Fatalf("fact not string: expected %q, got %q", want, got)
	}
}

// ============================================
// PhaseFour Function Tests
// ============================================

func TestPhaseFourWithURL_SingleFact(t *testing.T) {
	facts := []string{"Cats have 32 muscles in each ear."}
	srv := makeMockUpstreamMultiple(facts)
	defer srv.Close()

	// Use the helper function to test with custom URL
	result := PhaseFourWithURL(1, srv.URL+"/fact")

	if len(result) != 1 {
		t.Fatalf("expected 1 fact, got %d", len(result))
	}
	if result[0] != facts[0] {
		t.Fatalf("expected %q, got %q", facts[0], result[0])
	}
}
