package models

import "time"

type Stock struct {
	Article   string    `json:"article" db:"article"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
