package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"storeSystem/internal/database/mocks"

	"github.com/go-chi/chi/v5"
)

func setupCounterpartyHandler() *Handlers {
	return &Handlers{
		counterpartyStore: &mocks.MockCounterpartyStore{},
	}
}

func addCounterpartyURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)
}

func TestGetAllCounterparties_Success(t *testing.T) {

	h := setupCounterpartyHandler()

	req := httptest.NewRequest(http.MethodGet, "/counterparties", nil)
	rr := httptest.NewRecorder()

	h.GetAllCounterparties(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetCounterpartyByID_Success(t *testing.T) {

	h := setupCounterpartyHandler()

	req := httptest.NewRequest(http.MethodGet, "/counterparties/1", nil)
	req = addCounterpartyURLParam(req, "id", "1")

	rr := httptest.NewRecorder()

	h.GetCounterpartyByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetCounterpartyByID_InvalidID(t *testing.T) {

	h := setupCounterpartyHandler()

	req := httptest.NewRequest(http.MethodGet, "/counterparties/abc", nil)
	req = addCounterpartyURLParam(req, "id", "abc")

	rr := httptest.NewRecorder()

	h.GetCounterpartyByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestGetCounterpartyByID_NotFound(t *testing.T) {

	h := setupCounterpartyHandler()

	req := httptest.NewRequest(http.MethodGet, "/counterparties/999", nil)
	req = addCounterpartyURLParam(req, "id", "999")

	rr := httptest.NewRecorder()

	h.GetCounterpartyByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestCreateCounterparty_Success(t *testing.T) {

	h := setupCounterpartyHandler()

	body := `{
		"full_name":"test counterparty"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/counterparties",
		bytes.NewBufferString(body),
	)

	rr := httptest.NewRecorder()

	h.CreateCounterparty(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}

func TestCreateCounterparty_InvalidJSON(t *testing.T) {

	h := setupCounterpartyHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/counterparties",
		strings.NewReader("{invalid json"),
	)

	rr := httptest.NewRecorder()

	h.CreateCounterparty(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestCreateCounterparty_DBError(t *testing.T) {

	h := setupCounterpartyHandler()

	body := `{
		"full_name":"db-error"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/counterparties",
		bytes.NewBufferString(body),
	)

	rr := httptest.NewRecorder()

	h.CreateCounterparty(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}
