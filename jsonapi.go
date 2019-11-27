package jsonapi

var tagKey = "jsonapi"

type JSONAPI struct {
	Version string `json:"version,omitempty"`
	Meta    Meta   `json:"meta,omitempty"`
}
