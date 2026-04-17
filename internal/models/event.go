package models

type Event struct {
	EventId   int     `json:"event_id"`
	RfidId    *string `json:"rfid_id"`
	Article   *string `json:"article"`
	Scanner   int     `json:"scanner"`
	IsIn      *bool   `json:"is_in"`
	Error     *string `json:"error"`
	CreatedAt *string `json:"created_at"`
}
