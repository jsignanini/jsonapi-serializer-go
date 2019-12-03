package jsonapi

type Document struct {
	Data    Resource `json:"data,omitempty"`
	JSONAPI JSONAPI  `json:"jsonapi"`
	Meta    Meta     `json:"meta,omitempty"`
	// Errors
	// Links
	// Included
}

func NewDocument() *Document {
	return &Document{
		JSONAPI: JSONAPI{
			Version: "1.0",
		},
		Data: Resource{
			Attributes: Attributes{},
			Meta:       Meta{},
		},
	}
}
