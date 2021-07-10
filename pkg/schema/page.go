package schema

import "time"

// Page schema
type Page struct {
	Event        *Event         `json:"event,omitempty"`
	Name         string         `json:"name,omitempty"`
	Identifier   int            `json:"identifier"`
	Version      *Version       `json:"version,omitempty"`
	DateModified time.Time      `json:"date_modified,omitempty"`
	URL          string         `json:"url,omitempty"`
	IsPartOf     *Project       `json:"is_part_of,omitempty"`
	ArticleBody  []*ArticleBody `json:"article_body,omitempty"`
}

// ArticleBody content of the page
type ArticleBody struct {
	Text           string `json:"text"`
	EncodingFormat string `json:"encoding_format"`
}
