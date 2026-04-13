package handlers

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"

	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handlers) ScanSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	h.mu.Lock()
	h.scanClients = append(h.scanClients, conn)
	h.mu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Handlers) Broadcast(payload []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < len(h.scanClients); i++ {
		err := h.scanClients[i].WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			h.scanClients[i].Close()
			h.scanClients = append(h.scanClients[:i], h.scanClients[i+1:]...)
			i--
		}
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
			if n != nil {
				h.Broadcast([]byte(n.Extra))
			}
		}
	}
}
