package mocks

import (
	"errors"
	"storeSystem/internal/database"
	"storeSystem/internal/models"
	"time"
)

type MockShipmentListStore struct {
}

func (m *MockShipmentListStore) GetAll() ([]models.ShipmentList, error) {
	return []models.ShipmentList{
		{
			ID:             1,
			ShipmentId:     1,
			Article:        "test",
			ExpectedAmount: 10,
			Status:         "NEW",
		},
	}, nil
}

func (m *MockShipmentListStore) GetByID(id int) (*models.ShipmentList, error) {
	if id == 1 {
		return &models.ShipmentList{
			ID:      1,
			Article: "test",
			Status:  "NEW",
		}, nil
	}

	return nil, errors.New("shipment_list not found")
}

func (m *MockShipmentListStore) GetByShipmentID(id int) ([]models.ShipmentList, error) {
	if id == 1 {
		return []models.ShipmentList{
			{
				ID:         1,
				ShipmentId: 1,
				Article:    "test",
			},
		}, nil
	}

	return nil, errors.New("not found")
}

func (m *MockShipmentListStore) Create(input models.CreateShipmentListInput) (*models.ShipmentList, error) {
	if input.Article == "db error" {
		return nil, errors.New("database error")
	}
	return &models.ShipmentList{
		ID:             1,
		ShipmentId:     input.ShipmentId,
		CustomerId:     input.CustomerId,
		ExpectedAmount: input.ExpectedAmount,
		Article:        input.Article,
		Status:         "NEW",
		CreatedAt:      time.Now(),
	}, nil
}

func (m *MockShipmentListStore) ProcessScannerEvent(
	shipmentID int,
	evt models.Event,
	workerID int,
) (*database.ShipmentListUpdateDTO, error) {

	if evt.Error != nil && *evt.Error != "" {
		return nil, errors.New("event error")
	}

	if evt.Article == nil {
		return nil, errors.New("article is nil")
	}

	if *evt.Article == "fail" {
		return nil, errors.New("processing failed")
	}

	return &database.ShipmentListUpdateDTO{
		ShipmentListID: 1,
		ShipmentID:     shipmentID,
		CustomerID:     1,
		ExpectedAmount: 10,
		RealAmount:     1,
		StockAvailable: 5,
		StockReserved:  3,
		Article:        *evt.Article,
		Status:         "NOT_ENOUGH",
		UpdatedAt:      time.Now(),
	}, nil
}
