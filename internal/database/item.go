package database

import (
	"database/sql"
	"fmt"
	"storeSystem/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ItemStore struct {
	db *sqlx.DB
}

func NewItemStore(db *sqlx.DB) *ItemStore {
	return &ItemStore{db: db}
}

func (s *ItemStore) GetAll() ([]models.Item, error) {
	var items []models.Item
	query := `SELECT * FROM item order by created_at desc;`

	err := s.db.Select(&items, query)
	fmt.Println(items, err)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ItemStore) GetByID(id int) (*models.Item, error) {
	var item models.Item
	query := `SELECT item_id, rfid_id, delivery_list_id, supplier_id,
       name, article, created_by, updated_by,
       created_at, updated_at FROM item where item_id=$1;`

	err := s.db.Get(&item, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ItemStore) Create(input models.CreateItemInput) (*models.Item, error) {
	var item models.Item

	query := `
	INSERT INTO item (rfid_id, delivery_list_id, supplier_id, name, article, created_by, updated_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8, $9)
	returning item_id, rfid_id, delivery_list_id, supplier_id, name, article, created_by, updated_by, created_at, updated_at;`

	now := time.Now()

	err := s.db.QueryRowx(query, input.RfidID, input.DeliveryListID, input.SupplierID, input.Name, input.Article, input.CreatedBy, input.UpdatedBy, now, now).StructScan(&item)

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ItemStore) Update(item_id int, input models.UpdateItemInput) (*models.Item, error) {
	item, err := s.GetByID(item_id)
	if err != nil {
		return nil, err
	}
	if input.DeliveryListID != nil {
		item.DeliveryListID = *input.DeliveryListID
	}
	if input.SupplierID != nil {
		item.SupplierID = *input.SupplierID
	}
	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Article != nil {
		item.Article = *input.Article
	}
	item.UpdatedAt = time.Now()

	query := `UPDATE item SET delivery_list_id = $1, supplier_id = $2, name= $3, article = $4, updated_at=$5
	WHERE item_id = $6 returning item_id, delivery_list_id, supplier_id, name, article, updated_at;`

	var updatedItem models.Item

	err = s.db.QueryRowx(query, item.DeliveryListID, item.SupplierID, item.Name, item.Article, item.UpdatedAt, item.ID).StructScan(&updatedItem)
	if err != nil {
		return nil, err
	}
	return &updatedItem, nil
}
