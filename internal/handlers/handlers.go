package handlers

import (
	"encoding/json"
	"net/http"
	"storeSystem/internal/database"
	"storeSystem/internal/models"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	itemStore   *database.ItemStore
	workerStore *database.WorkerStore
}

func NewHandlers(
	itemStore *database.ItemStore,
	workerStore *database.WorkerStore,
) *Handlers {
	return &Handlers{
		itemStore:   itemStore,
		workerStore: workerStore,
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *Handlers) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.itemStore.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения товаров")
		return
	}
	respondWithJSON(w, http.StatusOK, items)
}

func (h *Handlers) GetItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id товара")
		return
	}
	item, err := h.itemStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}

func (h *Handlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	var input models.CreateItemInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	if strings.TrimSpace(input.Article) == "" {
		respondWithError(w, http.StatusBadRequest, "article item должен присутствовать")
		return
	}

	item, err := h.itemStore.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, item)
}

func (h *Handlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/items/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id задачи")
		return
	}

	var input models.UpdateItemInput

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if input.Article != nil && strings.TrimSpace(*input.Article) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок обязательный")
	}
	item, err := h.itemStore.Update(id, input)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}
