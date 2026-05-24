package mocks

import (
	"errors"
	"storeSystem/internal/models"
)

type MockItemStore struct{}

func (m *MockItemStore) GetAll() ([]models.Item, error) {
	return []models.Item{
		{
			ID:      1,
			Article: "Test",
		},
	}, nil
}

func (m *MockItemStore) GetByID(id int) (*models.Item, error) {
	if id == 1 {
		return &models.Item{
			ID:      1,
			Article: "Test item",
		}, nil
	}

	return nil, errors.New("record not found")
}

func (m *MockItemStore) Create(input models.CreateItemInput) (*models.Item, error) {
	if input.Article == "db-error" {
		return nil, errors.New("database error")
	}

	return &models.Item{
		ID:      1,
		Article: input.Article,
	}, nil
}

func (m *MockItemStore) Update(id int, input models.UpdateItemInput) (*models.Item, error) {
	if id == 999 {
		return nil, errors.New("record not found")
	}

	if input.Article != nil && *input.Article == "db-error" {
		return nil, errors.New("database error")
	}

	return &models.Item{
		ID: id,
	}, nil
}
