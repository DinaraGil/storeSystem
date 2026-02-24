package models

import "time"

type Item struct {
	ID             int       `json:"item_id" db:"item_id"`
	RfidID         string    `json:"rfid_id" db:"rfid_id"`
	DeliveryListID int       `json:"delivery_list_id" db:"delivery_list_id"`
	SupplierID     int       `json:"supplier_id" db:"supplier_id"`
	Name           string    `json:"name" db:"name"`
	Article        string    `json:"article" db:"article"`
	CreatedBy      int       `json:"created_by" db:"created_by"`
	UpdatedBy      *int      `json:"updated_by" db:"updated_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateItemInput struct {
	RfidID         string `json:"rfid_id"`
	DeliveryListID int    `json:"delivery_list_id"`
	SupplierID     int    `json:"supplier_id"`
	Name           string `json:"name"`
	Article        string `json:"article"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      *int   `json:"updated_by"`
}

type UpdateItemInput struct {
	DeliveryListID *int    `json:"delivery_list_id"`
	SupplierID     *int    `json:"supplier_id"`
	Name           *string `json:"name"`
	Article        *string `json:"article"`
}
