package jsonapi

// Information is an object that holds information about the JSON:API implementation
// a the top-level document.
// See https://jsonapi.org/format/#document-jsonapi-object.
type Information struct {
	Version string `json:"version,omitempty"`
	Meta    Meta   `json:"meta,omitempty"`
}
