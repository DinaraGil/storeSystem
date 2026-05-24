package database

import (
	"database/sql"
	"fmt"
	"storeSystem/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type DeliveryStore struct {
	db *sqlx.DB
}

func NewDeliveryStore(db *sqlx.DB) *DeliveryStore {
	return &DeliveryStore{db: db}
}

func (s *DeliveryStore) GetAll() ([]models.Delivery, error) {
	var del []models.Delivery
	query := `SELECT * FROM delivery order by delivery_id asc;`

	err := s.db.Select(&del, query)

	if err != nil {
		return nil, err
	}
	return del, nil
}

func (s *DeliveryStore) GetErrorDeliveries() ([]models.Delivery, error) {
	var del []models.Delivery
	query := `SELECT * FROM delivery WHERE status = 'ERROR' order by delivery_id asc;`

	err := s.db.Select(&del, query)

	if err != nil {
		return nil, err
	}
	return del, nil
}

func (s *DeliveryStore) GetByID(id int) (*models.Delivery, error) {
	var del models.Delivery
	query := `SELECT * FROM delivery where delivery_id=$1;`

	err := s.db.Get(&del, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("del with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &del, nil
}

func (s *DeliveryStore) Create(input models.CreateDeliveryInput) (*models.Delivery, error) {
	var del models.Delivery

	query := `
	INSERT INTO delivery (status, accepted_at, created_by, accepted_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5 ,$6, $7)
	returning delivery_id, status, accepted_at, created_by, accepted_by, created_at, updated_at;`

	now := time.Now()

	err := s.db.QueryRowx(query, input.Status, input.AcceptedAt, input.CreatedBy, input.AcceptedBy, now, now).StructScan(&del)

	if err != nil {
		return nil, err
	}
	return &del, nil
}

func (s *DeliveryStore) CompleteDelivery(deliveryID int) error {
	_, err := s.db.Exec(`
		UPDATE delivery
		SET status = CASE
			WHEN NOT EXISTS (
				SELECT 1
				FROM delivery_list
				WHERE delivery_id = $1
				  AND status != 'COMPLETED'
			)
			THEN 'COMPLETED'
		    WHEN NOT EXISTS (
				SELECT 1
				FROM delivery_list
				WHERE delivery_id = $1
				  AND status != 'NEW'
			)
		    THEN 'NEW'
			ELSE 'ERROR'
		END,
		updated_at = CURRENT_TIMESTAMP
		WHERE delivery_id = $1
	`, deliveryID)
	fmt.Println(err)
	return err
}
