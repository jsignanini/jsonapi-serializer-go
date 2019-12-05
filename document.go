package jsonapi

type document struct {
	JSONAPI JSONAPI `json:"jsonapi"`
	Meta    *Meta   `json:"meta,omitempty"`
	Links   *Links  `json:"links,omitempty"`
	// Errors
	// Included
}

type Document struct {
	Data *Resource `json:"data,omitempty"`
	document
}

type CompoundDocument struct {
	Data []*Resource `json:"data,omitempty"`
	document
}

func NewDocument() *Document {
	return &Document{
		Data: &Resource{
			Attributes: Attributes{},
			Meta:       Meta{},
		},
		document: document{
			JSONAPI: JSONAPI{
				Version: "1.0",
			},
		},
	}
}

func NewCompoundDocument() *CompoundDocument {
	return &CompoundDocument{
		Data: []*Resource{},
		document: document{
			JSONAPI: JSONAPI{
				Version: "1.0",
			},
		},
	}
}
