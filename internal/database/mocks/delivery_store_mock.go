package mocks

import (
	"errors"
	"storeSystem/internal/models"
)

type MockDeliveryStore struct{}

func (m *MockDeliveryStore) GetAll() ([]models.Delivery, error) {
	return []models.Delivery{
		{ID: 1, Status: "new"},
	}, nil
}

func (m *MockDeliveryStore) GetErrorDeliveries() ([]models.Delivery, error) {
	return []models.Delivery{
		{ID: 2, Status: "ERROR"},
	}, nil
}

func (m *MockDeliveryStore) GetByID(id int) (*models.Delivery, error) {
	if id == 1 {
		return &models.Delivery{ID: 1, Status: "new"}, nil
	}
	return nil, errors.New("not found")
}

func (m *MockDeliveryStore) Create(input models.CreateDeliveryInput) (*models.Delivery, error) {
	if input.Status == "db-error" {
		return nil, errors.New("database error")
	}

	return &models.Delivery{ID: 1, Status: input.Status}, nil
}

func (m *MockDeliveryStore) CompleteDelivery(id int) error {
	if id == 999 {
		return errors.New("db error")
	}
	return nil
}
