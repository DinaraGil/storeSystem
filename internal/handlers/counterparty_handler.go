package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"storeSystem/internal/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetAllCounterparties(w http.ResponseWriter, r *http.Request) {
	counterparties, err := h.counterpartyStore.GetAll()
	fmt.Println(err)
	fmt.Println(counterparties)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения контрагентов")
		return
	}
	respondWithJSON(w, http.StatusOK, counterparties)
}

func (h *Handlers) GetCounterpartyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id контрагента")
		return
	}
	counterparty, err := h.counterpartyStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, counterparty)
}

func (h *Handlers) CreateCounterparty(w http.ResponseWriter, r *http.Request) {
	var input models.CreateCounterpartyInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	counterparty, err := h.counterpartyStore.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, counterparty)
}
