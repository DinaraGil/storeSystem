package mocks

import (
	"errors"

	"storeSystem/internal/models"
)

type MockCounterpartyStore struct{}

func (m *MockCounterpartyStore) GetAll() ([]models.Counterparty, error) {
	return []models.Counterparty{
		{
			ID:       1,
			FullName: "Test Counterparty",
		},
	}, nil
}

func (m *MockCounterpartyStore) GetByID(id int) (*models.Counterparty, error) {

	if id == 1 {
		return &models.Counterparty{
			ID:       1,
			FullName: "Test Counterparty",
		}, nil
	}

	return nil, errors.New("counterparty not found")
}

func (m *MockCounterpartyStore) Create(input models.CreateCounterpartyInput) (*models.Counterparty, error) {

	if input.FullName == "db-error" {
		return nil, errors.New("database error")
	}

	return &models.Counterparty{
		ID:       1,
		FullName: input.FullName,
	}, nil
}
