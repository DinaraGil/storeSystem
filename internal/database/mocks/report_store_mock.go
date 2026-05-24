package mocks

import (
	"fmt"
	"storeSystem/internal/models"
	"time"
)

type MockReportStore struct {
	ShouldFail bool
}

func (m *MockReportStore) Create(userID int, reportType, fileName, objectID, bucketName, dateFrom, dateTo string) error {
	if m.ShouldFail {
		return fmt.Errorf("db error")
	}
	return nil
}

func (m *MockReportStore) GetByUserID(userID int) ([]models.Report, error) {
	return []models.Report{
		{
			ReportID:   1,
			UserID:     userID,
			ReportType: "stock",
			FileName:   "file.csv",
			ObjectID:   "obj-1",
			DateFrom:   "2024-01-01",
			DateTo:     "2024-01-02",
			CreatedAt:  time.Now(),
		},
	}, nil
}
