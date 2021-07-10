package schema

// Project schema
type Project struct {
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier"`
	URL        string `json:"url,omitempty"`
}
