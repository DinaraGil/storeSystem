package handlers

import (
	"encoding/json"
	"fmt"
	"storeSystem/internal/models"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"

	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handlers) ScanSocket(w http.ResponseWriter, r *http.Request) {
	deliveryStr := chi.URLParam(r, "delivery_id")
	scannerStr := chi.URLParam(r, "scanner_id")

	deliveryID, err := strconv.Atoi(deliveryStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный delivery_id")
		return
	}

	scannerID, err := strconv.Atoi(scannerStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный scanner_id")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	sub := scanSubscription{
		DeliveryID: deliveryID,
		Conn:       conn,
	}

	h.mu.Lock()
	h.scanClients[scannerID] = append(h.scanClients[scannerID], sub)
	h.mu.Unlock()

	defer func() {
		h.removeScanClient(scannerID, conn)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Handlers) removeScanClient(scannerID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.scanClients[scannerID]
	for i, sub := range clients {
		if sub.Conn == conn {
			h.scanClients[scannerID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	if len(h.scanClients[scannerID]) == 0 {
		delete(h.scanClients, scannerID)
	}
}

func (h *Handlers) BroadcastToScanner(scannerID int, payload []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.scanClients[scannerID]
	for i := 0; i < len(clients); i++ {
		err := clients[i].Conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			clients[i].Conn.Close()
			clients = append(clients[:i], clients[i+1:]...)
			i--
		}
	}

	if len(clients) == 0 {
		delete(h.scanClients, scannerID)
	} else {
		h.scanClients[scannerID] = clients
	}
}

func (h *Handlers) ListenEvents(dsn string) {
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, nil)

	err := listener.Listen("event_channel")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case n := <-listener.Notify:
			if n == nil {
				continue
			}

			var evt models.Event
			if err := json.Unmarshal([]byte(n.Extra), &evt); err != nil {
				continue
			}

			subs := h.getScannerSubscriptions(evt.Scanner)

			for _, sub := range subs {
				err := h.deliveryListStore.ProcessScannerEvent(sub.DeliveryID, evt)
				if err != nil {
					fmt.Println(err)
					// можно логировать, а клиенту отправить ошибку
				}
			}

			h.BroadcastToScanner(evt.Scanner, []byte(n.Extra))
		}
	}
}

func (h *Handlers) getScannerSubscriptions(scannerID int) []scanSubscription {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.scanClients[scannerID]
	result := make([]scanSubscription, len(clients))
	copy(result, clients)
	return result
}
