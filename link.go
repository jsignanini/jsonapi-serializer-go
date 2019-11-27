package jsonapi

type LinkString string

type LinkObject struct {
	HRef string `json:"href"`
	Meta Meta   `json:"meta,omitempty"`
}
