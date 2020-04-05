package jsonapi

// Document is a JSON:API document top-level object.
// See https://jsonapi.org/format/#document-top-level.
type Document struct {
	Data *Resource `json:"data,omitempty"`
	document
}

// NewDocumentParams are the parameters describing a top-level document.
type NewDocumentParams struct {
	documentParams
}

// NewDocument generates a new top-level document object.
func NewDocument(p *NewDocumentParams) *Document {
	d := document{
		JSONAPI: &Information{
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

// CompoundDocument represents a compound document's top-level object.
// See https://jsonapi.org/format/#document-compound-documents.
type CompoundDocument struct {
	Data []*Resource `json:"data"`
	document
}

// NewCompoundDocumentParams are the parameters describing a new compound top-level document.
type NewCompoundDocumentParams struct {
	documentParams
}

// NewCompoundDocument generates a new compound top-level document object.
func NewCompoundDocument(p *NewCompoundDocumentParams) *CompoundDocument {
	d := document{
		JSONAPI: &Information{
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
	JSONAPI  *Information `json:"jsonapi,omitempty"`
	Meta     *Meta        `json:"meta,omitempty"`
	Links    *Links       `json:"links,omitempty"`
	Errors   []Error      `json:"errors,omitempty"`
	Included []*Resource  `json:"included,omitempty"`
}

type documentParams struct {
	Links *Links
	Meta  *Meta
}
