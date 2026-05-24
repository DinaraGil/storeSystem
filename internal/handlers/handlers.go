//go:build !nocover
// +build !nocover

package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

type scanSubscription struct {
	DeliveryID int
	WorkerID   int
	Conn       Conn
}
type Handlers struct {
	itemStore         ItemStore
	workerStore       WorkerStore
	deliveryListStore DeliveryListStore
	deliveryStore     DeliveryStore
	shipmentListStore ShipmentListStore
	shipmentStore     ShipmentStore
	stockStore        StockStore
	counterpartyStore CounterpartyStore
	minioService      MinioService
	reportStore       ReportStore
	mu                sync.Mutex
	clients           map[int]map[string][]Subscription
}

func NewHandlers(
	itemStore ItemStore,
	workerStore WorkerStore,
	deliveryListStore DeliveryListStore,
	deliveryStore DeliveryStore,
	shipmentListStore ShipmentListStore,
	shipmentStore ShipmentStore,
	counterpartyStore CounterpartyStore,
	stockStore StockStore,
	minioService MinioService,
	reportStore ReportStore,
) *Handlers {
	return &Handlers{
		itemStore:         itemStore,
		workerStore:       workerStore,
		deliveryListStore: deliveryListStore,
		deliveryStore:     deliveryStore,
		shipmentListStore: shipmentListStore,
		shipmentStore:     shipmentStore,
		counterpartyStore: counterpartyStore,
		stockStore:        stockStore,
		minioService:      minioService,
		reportStore:       reportStore,
		clients:           make(map[int]map[string][]Subscription),
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
