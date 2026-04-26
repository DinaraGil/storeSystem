package models

import "time"

type Report struct {
	ReportID   int       `json:"report_id" db:"report_id"`
	UserID     int       `json:"user_id" db:"user_id"`
	ReportType string    `json:"report_type" db:"report_type"`
	FileName   string    `json:"file_name" db:"file_name"`
	ObjectID   string    `json:"object_id" db:"object_id"`
	BucketName string    `json:"bucket_name" db:"bucket_name"`
	DateFrom   string    `json:"date_from" db:"date_from"`
	DateTo     string    `json:"date_to" db:"date_to"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
