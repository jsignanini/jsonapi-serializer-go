package jsonapi

import "encoding/json"

type Error struct {
	ID     string            `json:"id,omitempty"`
	Links  Links             `json:"links,omitempty"`
	Status string            `json:"status,omitempty"`
	Code   string            `json:"code,omitempty"`
	Title  string            `json:"title,omitempty"`
	Detail string            `json:"detail,omitempty"`
	Source map[string]string `json:"source,omitempty"`
	Meta   Meta              `json:"meta,omitempty"`
}

func (e *Error) Error() string {
	if b, err := json.MarshalIndent(e, jsonPrefix, jsonIndent); err == nil {
		return string(b)
	}
	return ""
}