package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"storeSystem/internal/database/mocks"
	"storeSystem/internal/models"
	"testing"
)

func setupDeliveryListHandler() *Handlers {
	return &Handlers{
		deliveryListStore: &mocks.MockDeliveryListStore{},
	}
}

func TestGetDeliveryListByID_Table(t *testing.T) {

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

			h := setupDeliveryListHandler()

			req := httptest.NewRequest(http.MethodGet, "/delivery-lists/"+tt.id, nil)
			req = addParam(req, "id", tt.id)

			rr := httptest.NewRecorder()

			h.GetDeliveryListByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCreateDeliveryList_Table(t *testing.T) {

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "success",
			body:           `{"article":"test"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           "{bad json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty article",
			body:           `{"article":""}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "db error",
			body:           `{"article":"db-error"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := setupDeliveryListHandler()

			req := httptest.NewRequest(
				http.MethodPost,
				"/delivery-lists",
				bytes.NewBufferString(tt.body),
			)

			rr := httptest.NewRecorder()

			h.CreateDeliveryList(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAddFromFile_Table(t *testing.T) {

	tests := []struct {
		name      string
		line      string
		userID    int
		expectErr bool
	}{
		{"success", "1,2,article,10", 1, false},
		{"empty line", "", 1, true},
		{"bad format", "1,2,article", 1, true},
		{"db error", "1,2,db-error,10", 1, true},
	}

	h := setupDeliveryListHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := h.AddFromFile(tt.line, tt.userID)

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

func createMultipartRequest(fieldName, fileName, content string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func TestUploadDeliveryList_Success(t *testing.T) {
	h := &Handlers{
		deliveryListStore: &mocks.MockDeliveryListStore{},
	}

	fileContent := "header\n1,2,article,10\n"

	req, err := createMultipartRequest("file", "test.csv", fileContent)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.WithValue(req.Context(), UserClaimsKey, &models.Claims{
		UserID: 1,
	}))

	rr := httptest.NewRecorder()

	h.UploadDeliveryList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestUploadDeliveryList_NoFile(t *testing.T) {
	h := &Handlers{}

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rr := httptest.NewRecorder()

	h.UploadDeliveryList(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}
