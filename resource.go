package jsonapi

type Resource struct {
	// Exception: The id member is not required when the resource object originates at the client and represents a new resource to be created on the server.
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`

	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
	// Links `json:"links,omitempty"`
	Meta Meta `json:"meta,omitempty"`
}
