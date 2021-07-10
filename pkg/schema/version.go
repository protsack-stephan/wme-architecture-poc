package schema

// Version schema
type Version struct {
	Event      *Event `json:"event,omitempty"`
	Identifier int    `json:"identifier"`
	Comment    string `json:"comment,omitempty"`
}
