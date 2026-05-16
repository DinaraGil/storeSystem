package mocks

import (
	"errors"
	"storeSystem/internal/database"
	"storeSystem/internal/models"
	"time"
)

type MockDeliveryListStore struct{}

// GET ALL
func (m *MockDeliveryListStore) GetAll() ([]models.DeliveryList, error) {
	return []models.DeliveryList{
		{
			ID:      1,
			Article: "Test article",
			Status:  "NEW",
		},
	}, nil
}

// GET BY ID
func (m *MockDeliveryListStore) GetByID(id int) (*models.DeliveryList, error) {
	if id == 1 {
		return &models.DeliveryList{
			ID:      1,
			Article: "Test article",
			Status:  "NEW",
		}, nil
	}

	return nil, errors.New("not found")
}

// GET BY DELIVERY ID
func (m *MockDeliveryListStore) GetByDeliveryID(id int) ([]models.DeliveryList, error) {
	if id == 1 {
		return []models.DeliveryList{
			{
				ID:      1,
				Article: "10",
			},
		}, nil
	}

	return nil, errors.New("not found")
}

// CREATE
func (m *MockDeliveryListStore) Create(input models.CreateDeliveryListInput) (*models.DeliveryList, error) {

	if input.Article == "db-error" {
		return nil, errors.New("database error")
	}

	return &models.DeliveryList{
		ID:             1,
		DeliveryId:     input.DeliveryId,
		SupplierId:     input.SupplierId,
		Article:        input.Article,
		ExpectedAmount: input.ExpectedAmount,
		Status:         "NEW",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// PROCESS SCANNER EVENT (упрощённый mock)
func (m *MockDeliveryListStore) ProcessScannerEvent(deliveryID int, evt models.Event, workerID int) (*database.DeliveryListUpdateDTO, error) {

	if evt.Error != nil && *evt.Error != "" {
		return nil, errors.New("event error")
	}

	if evt.Article == nil {
		return nil, errors.New("article is nil")
	}

	if *evt.Article == "fail" {
		return nil, errors.New("processing failed")
	}

	return &database.DeliveryListUpdateDTO{
		DeliveryListID: 1,
		DeliveryID:     deliveryID,
		SupplierID:     1,
		ExpectedAmount: 10,
		RealAmount:     1,
		Article:        *evt.Article,
		Status:         "NOT_ENOUGH",
		UpdatedAt:      time.Now(),
	}, nil
}
