package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"storeSystem/internal/models"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetAllDeliveries(w http.ResponseWriter, r *http.Request) {
	del, err := h.deliveryStore.GetAll()

	fmt.Println(del, err)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения поставок")
		return
	}
	respondWithJSON(w, http.StatusOK, del)
}

func (h *Handlers) GetDeliveryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id поставки")
		return
	}
	del, err := h.deliveryStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, del)
}

func (h *Handlers) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	var input models.CreateDeliveryInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	if strings.TrimSpace(input.Status) == "" {
		respondWithError(w, http.StatusBadRequest, "status должен присутствовать")
		return
	}

	del, err := h.deliveryStore.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, del)
}

func (h *Handlers) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id поставки")
		return
	}

	var input models.UpdateDeliveryInput

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if input.Status != nil && strings.TrimSpace(*input.Status) == "" {
		respondWithError(w, http.StatusBadRequest, "Статус обязательный")
	}
	delivery, err := h.deliveryStore.Update(id, input)
	fmt.Println(delivery, err)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, delivery)
}
