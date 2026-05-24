package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"storeSystem/internal/database/mocks"
)

func TestGetAllStocks_Success(t *testing.T) {

	h := &Handlers{
		stockStore: &mocks.MockStockStore{},
	}

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	rr := httptest.NewRecorder()

	h.GetAllStocks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetAllStocks_Empty(t *testing.T) {

	h := &Handlers{
		stockStore: &mocks.EmptyStockStoreMock{},
	}

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	rr := httptest.NewRecorder()

	h.GetAllStocks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestGetAllStocks_DBError(t *testing.T) {

	h := &Handlers{
		stockStore: &mocks.ErrorStockStoreMock{},
	}

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	rr := httptest.NewRecorder()

	h.GetAllStocks(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}
