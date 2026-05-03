package models

import "time"

type Shipment struct {
	ID int `json:"shipment_id" db:"shipment_id"`
	//ShipmentNumber int        `json:"shipment_number" db:"shipment_number"`
	Status      string     `json:"status" db:"status"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	CreatedBy   int        `json:"created_by" db:"created_by"`
	CompletedBy *int       `json:"completed_by" db:"completed_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateShipmentInput struct {
	//ShipmentNumber int        `json:"shipment_number" db:"shipment_number"`
	Status    string `json:"status" db:"status"`
	CreatedBy int    `json:"created_by" db:"created_by"`
}
