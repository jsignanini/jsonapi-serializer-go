package jsonapi

type JSONAPIInformation struct {
	Version string `json:"version,omitempty"`
	Meta    Meta   `json:"meta,omitempty"`
}
