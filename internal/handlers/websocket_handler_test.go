package handlers

import (
	"fmt"
	"storeSystem/internal/database/mocks"
	"storeSystem/internal/models"
	"testing"
)

type mockConn struct {
	written [][]byte
	closed  bool
}

func (m *mockConn) WriteMessage(_ int, data []byte) error {
	m.written = append(m.written, data)
	return nil
}

func (m *mockConn) ReadMessage() (int, []byte, error) {
	return 0, nil, fmt.Errorf("closed")
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

func TestRemoveClient(t *testing.T) {
	h := &Handlers{
		clients: make(map[int]map[string][]Subscription),
	}

	conn1 := &mockConn{}
	conn2 := &mockConn{}

	h.clients[1] = map[string][]Subscription{
		"delivery": {
			{Conn: conn1},
			{Conn: conn2},
		},
	}

	h.removeClient(1, "delivery", conn1)

	if len(h.clients[1]["delivery"]) != 1 {
		t.Fatal("client was not removed")
	}

	if h.clients[1]["delivery"][0].Conn != conn2 {
		t.Fatal("wrong client removed")
	}
}

type mockDeliveryStore struct{}

func (m *mockDeliveryStore) ProcessScannerEvent(id int, evt models.Event, workerID int) (any, error) {
	return map[string]string{"ok": "delivery"}, nil
}

func TestProcessEvent_Delivery(t *testing.T) {
	conn := &mockConn{}

	h := &Handlers{
		clients:           make(map[int]map[string][]Subscription),
		deliveryListStore: &mocks.MockDeliveryListStore{},
		shipmentListStore: nil,
	}

	h.clients[1] = map[string][]Subscription{
		"delivery": {
			{
				ObjectID: 10,
				WorkerID: 1,
				Conn:     conn,
			},
		},
	}

	evt := models.Event{
		Scanner: 1,
	}

	h.processEvent(evt)

	if len(conn.written) == 0 {
		t.Fatal("no message sent")
	}
}
