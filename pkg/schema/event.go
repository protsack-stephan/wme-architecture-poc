package schema

import "time"

const (
	EventTypeUpdate = "update"
	EventTypeCreate = "create"
	EventTypeDelete = "delete"
)

// Event representation of event in kaka topic
type Event struct {
	UUID string    `json:"uuid"`
	Type string    `json:"type"`
	Date time.Time `json:"date"`
}
