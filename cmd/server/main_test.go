package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
