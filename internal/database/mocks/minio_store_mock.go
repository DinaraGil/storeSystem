package mocks

import (
	"fmt"
	"storeSystem/internal/helpers"
)

type MockMinio struct {
	ShouldFail bool
}

func (m *MockMinio) CreateOne(data helpers.FileDataType) (*helpers.UploadedFile, error) {
	if m.ShouldFail {
		return nil, fmt.Errorf("minio error")
	}

	return &helpers.UploadedFile{
		ObjectID: "obj-1",
		Link:     "http://file",
	}, nil
}

func (m *MockMinio) GetOne(objectID string) (string, error) {
	if m.ShouldFail {
		return "", fmt.Errorf("minio error")
	}

	return "http://file", nil
}
