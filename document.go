package jsonapi

type document struct {
	JSONAPI *JSONAPIInformation `json:"jsonapi,omitempty"`
	Meta    *Meta               `json:"meta,omitempty"`
	Links   *Links              `json:"links,omitempty"`
	Errors  []Error             `json:"errors,omitempty"`
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
		document: document{
			JSONAPI: &JSONAPIInformation{
				Version: "1.0",
			},
		},
	}
}

func NewCompoundDocument() *CompoundDocument {
	return &CompoundDocument{
		Data: []*Resource{},
		document: document{
			JSONAPI: &JSONAPIInformation{
				Version: "1.0",
			},
		},
	}
}
