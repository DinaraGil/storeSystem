package models

import "time"

type Stock struct {
	Article   string    `json:"article" db:"article"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Reserved  int       `json:"reserved" db:"reserved"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
