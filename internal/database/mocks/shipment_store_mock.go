package mocks

import (
	"errors"
	"storeSystem/internal/models"
)

type MockShipmentStore struct{}

func (m *MockShipmentStore) GetAll() ([]models.Shipment, error) {
	return []models.Shipment{
		{
			ID:     1,
			Status: "NEW",
		},
	}, nil
}

func (m *MockShipmentStore) GetErrorShipments() ([]models.Shipment, error) {
	return []models.Shipment{
		{
			ID:     2,
			Status: "ERROR",
		},
	}, nil
}

func (m *MockShipmentStore) GetByID(id int) (*models.Shipment, error) {
	if id == 1 {
		return &models.Shipment{
			ID:     1,
			Status: "NEW",
		}, nil
	}

	return nil, errors.New("shipment not found")
}

func (m *MockShipmentStore) Create(input models.CreateShipmentInput) (*models.Shipment, error) {
	if input.Status == "db error" {
		return nil, errors.New("database error")
	}

	return &models.Shipment{
		ID:     1,
		Status: input.Status,
	}, nil
}

func (m *MockShipmentStore) CompleteShipment(id int) error {
	if id == 999 {
		return errors.New("db error")
	}
	return nil
}
