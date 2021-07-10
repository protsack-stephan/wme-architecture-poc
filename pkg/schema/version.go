package schema

import "time"

const (
	EventVersionCreate = "version_create"
)

// VersionEvent representation of event in kaka topic
type VersionEvent struct {
	UID     string    `json:"uid"`
	Type    string    `json:"type"`
	Date    time.Time `json:"date"`
	Payload *Version  `json:"payload"`
}

// Version schema
type Version struct {
	Identifier int    `json:"identifier"`
	Comment    string `json:"comment,omitempty"`
}
