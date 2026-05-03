package models

import "time"

type ShipmentList struct {
	ID             int       `json:"shipment_list_id" db:"shipment_list_id"`
	ShipmentId     int       `json:"shipment_id" db:"shipment_id"`
	CustomerId     int       `json:"customer_id" db:"customer_id"`
	ExpectedAmount int       `json:"expected_amount" db:"expected_amount"`
	RealAmount     int       `json:"real_amount" db:"real_amount"`
	Status         string    `json:"status" db:"status"`
	Article        string    `json:"article" db:"article"`
	CreatedBy      int       `json:"created_by" db:"created_by"`
	UpdatedBy      *int      `json:"updated_by" db:"updated_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateShipmentListInput struct {
	ShipmentId     int    `json:"shipment_id" db:"shipment_id"`
	CustomerId     int    `json:"customer_id" db:"customer_id"`
	ExpectedAmount int    `json:"expected_amount" db:"expected_amount"`
	Article        string `json:"article" db:"article"`
	CreatedBy      int    `json:"created_by" db:"created_by"`
	UpdatedBy      *int   `json:"updated_by" db:"updated_by"`
}
