package jsonapi

type Document struct {
	Data    Resource `json:"data,omitempty"`
	JSONAPI JSONAPI  `json:"jsonapi"`
	Meta    Meta     `json:"meta,omitempty"`
	// Errors
	// Links
	// Included
}
