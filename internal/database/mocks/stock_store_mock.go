package mocks

import (
	"database/sql"
	"errors"

	"storeSystem/internal/models"
)

type MockStockStore struct{}

func (m *MockStockStore) GetAll() ([]models.Stock, error) {

	// success
	return []models.Stock{
		{
			Article:  "1",
			Quantity: 10,
		},
	}, nil
}

type EmptyStockStoreMock struct{}

func (m *EmptyStockStoreMock) GetAll() ([]models.Stock, error) {
	return nil, sql.ErrNoRows
}

type ErrorStockStoreMock struct{}

func (m *ErrorStockStoreMock) GetAll() ([]models.Stock, error) {
	return nil, errors.New("database error")
}
