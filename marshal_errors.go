package jsonapi

import "encoding/json"

func MarshalErrors(p *MarshalParams, errs ...Error) ([]byte, error) {
	ndp := NewDocumentParams{}
	if p != nil {
		ndp.Links = p.Links
		ndp.Meta = p.Meta
	}
	document := NewDocument(&ndp)
	document.Errors = errs
	return json.MarshalIndent(&document, jsonPrefix, jsonIndent)
}
