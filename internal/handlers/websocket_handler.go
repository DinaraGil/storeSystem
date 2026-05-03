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

type Subscription struct {
	ObjectID int // delivery_id / shipment_id
	WorkerID int
	Conn     *websocket.Conn
	Handler  func(evt models.Event) (any, error)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handlers) Socket(w http.ResponseWriter, r *http.Request) {
	objectType := chi.URLParam(r, "object_type") // delivery / shipment
	objectIDStr := chi.URLParam(r, "object_id")
	scannerStr := chi.URLParam(r, "scanner_id")

	claims, ok := GetUserClaimsFromContext(r.Context())
	if !ok {
		return
	}

	objectID, err := strconv.Atoi(objectIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid object_id")
		return
	}

	scannerID, err := strconv.Atoi(scannerStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid scanner_id")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	sub := Subscription{
		ObjectID: objectID,
		WorkerID: claims.UserID,
		Conn:     conn,
	}

	h.mu.Lock()
	if h.clients[scannerID] == nil {
		h.clients[scannerID] = make(map[string][]Subscription)
	}
	h.clients[scannerID][objectType] = append(h.clients[scannerID][objectType], sub)
	h.mu.Unlock()

	defer func() {
		h.removeClient(scannerID, objectType, conn)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Handlers) removeClient(scannerID int, objectType string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	list := h.clients[scannerID][objectType]

	for i, sub := range list {
		if sub.Conn == conn {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}

	if len(list) == 0 {
		delete(h.clients[scannerID], objectType)
	} else {
		h.clients[scannerID][objectType] = list
	}
}

func (h *Handlers) Broadcast(scannerID int, objectType string, payload []byte) {
	h.mu.Lock()
	clients := append([]Subscription(nil), h.clients[scannerID][objectType]...)
	h.mu.Unlock()

	var alive []Subscription

	for _, sub := range clients {
		err := sub.Conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			sub.Conn.Close()
			continue
		}
		alive = append(alive, sub)
	}

	h.mu.Lock()
	if len(alive) == 0 {
		delete(h.clients[scannerID], objectType)
	} else {
		h.clients[scannerID][objectType] = alive
	}
	h.mu.Unlock()
}

func (h *Handlers) ListenEvents(dsn string) {
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, nil)

	_ = listener.Listen("event_channel")

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

			h.processEvent(evt)
		}
	}
}

func (h *Handlers) processEvent(evt models.Event) {
	h.mu.Lock()
	scannerClients := h.clients[evt.Scanner]
	h.mu.Unlock()

	for objectType, subs := range scannerClients {

		for _, sub := range subs {

			var result any
			var err error

			switch objectType {
			case "delivery":
				result, err = h.deliveryListStore.ProcessScannerEvent(sub.ObjectID, evt, sub.WorkerID)

			case "shipment":
				result, err = h.shipmentListStore.ProcessScannerEvent(sub.ObjectID, evt, sub.WorkerID)
			}

			if err != nil {
				fmt.Println(err)
				errorPayload, _ := json.Marshal(map[string]string{
					"error": err.Error(),
				})
				_ = sub.Conn.WriteMessage(websocket.TextMessage, errorPayload)
				continue
			}

			payload, _ := json.Marshal(result)

			_ = sub.Conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}
