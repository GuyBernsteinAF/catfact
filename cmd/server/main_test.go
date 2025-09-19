package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPhaseFourAPI_ValidRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/cat-facts?name=TestUser&amount=2", nil)
	w := httptest.NewRecorder()

	phaseFourAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response SuccessResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Facts) != 2 {
		t.Errorf("expected 2 facts, got %d", len(response.Facts))
	}
}

func TestValidate_MissingName(t *testing.T) {
	w := httptest.NewRecorder()

	valid := validate(w, "1", "")

	if valid {
		t.Error("expected validation to fail for missing name")
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
