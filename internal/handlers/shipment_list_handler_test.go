package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"storeSystem/internal/database/mocks"
)

func setupShipmentListHandler() *Handlers {
	return &Handlers{
		shipmentListStore: &mocks.MockShipmentListStore{},
	}
}

func TestGetAllShipmentLists_Table(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{"success", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupShipmentListHandler()

			req := httptest.NewRequest(http.MethodGet, "/shipment-lists", nil)
			rr := httptest.NewRecorder()

			h.GetAllShipmentLists(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetShipmentListByID_Table(t *testing.T) {

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

			h := setupShipmentListHandler()

			req := httptest.NewRequest(http.MethodGet, "/shipment-lists/"+tt.id, nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.GetShipmentListByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCreateShipmentList_Table(t *testing.T) {

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{"success", `{"article":"test"}`, http.StatusCreated},
		{"invalid json", "{bad json", http.StatusBadRequest},
		{"empty article", `{"article":""}`, http.StatusBadRequest},
		{"db error", `{"article":"db error"}`, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupShipmentListHandler()

			req := httptest.NewRequest(
				http.MethodPost,
				"/shipment-lists",
				bytes.NewBufferString(tt.body),
			)

			rr := httptest.NewRecorder()

			h.CreateShipmentList(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAddShipmentFromFile_Table(t *testing.T) {

	tests := []struct {
		name      string
		line      string
		userID    int
		expectErr bool
	}{
		{"success", "1,2,article,10", 1, false},
		{"empty line", "", 1, true},
		{"bad format", "1,2,article", 1, true},
	}

	h := setupShipmentListHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := h.AddShipmentFromFile(tt.line, tt.userID)

			if tt.expectErr && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.expectErr && res == nil {
				t.Fatalf("expected result but got nil")
			}
		})
	}
}
