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
	query := `SELECT * FROM delivery order by created_at desc;`

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
	INSERT INTO delivery (status, planned_arrival_at, accepted_at, created_by, accepted_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5 ,$6, $7)
	returning delivery_id, status, planned_arrival_at, accepted_at, created_by, accepted_by, created_at, updated_at;`

	now := time.Now()

	err := s.db.QueryRowx(query, input.Status, input.PlannedArrivalAt, input.AcceptedAt, input.CreatedBy, input.AcceptedBy, now, now).StructScan(&del)

	if err != nil {
		return nil, err
	}
	return &del, nil
}

func (s *DeliveryStore) Update(delivery_id int, input models.UpdateDeliveryInput) (*models.Delivery, error) {
	del, err := s.GetByID(delivery_id)
	fmt.Println("get by id", del, err)
	if err != nil {
		return nil, err
	}
	if input.Status != nil {
		del.Status = *input.Status
	}
	if input.PlannedArrivalAt != nil {
		del.PlannedArrivalAt = *input.PlannedArrivalAt
	}
	if input.AcceptedAt != nil {
		del.AcceptedAt = input.AcceptedAt
	}
	if input.AcceptedBy != nil {
		del.AcceptedBy = *input.AcceptedBy
	}
	del.UpdatedAt = time.Now()

	query := `UPDATE delivery SET status = $1, planned_arrival_at= $2, accepted_at = $3, accepted_by=$4, updated_at=$5
	WHERE delivery_id = $6 returning delivery_id, status, planned_arrival_at, accepted_at, accepted_by, updated_at;`

	var updatedDel models.Delivery

	err = s.db.QueryRowx(query, del.Status, del.PlannedArrivalAt, del.AcceptedAt, del.AcceptedBy, del.UpdatedAt, del.ID).StructScan(&updatedDel)
	fmt.Println("query rowx", updatedDel, err)
	if err != nil {
		return nil, err
	}
	return &updatedDel, nil
}
