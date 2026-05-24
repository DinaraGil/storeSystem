package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"storeSystem/internal/database/mocks"

	"github.com/go-chi/chi/v5"
)

func setupDeliveryHandler() *Handlers {
	return &Handlers{
		deliveryStore:     &mocks.MockDeliveryStore{},
		deliveryListStore: &mocks.MockDeliveryListStore{},
	}
}

func addParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestGetAllDeliveries_Table(t *testing.T) {
	tests := []struct {
		name           string
		storeError     error
		expectedStatus int
	}{
		{
			name:           "success",
			storeError:     nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupDeliveryHandler()

			req := httptest.NewRequest(http.MethodGet, "/deliveries", nil)
			rr := httptest.NewRecorder()

			h.GetAllDeliveries(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetDeliveryByID_Table(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"success", "1", http.StatusOK},
		{"invalid id", "abc", http.StatusBadRequest},
		{"not found", "999", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupDeliveryHandler()

			req := httptest.NewRequest(http.MethodGet, "/deliveries/"+tt.id, nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.GetDeliveryByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCreateDelivery_Table(t *testing.T) {

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "success",
			body:           `{"status":"new"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           "{bad json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty status",
			body:           `{"status":""}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupDeliveryHandler()

			req := httptest.NewRequest(
				http.MethodPost,
				"/deliveries",
				bytes.NewBufferString(tt.body),
			)

			rr := httptest.NewRecorder()

			h.CreateDelivery(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetDeliveryLists_Table(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"success", "1", http.StatusOK},
		{"invalid id", "abc", http.StatusBadRequest},
		{"not found", "999", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupDeliveryHandler()

			req := httptest.NewRequest(http.MethodGet, "/deliveries/"+tt.id+"/lists", nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.GetDeliveryListsByDeliveryID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCompleteDelivery_Table(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"success", "1", http.StatusOK},
		{"invalid id", "abc", http.StatusBadRequest},
		{"db error", "999", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupDeliveryHandler()

			req := httptest.NewRequest(http.MethodPost, "/deliveries/"+tt.id+"/complete", nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.CompleteDelivery(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
