package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"storeSystem/internal/helpers"
	"strconv"

	"github.com/google/uuid"
)

type ReportConfig struct {
	ReportType string `json:"report_type"`
	DateFrom   string `json:"date_from"`
	DateTo     string `json:"date_to"`
	//clientID   int
	//date       string
}

func (h *Handlers) GenerateStockReport(config *ReportConfig) (string, error) {
	stocks, err := h.stockStore.GetAll()
	if err != nil {
		return "", fmt.Errorf("ошибка получения остатков")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{"Артикул", "Количество", "Время изменения"}); err != nil {
		return "", err
	}

	for _, stock := range stocks {
		record := []string{
			stock.Article,
			strconv.Itoa(stock.Quantity),
			stock.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	filename := "stock_report_" + uuid.NewString() + ".csv"

	fileData := helpers.FileDataType{
		FileName: filename,
		Data:     buf.Bytes(),
	}

	link, err := h.minioService.CreateOne(fileData)
	if err != nil {
		return "", fmt.Errorf("unable to save the file: %w", err)
	}

	return link, nil
}

func (h *Handlers) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var config ReportConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var link string
	var err error

	switch config.ReportType {
	case "stock":
		link, err = h.GenerateStockReport(&config)
	default:
		err = fmt.Errorf("unknown report type")
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"link": link,
	})
}
