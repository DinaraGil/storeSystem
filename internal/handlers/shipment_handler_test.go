package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"storeSystem/internal/database/mocks"
)

func setupShipmentHandler() *Handlers {
	return &Handlers{
		shipmentStore:     &mocks.MockShipmentStore{},
		shipmentListStore: &mocks.MockShipmentListStore{},
	}
}

func TestGetAllShipments_Table(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{"success", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupShipmentHandler()

			req := httptest.NewRequest(http.MethodGet, "/shipments", nil)
			rr := httptest.NewRecorder()

			h.GetAllShipments(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetShipmentByID_Table(t *testing.T) {

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

			h := setupShipmentHandler()

			req := httptest.NewRequest(http.MethodGet, "/shipments/"+tt.id, nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.GetShipmentByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCreateShipment_Table(t *testing.T) {

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{"success", `{"status":"new"}`, http.StatusCreated},
		{"invalid json", "{bad json", http.StatusBadRequest},
		{"empty status", `{"status":""}`, http.StatusBadRequest},
		{"db error", `{"status":"db error"}`, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupShipmentHandler()

			req := httptest.NewRequest(
				http.MethodPost,
				"/shipments",
				bytes.NewBufferString(tt.body),
			)

			rr := httptest.NewRecorder()

			h.CreateShipment(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetShipmentListsByShipmentID_Table(t *testing.T) {

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

			h := setupShipmentHandler()

			req := httptest.NewRequest(http.MethodGet, "/shipments/"+tt.id+"/lists", nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.GetShipmentListsByShipmentID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCompleteShipment_Table(t *testing.T) {

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

			h := setupShipmentHandler()

			req := httptest.NewRequest(http.MethodPost, "/shipments/"+tt.id+"/complete", nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.CompleteShipment(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
