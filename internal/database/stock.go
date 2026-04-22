package database

import (
	"storeSystem/internal/models"

	"github.com/jmoiron/sqlx"
)

type StockStore struct {
	db *sqlx.DB
}

func NewStockStore(db *sqlx.DB) *StockStore {
	return &StockStore{db: db}
}

func (s *StockStore) GetAll() ([]models.Stock, error) {
	var stocks []models.Stock
	query := `SELECT * FROM stock_balance order by quantity desc;`

	err := s.db.Select(&stocks, query)

	if err != nil {
		return nil, err
	}
	return stocks, nil
}
