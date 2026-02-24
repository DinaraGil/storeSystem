package handlers

import (
	"encoding/json"
	"net/http"
	"storeSystem/internal/models"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetAllDeliveryLists(w http.ResponseWriter, r *http.Request) {
	delLists, err := h.deliveryListStore.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения листов поставок")
		return
	}
	respondWithJSON(w, http.StatusOK, delLists)
}

func (h *Handlers) GetDeliveryListByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id листа поставки")
		return
	}
	delList, err := h.deliveryListStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, delList)
}

func (h *Handlers) CreateDeliveryList(w http.ResponseWriter, r *http.Request) {
	var input models.CreateDeliveryListInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	if strings.TrimSpace(input.Article) == "" {
		respondWithError(w, http.StatusBadRequest, "article должен присутствовать")
		return
	}

	delList, err := h.deliveryListStore.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, delList)
}
