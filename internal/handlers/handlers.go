package handlers

import (
	"encoding/json"
	"net/http"
	"storeSystem/internal/database"
	"storeSystem/internal/minio"
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
	shipmentListStore *database.ShipmentListStore
	shipmentStore     *database.ShipmentStore
	stockStore        *database.StockStore
	counterpartyStore *database.CounterpartyStore
	minioService      minio.Client
	reportStore       *database.ReportStore
	mu                sync.Mutex
	clients           map[int]map[string][]Subscription
}

func NewHandlers(
	itemStore *database.ItemStore,
	workerStore *database.WorkerStore,
	deliveryListStore *database.DeliveryListStore,
	deliveryStore *database.DeliveryStore,
	shipmentListStore *database.ShipmentListStore,
	shipmentStore *database.ShipmentStore,
	counterpartyStore *database.CounterpartyStore,
	stockStore *database.StockStore,
	minioService minio.Client,
	reportStore *database.ReportStore,
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
