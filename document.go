package jsonapi

type Document struct {
	Data *Resource `json:"data,omitempty"`
	document
}

type NewDocumentParams struct {
	documentParams
}

func NewDocument(p *NewDocumentParams) *Document {
	d := document{
		JSONAPI: &JSONAPIInformation{
			Version: "1.0",
		},
	}
	if p != nil {
		d.Links = p.Links
		d.Meta = p.Meta
	}
	return &Document{
		document: d,
	}
}

type NewCompoundDocumentParams struct {
	documentParams
}

type CompoundDocument struct {
	Data []*Resource `json:"data"`
	document
}

func NewCompoundDocument(p *NewCompoundDocumentParams) *CompoundDocument {
	d := document{
		JSONAPI: &JSONAPIInformation{
			Version: "1.0",
		},
	}
	if p != nil {
		d.Links = p.Links
		d.Meta = p.Meta
	}
	return &CompoundDocument{
		Data:     []*Resource{},
		document: d,
	}
}

type document struct {
	JSONAPI  *JSONAPIInformation `json:"jsonapi,omitempty"`
	Meta     *Meta               `json:"meta,omitempty"`
	Links    *Links              `json:"links,omitempty"`
	Errors   []Error             `json:"errors,omitempty"`
	Included []*Resource         `json:"included,omitempty"`
}

type documentParams struct {
	Links *Links
	Meta  *Meta
}
