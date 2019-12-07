package jsonapi

import "encoding/json"

func MarshalErrors(p *MarshalParams, errs ...Error) ([]byte, error) {
	document := NewDocument()
	document.Errors = errs

	// handle optional params
	if p != nil && p.Links != nil {
		document.Links = p.Links
	}
	if p != nil && p.Meta != nil {
		document.Meta = p.Meta
	}

	return json.MarshalIndent(&document, jsonPrefix, jsonIndent)
}
