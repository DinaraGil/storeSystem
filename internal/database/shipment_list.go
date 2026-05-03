package database

import (
	"database/sql"
	"fmt"
	"storeSystem/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ShipmentListStore struct {
	db *sqlx.DB
}

func NewShipmentListStore(db *sqlx.DB) *ShipmentListStore {
	return &ShipmentListStore{db: db}
}

func (s *ShipmentListStore) GetAll() ([]models.ShipmentList, error) {
	var lists []models.ShipmentList
	query := `SELECT * FROM shipment_list ORDER BY created_at DESC;`

	err := s.db.Select(&lists, query)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (s *ShipmentListStore) GetByID(id int) (*models.ShipmentList, error) {
	var list models.ShipmentList
	query := `SELECT * FROM shipment_list WHERE shipment_list_id=$1;`

	err := s.db.Get(&list, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shipment_list with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ShipmentListStore) GetByShipmentID(id int) ([]models.ShipmentList, error) {
	var lists []models.ShipmentList
	query := `SELECT * FROM shipment_list WHERE shipment_id=$1 ORDER BY shipment_list_id;`

	err := s.db.Select(&lists, query, id)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (s *ShipmentListStore) Create(input models.CreateShipmentListInput) (*models.ShipmentList, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var list models.ShipmentList

	query := `
	INSERT INTO shipment_list (
		shipment_id, customer_id,
		expected_amount, article,
		created_by, updated_by, created_at, updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING 
		shipment_list_id, shipment_id, customer_id,
		expected_amount, status, article,
		created_by, updated_by, created_at, updated_at;
	`

	now := time.Now()

	err = tx.QueryRowx(
		query,
		input.ShipmentId,
		input.CustomerId,
		input.ExpectedAmount,
		input.Article,
		input.CreatedBy,
		input.UpdatedBy,
		now,
		now,
	).StructScan(&list)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO stock_balance (article, reserved, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (article)
		DO UPDATE SET
			reserved = stock_balance.reserved + EXCLUDED.reserved,
			updated_at = CURRENT_TIMESTAMP
    `, list.Article, list.ExpectedAmount)
	if err != nil {
		return nil, err
	}

	// фиксируем
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &list, nil
}

type ShipmentListUpdateDTO struct {
	ShipmentListID int       `json:"shipment_list_id"`
	ShipmentID     int       `json:"shipment_id"`
	CustomerID     int       `json:"customer_id"`
	ExpectedAmount int       `json:"expected_amount"`
	RealAmount     int       `json:"real_amount"`
	StockAvailable int       `json:"stock_available"`
	StockReserved  int       `json:"stock_reserved"`
	Article        string    `json:"article"`
	Status         string    `json:"status"`
	UpdatedBy      int       `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s *ShipmentListStore) ProcessScannerEvent(shipmentID int, evt models.Event, workerID int) (*ShipmentListUpdateDTO, error) {

	if evt.Error != nil && *evt.Error != "" {
		return nil, fmt.Errorf("error in event")
	}

	// отгрузка = товар выходит
	if evt.IsIn == nil || *evt.IsIn {
		return nil, fmt.Errorf("the rfid is not outgoing")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var shipmentListID int

	if evt.Article == nil {
		return nil, fmt.Errorf("article is nil in event")
	}
	
	err = tx.QueryRow(`
		SELECT shipment_list_id
		FROM shipment_list
		WHERE shipment_id = $1 AND article = $2
	`, shipmentID, *evt.Article).Scan(&shipmentListID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shipment_list not found")
	}
	if err != nil {
		return nil, err
	}

	// 🔎 проверяем товар
	var itemID int
	err = tx.QueryRow(`
		SELECT item_id
		FROM item
		WHERE rfid_id = $1 AND status = 'STOCKED'
	`, evt.RfidId).Scan(&itemID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found or already shipped")
	}
	if err != nil {
		return nil, err
	}

	// 🔥 блокируем строку склада
	var quantity, reserved int
	err = tx.QueryRow(`
		SELECT quantity, reserved
		FROM stock_balance
		WHERE article = $1
		FOR UPDATE
	`, *evt.Article).Scan(&quantity, &reserved)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("stock not found")
	}
	if err != nil {
		return nil, err
	}

	//available := quantity - reserved
	if quantity <= 0 {
		return nil, fmt.Errorf("no available stock")
	}

	// 🔥 атомарно уменьшаем склад и резерв
	res, err := tx.Exec(`
		UPDATE stock_balance
		SET 
		    quantity = quantity - 1,
		    reserved = reserved - 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE article = $1
		  AND quantity > 0
		  AND reserved > 0
	`, *evt.Article)

	if err != nil {
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no reserved stock available")
	}

	// 📦 помечаем товар как отгруженный
	_, err = tx.Exec(`
		UPDATE item
		SET status = 'SHIPPED',
		    updated_at = CURRENT_TIMESTAMP
		WHERE item_id = $1
	`, itemID)

	if err != nil {
		return nil, err
	}

	// 📊 обновляем shipment_list
	var updated ShipmentListUpdateDTO

	err = tx.QueryRow(`
		UPDATE shipment_list
		SET 
		    real_amount = real_amount + 1,
		    status = CASE
		        WHEN real_amount + 1 = expected_amount THEN 'COMPLETED'
		        WHEN real_amount + 1 > expected_amount THEN 'OVER_SHIPPED'
		        ELSE 'NOT_ENOUGH'
		    END,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $1
		WHERE shipment_list_id = $2
		RETURNING 
			shipment_list_id,
			shipment_id,
			customer_id,
			expected_amount,
			real_amount,
			article,
			status,
			updated_by,
			updated_at
	`, workerID, shipmentListID).Scan(
		&updated.ShipmentListID,
		&updated.ShipmentID,
		&updated.CustomerID,
		&updated.ExpectedAmount,
		&updated.RealAmount,
		&updated.Article,
		&updated.Status,
		&updated.UpdatedBy,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// 📦 обновляем shipment
	_, err = tx.Exec(`
		UPDATE shipment
		SET updated_at = CURRENT_TIMESTAMP
		WHERE shipment_id = $1
	`, shipmentID)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	updated.StockAvailable = quantity - 1
	updated.StockReserved = reserved - 1
	return &updated, nil
}
