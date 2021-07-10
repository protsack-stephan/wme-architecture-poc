package schema

// Message kafka message event
type Message struct {
	Event *Event `json:"event"`
}
