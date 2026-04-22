package handlers

import (
	"net/http"
)

func (h *Handlers) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.stockStore.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения остатков")
		return
	}
	respondWithJSON(w, http.StatusOK, stocks)
}
