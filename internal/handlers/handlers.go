package handlers

import (
	"encoding/json"
	"net/http"
	"storeSystem/internal/database"
)

type Handlers struct {
	itemStore         *database.ItemStore
	workerStore       *database.WorkerStore
	deliveryListStore *database.DeliveryListStore
	deliveryStore     *database.DeliveryStore
}

func NewHandlers(
	itemStore *database.ItemStore,
	workerStore *database.WorkerStore,
	deliveryListStore *database.DeliveryListStore,
	deliveryStore *database.DeliveryStore,
) *Handlers {
	return &Handlers{
		itemStore:         itemStore,
		workerStore:       workerStore,
		deliveryListStore: deliveryListStore,
		deliveryStore:     deliveryStore,
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
