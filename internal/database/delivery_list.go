package database

import (
	"database/sql"
	"fmt"
	"storeSystem/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type DeliveryListStore struct {
	db *sqlx.DB
}

func NewDeliveryListStore(db *sqlx.DB) *DeliveryListStore {
	return &DeliveryListStore{db: db}
}

func (s *DeliveryListStore) GetAll() ([]models.DeliveryList, error) {
	var delLists []models.DeliveryList
	query := `SELECT * FROM delivery_list order by created_at desc;`

	err := s.db.Select(&delLists, query)

	if err != nil {
		return nil, err
	}
	return delLists, nil
}

func (s *DeliveryListStore) GetByID(id int) (*models.DeliveryList, error) {
	var delList models.DeliveryList
	query := `SELECT * FROM delivery_list where delivery_list_id=$1;`

	err := s.db.Get(&delList, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("delList with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &delList, nil
}

func (s *DeliveryListStore) Create(input models.CreateDeliveryListInput) (*models.DeliveryList, error) {
	var delList models.DeliveryList

	query := `
	INSERT INTO delivery_list (delivery_id, supplier_id, amount, article, created_by, updated_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8)
	returning delivery_list_id, delivery_id, supplier_id, amount, article, created_by, updated_by, created_at, updated_at;`

	now := time.Now()

	err := s.db.QueryRowx(query, input.DeliveryId, input.SupplierId, input.Amount, input.Article, input.CreatedBy, input.UpdatedBy, now, now).StructScan(&delList)

	if err != nil {
		return nil, err
	}
	return &delList, nil
}
