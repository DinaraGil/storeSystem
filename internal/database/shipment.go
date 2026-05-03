package database

import (
	"database/sql"
	"fmt"
	"storeSystem/internal/models"

	"github.com/jmoiron/sqlx"
)

type ShipmentStore struct {
	db *sqlx.DB
}

func NewShipmentStore(db *sqlx.DB) *ShipmentStore {
	return &ShipmentStore{db: db}
}

func (s *ShipmentStore) GetAll() ([]models.Shipment, error) {
	var ships []models.Shipment
	query := `SELECT * FROM shipment order by shipment_id asc;`

	err := s.db.Select(&ships, query)

	if err != nil {
		return nil, err
	}
	return ships, nil
}

func (s *ShipmentStore) GetErrorShipments() ([]models.Shipment, error) {
	var ships []models.Shipment
	query := `SELECT * FROM shipment WHERE status = 'ERROR' order by shipment_id asc;`

	err := s.db.Select(&ships, query)

	if err != nil {
		return nil, err
	}
	return ships, nil
}

func (s *ShipmentStore) GetByID(id int) (*models.Shipment, error) {
	var ships models.Shipment
	query := `SELECT * FROM shipment where shipment_id=$1;`

	err := s.db.Get(&ships, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ship with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &ships, nil
}

func (s *ShipmentStore) Create(input models.CreateShipmentInput) (*models.Shipment, error) {
	var ship models.Shipment

	query := `
	INSERT INTO shipment (status, created_by)
	VALUES ($1, $2)
	returning shipment_id, status, completed_at, created_by, completed_by, created_at, updated_at;`

	err := s.db.QueryRowx(query, input.Status, input.CreatedBy).StructScan(&ship)

	if err != nil {
		return nil, err
	}
	return &ship, nil
}

func (s *ShipmentStore) CompleteShipment(shipmentID int) error {
	_, err := s.db.Exec(`
		UPDATE shipment
		SET status = CASE
			WHEN NOT EXISTS (
				SELECT 1
				FROM shipment_list
				WHERE shipment_id = $1
				  AND status != 'COMPLETED'
			)
			THEN 'COMPLETED'

			WHEN NOT EXISTS (
				SELECT 1
				FROM shipment_list
				WHERE shipment_id = $1
				  AND status != 'NEW'
			)
			THEN 'NEW'

			ELSE 'ERROR'
		END,
		completed_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
		WHERE shipment_id = $1
	`, shipmentID)
	fmt.Println(err)
	return err
}
