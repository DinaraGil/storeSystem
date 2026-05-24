package mocks

import (
	"errors"
	"storeSystem/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type MockWorkerStore struct{}

func (m *MockWorkerStore) GetByUsername(username string) (*models.Worker, error) {
	if username == "valid" {
		hash, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

		return &models.Worker{
			ID:           1,
			Username:     "valid",
			PasswordHash: string(hash),
			RoleId:       1,
		}, nil
	}

	return nil, errors.New("not found")
}

func (m *MockWorkerStore) Create(input models.CreateWorkerInput) (*models.Worker, error) {
	if input.Username == "exists" {
		return nil, errors.New("already exists")
	}

	return &models.Worker{
		ID:       1,
		Username: input.Username,
		RoleId:   input.RoleId,
	}, nil
}
