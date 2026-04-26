package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"storeSystem/internal/config"
	"storeSystem/internal/helpers"
	"strconv"
	"time"
)

type ReportConfig struct {
	ReportType string `json:"report_type"`
	DateFrom   string `json:"date_from"`
	DateTo     string `json:"date_to"`
	UserID     int
}

func (h *Handlers) writeStockCSV(writer *csv.Writer) error {
	stocks, err := h.stockStore.GetAll()
	if err != nil {
		return fmt.Errorf("ошибка получения остатков")
	}

	if err := writer.Write([]string{"Артикул", "Количество", "Время изменения"}); err != nil {
		return err
	}

	for _, stock := range stocks {
		record := []string{
			stock.Article,
			strconv.Itoa(stock.Quantity),
			stock.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handlers) writeErrorDeliveriesCSV(writer *csv.Writer) error {
	errorDeliveries, err := h.deliveryStore.GetErrorDeliveries()
	if err != nil {
		return fmt.Errorf("ошибка получения ошибочных поставок")
	}

	if err := writer.Write([]string{"Статус", "Время приема", "Создано сотрудником", "Принято сотрудником", "Время создания", "Время обновления"}); err != nil {
		return err
	}

	for _, del := range errorDeliveries {
		record := []string{
			del.Status,
			del.AcceptedAt.Format("2006-01-02 15:04:05"),
			strconv.Itoa(del.CreatedBy),
			strconv.Itoa(del.AcceptedBy),
			del.CreatedAt.Format("2006-01-02 15:04:05"),
			del.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handlers) GenerateReportFile(reportConfig *ReportConfig) (*map[string]string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	var writeCSVError error
	switch reportConfig.ReportType {
	case "stock":
		writeCSVError = h.writeStockCSV(writer)
	case "delivery_errors":
		writeCSVError = h.writeErrorDeliveriesCSV(writer)
	default:
		writeCSVError = fmt.Errorf("Неизвестный тип отчета")
	}

	if writeCSVError != nil {
		return nil, writeCSVError
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	filename := reportConfig.ReportType + "_" + time.Now().Format("2006-01-02_15-04-05") + ".csv"

	fileData := helpers.FileDataType{
		FileName: filename,
		Data:     buf.Bytes(),
	}

	uploaded, err := h.minioService.CreateOne(fileData)
	if err != nil {
		return nil, fmt.Errorf("unable to save the file: %w", err)
	}

	err = h.reportStore.Create(
		reportConfig.UserID,
		reportConfig.ReportType,
		filename,
		uploaded.ObjectID,
		config.AppConfig.BucketName,
		reportConfig.DateFrom,
		reportConfig.DateTo,
	)
	if err != nil {
		return nil, err
	}

	return &map[string]string{
		"link":     uploaded.Link,
		"filename": filename,
	}, nil
}

func (h *Handlers) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var config ReportConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	claims, ok := GetUserClaimsFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusBadRequest, "no claims")
		return
	}
	config.UserID = claims.UserID

	result, err := h.GenerateReportFile(&config)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"link":     (*result)["link"],
		"filename": (*result)["filename"],
	})
}

func (h *Handlers) GetUsersReports(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserClaimsFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusBadRequest, "no claims")
		return
	}
	userID := claims.UserID

	reports, err := h.reportStore.GetByUserID(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ReportResponse struct {
		ReportID   int       `json:"report_id"`
		ReportType string    `json:"report_type"`
		FileName   string    `json:"file_name"`
		Link       string    `json:"link"`
		DateFrom   string    `json:"date_from"`
		DateTo     string    `json:"date_to"`
		CreatedAt  time.Time `json:"created_at"`
	}

	result := make([]ReportResponse, 0, len(reports))

	for _, report := range reports {
		link, err := h.minioService.GetOne(report.ObjectID)
		if err != nil {
			continue
		}

		result = append(result, ReportResponse{
			ReportID:   report.ReportID,
			ReportType: report.ReportType,
			FileName:   report.FileName,
			Link:       link,
			DateFrom:   report.DateFrom,
			DateTo:     report.DateTo,
			CreatedAt:  report.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, result)
}
