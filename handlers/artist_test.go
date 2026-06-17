package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandleArtistPage_InvalidID(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/artist/{id}", HandleArtistPage)

	req := httptest.NewRequest(http.MethodGet, "/artist/invalid", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid ID, got %d", rr.Code)
	}
}

func TestHandleArtistPage_NegativeID(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/artist/{id}", HandleArtistPage)

	req := httptest.NewRequest(http.MethodGet, "/artist/-1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for negative ID, got %d", rr.Code)
	}
}

func TestHandleArtistPage_ZeroID(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/artist/{id}", HandleArtistPage)

	req := httptest.NewRequest(http.MethodGet, "/artist/0", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for zero ID, got %d", rr.Code)
	}
}

func TestNotFoundHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/artist/{id}", HandleArtistPage)
	r.NotFoundHandler = http.HandlerFunc(HandleNotFound)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}
