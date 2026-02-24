package models

import "time"

type DeliveryList struct {
	ID         int       `json:"delivery_list_id" db:"delivery_list_id"`
	DeliveryId int       `json:"delivery_id" db:"delivery_id"`
	SupplierId int       `json:"supplier_id" db:"supplier_id"`
	Amount     int       `json:"amount" db:"amount"`
	Article    string    `json:"article" db:"article"`
	CreatedBy  int       `json:"created_by" db:"created_by"`
	UpdatedBy  *int      `json:"updated_by" db:"updated_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
type CreateDeliveryListInput struct {
	DeliveryId int    `json:"delivery_id" db:"delivery_id"`
	SupplierId int    `json:"supplier_id" db:"supplier_id"`
	Amount     int    `json:"amount" db:"amount"`
	Article    string `json:"article" db:"article"`
	CreatedBy  int    `json:"created_by" db:"created_by"`
	UpdatedBy  *int   `json:"updated_by" db:"updated_by"`
}
