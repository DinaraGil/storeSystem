package database

import (
	"storeSystem/internal/models"

	"github.com/jmoiron/sqlx"
)

type ReportStore struct {
	db *sqlx.DB
}

func NewReportStore(db *sqlx.DB) *ReportStore {
	return &ReportStore{db: db}
}

func (s *ReportStore) Create(userID int, reportType, fileName, objectID, bucketName, dateFrom, dateTo string) error {
	_, err := s.db.Exec(`
		INSERT INTO report (
			user_id, report_type, file_name, object_id, bucket_name, date_from, date_to
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, reportType, fileName, objectID, bucketName, dateFrom, dateTo)

	return err
}

func (s *ReportStore) GetByUserID(userID int) ([]models.Report, error) {
	rows, err := s.db.Query(`
		SELECT report_id, user_id, report_type, file_name, object_id, bucket_name,
		       date_from, date_to, created_at
		FROM report
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var r models.Report
		err := rows.Scan(
			&r.ReportID,
			&r.UserID,
			&r.ReportType,
			&r.FileName,
			&r.ObjectID,
			&r.BucketName,
			&r.DateFrom,
			&r.DateTo,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, rows.Err()
}
