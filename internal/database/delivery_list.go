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

func (s *DeliveryListStore) GetByDeliveryID(id int) ([]models.DeliveryList, error) {
	var delLists []models.DeliveryList
	query := `SELECT * FROM delivery_list where delivery_id=$1;`

	err := s.db.Select(&delLists, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("delList with delivery id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return delLists, nil
}

func (s *DeliveryListStore) Create(input models.CreateDeliveryListInput) (*models.DeliveryList, error) {
	var delList models.DeliveryList

	query := `
	INSERT INTO delivery_list (delivery_id, supplier_id, expected_amount, article, created_by, updated_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8)
	returning delivery_list_id, delivery_id, supplier_id, expected_amount, status, article, created_by, updated_by, created_at, updated_at;`

	now := time.Now()

	err := s.db.QueryRowx(query, input.DeliveryId, input.SupplierId, input.ExpectedAmount, input.Article, input.CreatedBy, input.UpdatedBy, now, now).StructScan(&delList)

	if err != nil {
		return nil, err
	}
	return &delList, nil
}

func (s *DeliveryListStore) ProcessScannerEvent(deliveryID int, evt models.Event) error {
	if evt.Error != nil && *evt.Error != "" {
		return fmt.Errorf("error in event")
	}

	if evt.IsIn == nil || !*evt.IsIn {
		return fmt.Errorf("the rfid is out")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var supplierID int
	//var expectedAmount int
	//var realAmount int

	var deliveryListID int

	err = tx.QueryRow(`
		SELECT delivery_list_id, supplier_id
		FROM delivery_list
		WHERE delivery_id = $1 AND article = $2
	`, deliveryID, evt.Article).Scan(&deliveryListID, &supplierID)

	if err != nil {
		return err
	}

	var existingItemID int
	err = tx.QueryRow(`
		SELECT item_id
		FROM item
		WHERE rfid_id = $1
	`, evt.RfidId).Scan(&existingItemID)

	if err == nil {
		return fmt.Errorf("item already exist")
	}
	if err != sql.ErrNoRows {
		return err
	}

	var itemID int
	err = tx.QueryRow(`
		INSERT INTO item (
			rfid_id, delivery_list_id, supplier_id, name, article
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING item_id
	`, evt.RfidId, deliveryListID, supplierID, "aboba", evt.Article).Scan(&itemID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE delivery_list
		SET real_amount = real_amount + 1,
		    status = CASE
		        WHEN real_amount + 1 = expected_amount THEN 'COMPLETED'
		        WHEN real_amount + 1 > expected_amount THEN 'OVERMUCH'
		        ELSE 'IN_PROGRESS'
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE delivery_list_id = $1
	`, deliveryListID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO stock_balance (article, quantity, updated_at)
		VALUES ($1, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (article)
		DO UPDATE SET
			quantity = stock_balance.quantity + 1,
			updated_at = CURRENT_TIMESTAMP
	`, evt.Article)
	if err != nil {
		return err
	}

	return tx.Commit()
}
