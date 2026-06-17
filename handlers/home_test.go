package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Change to project root so template paths resolve
	os.Chdir("..")
	code := m.Run()
	os.Exit(code)
}

func TestHandleHome_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for POST request, got %d", rr.Code)
	}
}

func TestHandleHome_ValidGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleHome_WithPageParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=1", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleHome_CardDisplay(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?display=card", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleHome_BlockDisplay(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?display=block", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleHome_TableDisplay(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?display=table", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleHome_FilterByName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?q=queen", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleHome_FilterByMember(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?member=freddie", nil)
	rr := httptest.NewRecorder()

	HandleHome(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}
