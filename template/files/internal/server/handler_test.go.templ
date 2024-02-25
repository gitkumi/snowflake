package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexRoute(t *testing.T) {
	s := New()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHealthRoute(t *testing.T) {
	s := New()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	s.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want %v", w.Code, http.StatusOK)
	}

	expected := "{\"message\":\"OK\"}"
	if w.Body.String() != expected {
		t.Errorf("got %v want %v", w.Body.String(), expected)
	}
}
