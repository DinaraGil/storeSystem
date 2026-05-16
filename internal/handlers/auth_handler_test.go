package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"storeSystem/internal/database/mocks"
	"storeSystem/internal/models"
	"testing"
)

func TestLogin_Table(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "invalid json",
			body:           "{bad json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "user not found",
			body:           `{"username":"unknown","password":"123"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "success case",
			body:           `{"username":"valid","password":"123"}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := &Handlers{
				workerStore: &mocks.MockWorkerStore{},
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			h.Login(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	h := &Handlers{}

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rr := httptest.NewRecorder()

	h.Logout(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestMe_Table(t *testing.T) {
	tests := []struct {
		name           string
		withClaims     bool
		expectedStatus int
	}{
		{"no claims", false, http.StatusUnauthorized},
		{"success", true, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := &Handlers{}

			req := httptest.NewRequest(http.MethodGet, "/me", nil)

			if tt.withClaims {
				req = req.WithContext(
					context.WithValue(req.Context(), UserClaimsKey, &models.Claims{
						UserID:   1,
						Username: "test",
						RoleID:   1,
					}),
				)
			}

			rr := httptest.NewRecorder()

			h.Me(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCreateWorker_Table(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "invalid json",
			body:           "{bad json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation error",
			body:           `{"username":""}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "user exists",
			body:           `{"username":"exists","password":"123","role_id":1}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "success",
			body:           `{"username":"newuser","password":"123","role_id":1}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := &Handlers{
				workerStore: &mocks.MockWorkerStore{},
			}

			req := httptest.NewRequest(http.MethodPost, "/workers", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			h.CreateWorker(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
