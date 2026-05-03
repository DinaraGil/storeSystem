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

func (h *Handlers) GetAllShipments(w http.ResponseWriter, r *http.Request) {
	sh, err := h.shipmentStore.GetAll()

	fmt.Println(sh, err)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения отгрузок")
		return
	}

	respondWithJSON(w, http.StatusOK, sh)
}

func (h *Handlers) GetShipmentByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id отгрузки")
		return
	}

	sh, err := h.shipmentStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sh)
}

func (h *Handlers) CreateShipment(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShipmentInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	if strings.TrimSpace(input.Status) == "" {
		respondWithError(w, http.StatusBadRequest, "status должен присутствовать")
		return
	}

	sh, err := h.shipmentStore.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, sh)
}

//func (h *Handlers) UpdateShipment(w http.ResponseWriter, r *http.Request) {
//	idStr := chi.URLParam(r, "id")
//	id, err := strconv.Atoi(idStr)
//
//	if err != nil {
//		respondWithError(w, http.StatusBadRequest, "Некорректный id отгрузки")
//		return
//	}
//
//	var input models.UpdateShipmentInput
//
//	err = json.NewDecoder(r.Body).Decode(&input)
//	if err != nil {
//		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
//		return
//	}
//
//	if input.Status != nil && strings.TrimSpace(*input.Status) == "" {
//		respondWithError(w, http.StatusBadRequest, "Статус обязателен")
//		return
//	}
//
//	shipment, err := h.shipmentStore.Update(id, input)
//
//	fmt.Println(shipment, err)
//
//	if err != nil {
//		if strings.Contains(err.Error(), "record not found") {
//			respondWithError(w, http.StatusNotFound, err.Error())
//		} else {
//			respondWithError(w, http.StatusInternalServerError, err.Error())
//		}
//		return
//	}
//
//	respondWithJSON(w, http.StatusOK, shipment)
//}

func (h *Handlers) GetShipmentListsByShipmentID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id отгрузки")
		return
	}

	lists, err := h.shipmentListStore.GetByShipmentID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, lists)
}

func (h *Handlers) CompleteShipment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id")
		return
	}

	err = h.shipmentStore.CompleteShipment(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "completed",
	})
}
