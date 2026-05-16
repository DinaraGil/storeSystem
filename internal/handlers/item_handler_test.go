package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"storeSystem/internal/database/mocks"
	"storeSystem/internal/models"

	"github.com/go-chi/chi/v5"
)

func setupHandler() *Handlers {
	return &Handlers{
		itemStore: &mocks.MockItemStore{},
	}
}

func addURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)
}

func TestGetAllItems_Success(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rr := httptest.NewRecorder()

	h.GetAllItems(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetItemByID_Success(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	req = addURLParam(req, "id", "1")

	rr := httptest.NewRecorder()

	h.GetItemByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetItemByID_InvalidID(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/items/abc", nil)
	req = addURLParam(req, "id", "abc")

	rr := httptest.NewRecorder()

	h.GetItemByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestGetItemByID_NotFound(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/items/999", nil)
	req = addURLParam(req, "id", "999")

	rr := httptest.NewRecorder()

	h.GetItemByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestCreateItem_Success(t *testing.T) {
	h := setupHandler()

	body := `{
		"article":"test article"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/items",
		bytes.NewBufferString(body),
	)

	rr := httptest.NewRecorder()

	h.CreateItem(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}

func TestCreateItem_InvalidJSON(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/items",
		strings.NewReader("{invalid json"),
	)

	rr := httptest.NewRecorder()

	h.CreateItem(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestCreateItem_EmptyArticle(t *testing.T) {
	h := setupHandler()

	body := `{
		"article":""
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/items",
		bytes.NewBufferString(body),
	)

	rr := httptest.NewRecorder()

	h.CreateItem(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestCreateItem_DBError(t *testing.T) {
	h := setupHandler()

	body := `{
		"article":"db-error"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/items",
		bytes.NewBufferString(body),
	)

	rr := httptest.NewRecorder()

	h.CreateItem(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}

func TestUpdateItem_Success(t *testing.T) {
	h := setupHandler()

	article := "updated"

	body, _ := json.Marshal(models.UpdateItemInput{
		Article: &article,
	})

	req := httptest.NewRequest(
		http.MethodPut,
		"/items/1",
		bytes.NewBuffer(body),
	)

	req = addURLParam(req, "id", "1")

	rr := httptest.NewRecorder()

	h.UpdateItem(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestUpdateItem_InvalidID(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodPut,
		"/items/abc",
		nil,
	)

	req = addURLParam(req, "id", "abc")

	rr := httptest.NewRecorder()

	h.UpdateItem(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestUpdateItem_InvalidJSON(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodPut,
		"/items/1",
		strings.NewReader("{invalid"),
	)

	req = addURLParam(req, "id", "1")

	rr := httptest.NewRecorder()

	h.UpdateItem(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestUpdateItem_EmptyArticle(t *testing.T) {
	h := setupHandler()

	article := ""

	body, _ := json.Marshal(models.UpdateItemInput{
		Article: &article,
	})

	req := httptest.NewRequest(
		http.MethodPut,
		"/items/1",
		bytes.NewBuffer(body),
	)

	req = addURLParam(req, "id", "1")

	rr := httptest.NewRecorder()

	h.UpdateItem(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestUpdateItem_NotFound(t *testing.T) {
	h := setupHandler()

	article := "updated"

	body, _ := json.Marshal(models.UpdateItemInput{
		Article: &article,
	})

	req := httptest.NewRequest(
		http.MethodPut,
		"/items/999",
		bytes.NewBuffer(body),
	)

	req = addURLParam(req, "id", "999")

	rr := httptest.NewRecorder()

	h.UpdateItem(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rr.Code)
	}
}

func TestUpdateItem_DBError(t *testing.T) {
	h := setupHandler()

	article := "db-error"

	body, _ := json.Marshal(models.UpdateItemInput{
		Article: &article,
	})

	req := httptest.NewRequest(
		http.MethodPut,
		"/items/1",
		bytes.NewBuffer(body),
	)

	req = addURLParam(req, "id", "1")

	rr := httptest.NewRecorder()

	h.UpdateItem(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}
