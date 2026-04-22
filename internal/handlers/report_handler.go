package handlers

//
//import (
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"os"
//)
//
//type ReportConfig struct {
//	reportType string
//	dateFrom   string
//	dateTo     string
//	clientID   int
//	date       string
//}
//
//func (h *Handlers) GenerateStockReport(config *ReportConfig) error {
//	stocks, err := h.stockStore.GetAll()
//
//	if err != nil {
//		return fmt.Errorf("Ошибка получения остатков")
//	}
//
//	file, err := os.Create(config.reportType + config.)
//
//}
//
//func (h *Handlers) GenerateReport(w http.ResponseWriter, r *http.Request) {
//	var config *ReportConfig
//
//	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
//		respondWithError(w, http.StatusBadRequest, err.Error())
//		return
//	}
//
//	var generationError error
//	switch config.reportType {
//	case "stock":
//		generationError = GenerateStockReport()
//	default:
//		generationError = fmt.Errorf("unknown report type")
//	}
//	if generationError != nil {
//		respondWithError(w, http.StatusInternalServerError, generationError.Error())
//		return
//	}
//	respondWithJSON(w, http.StatusOK, "")
//}
