package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"storeSystem/internal/config"
	"storeSystem/internal/database/mocks"
	"storeSystem/internal/models"
	"testing"
)

func setupReportHandler() *Handlers {
	return &Handlers{
		stockStore:    &mocks.MockStockStore{},
		deliveryStore: &mocks.MockDeliveryStore{},
		minioService:  &mocks.MockMinio{},
		reportStore:   &mocks.MockReportStore{},
	}
}

func init() {
	config.AppConfig = &config.Config{
		BucketName: "test",
	}

}

func TestGenerateReportFile_Table(t *testing.T) {
	tests := []struct {
		name        string
		reportType  string
		expectError bool
	}{
		{"stock success", "stock", false},
		{"unknown type", "bad", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupReportHandler()

			res, err := h.GenerateReportFile(&ReportConfig{
				ReportType: tt.reportType,
				UserID:     1,
				DateFrom:   "2026-01-01",
				DateTo:     "2026-02-02",
			})

			if tt.expectError && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.expectError {
				if res == nil {
					t.Fatal("expected result")
				}
				if (*res)["filename"] == "" {
					t.Fatal("empty filename")
				}
			}
		})
	}
}

func TestGetUsersReports_Table(t *testing.T) {
	tests := []struct {
		name           string
		withClaims     bool
		minioFail      bool
		expectedStatus int
	}{
		{"no claims", false, false, http.StatusBadRequest},
		{"minio fail skipped", true, true, http.StatusOK},
		{"success", true, false, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupReportHandler()

			h.minioService = &mocks.MockMinio{ShouldFail: tt.minioFail}

			req := httptest.NewRequest(http.MethodGet, "/reports", nil)

			if tt.withClaims {
				req = req.WithContext(context.WithValue(req.Context(), UserClaimsKey, &models.Claims{
					UserID: 1,
				}))
			}

			rr := httptest.NewRecorder()

			h.GetUsersReports(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
