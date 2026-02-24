package models

import "time"

type Delivery struct {
	ID               int        `json:"delivery_id" db:"delivery_id"`
	Status           string     `json:"status" db:"status"`
	PlannedArrivalAt time.Time  `json:"planned_arrival_at" db:"planned_arrival_at"`
	AcceptedAt       *time.Time `json:"accepted_at" db:"accepted_at"`
	CreatedBy        int        `json:"created_by" db:"created_by"`
	AcceptedBy       int        `json:"accepted_by" db:"accepted_by"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}
type CreateDeliveryInput struct {
	Status           string     `json:"status" db:"status"`
	PlannedArrivalAt time.Time  `json:"planned_arrival_at" db:"planned_arrival_at"`
	AcceptedAt       *time.Time `json:"accepted_at" db:"accepted_at"`
	CreatedBy        int        `json:"created_by" db:"created_by"`
	AcceptedBy       *int       `json:"accepted_by" db:"accepted_by"`
}

type UpdateDeliveryInput struct {
	Status           *string    `json:"status" db:"status"`
	PlannedArrivalAt *time.Time `json:"planned_arrival_at" db:"planned_arrival_at"`
	AcceptedAt       *time.Time `json:"accepted_at" db:"accepted_at"`
	AcceptedBy       *int       `json:"accepted_by" db:"accepted_by"`
}
