package handlers

import (
	"encoding/json"
	"net/http"
	"storeSystem/internal/database"
	"sync"

	"github.com/gorilla/websocket"
)

type scanSubscription struct {
	DeliveryID int
	WorkerID   int
	Conn       *websocket.Conn
}
type Handlers struct {
	itemStore         *database.ItemStore
	workerStore       *database.WorkerStore
	deliveryListStore *database.DeliveryListStore
	deliveryStore     *database.DeliveryStore
	counterpartyStore *database.CounterpartyStore
	scanClients       map[int][]scanSubscription
	mu                sync.Mutex
}

func NewHandlers(
	itemStore *database.ItemStore,
	workerStore *database.WorkerStore,
	deliveryListStore *database.DeliveryListStore,
	deliveryStore *database.DeliveryStore,
	counterpartyStore *database.CounterpartyStore,
) *Handlers {
	return &Handlers{
		itemStore:         itemStore,
		workerStore:       workerStore,
		deliveryListStore: deliveryListStore,
		deliveryStore:     deliveryStore,
		counterpartyStore: counterpartyStore,
		scanClients:       make(map[int][]scanSubscription),
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
