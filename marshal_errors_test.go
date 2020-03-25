package jsonapi

import (
	"bytes"
	"testing"
)

func TestMarshalErrors(t *testing.T) {
	// simple error
	simpleError := Error{
		ID:     "NOT_FOUND",
		Status: "404",
		Title:  "not-found",
	}
	simpleErrorExpected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"id": "NOT_FOUND",
			"status": "404",
			"title": "not-found"
		}
	]
}`)
	if b, err := MarshalErrors(nil, simpleError); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(simpleErrorExpected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(simpleErrorExpected), string(b))
		}
	}

	// error with links
	simpleErrorWithLinks := Error{
		ID:     "NOT_FOUND",
		Status: "404",
		Title:  "not-found",
		Links: Links{
			"about": "/errors/NOT_FOUND",
		},
	}
	simpleErrorWithLinksExpected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"id": "NOT_FOUND",
			"links": {
				"about": "/errors/NOT_FOUND"
			},
			"status": "404",
			"title": "not-found"
		}
	]
}`)
	if b, err := MarshalErrors(nil, simpleErrorWithLinks); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(simpleErrorWithLinksExpected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(simpleErrorWithLinksExpected), string(b))
		}
	}

	// error with meta
	simpleErrorWithMeta := Error{
		ID:     "NOT_FOUND",
		Status: "404",
		Title:  "not-found",
		Meta: Meta{
			"errors_documentation_url": "https://example.com",
		},
	}
	simpleErrorWithMetaExpected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"id": "NOT_FOUND",
			"status": "404",
			"title": "not-found",
			"meta": {
				"errors_documentation_url": "https://example.com"
			}
		}
	]
}`)
	if b, err := MarshalErrors(nil, simpleErrorWithMeta); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(simpleErrorWithMetaExpected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(simpleErrorWithMetaExpected), string(b))
		}
	}

	// error with meta and links
	simpleErrorWithMetaAndLinks := Error{
		ID:     "NOT_FOUND",
		Status: "404",
		Title:  "not-found",
		Links: Links{
			"about": "/errors/NOT_FOUND",
		},
		Meta: Meta{
			"errors_documentation_url": "https://example.com",
		},
	}
	simpleErrorWithMetaAndLinksExpected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"id": "NOT_FOUND",
			"links": {
				"about": "/errors/NOT_FOUND"
			},
			"status": "404",
			"title": "not-found",
			"meta": {
				"errors_documentation_url": "https://example.com"
			}
		}
	]
}`)
	if b, err := MarshalErrors(nil, simpleErrorWithMetaAndLinks); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(simpleErrorWithMetaAndLinksExpected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(simpleErrorWithMetaAndLinksExpected), string(b))
		}
	}

	// error with document meta and links
	simpleErrorWithDocumentMetaAndLinks := Error{
		ID:     "NOT_FOUND",
		Status: "404",
		Title:  "not-found",
	}
	simpleErrorWithDocumentMetaAndLinksExpected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"meta": {
		"errors_documentation_url": "https://example.com"
	},
	"links": {
		"about": "/errors/NOT_FOUND"
	},
	"errors": [
		{
			"id": "NOT_FOUND",
			"status": "404",
			"title": "not-found"
		}
	]
}`)
	if b, err := MarshalErrors(&MarshalParams{
		Links: &Links{
			"about": "/errors/NOT_FOUND",
		},
		Meta: &Meta{
			"errors_documentation_url": "https://example.com",
		},
	}, simpleErrorWithDocumentMetaAndLinks); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(simpleErrorWithDocumentMetaAndLinksExpected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(simpleErrorWithDocumentMetaAndLinksExpected), string(b))
		}
	}
}

func TestMarshalManyErrors(t *testing.T) {
	e1 := Error{
		ID:     "errorOne",
		Status: "404",
		Title:  "not-found",
	}
	e2 := Error{
		ID:     "errorTwo",
		Status: "500",
		Title:  "server-error",
	}
	expected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"id": "errorOne",
			"status": "404",
			"title": "not-found"
		},
		{
			"id": "errorTwo",
			"status": "500",
			"title": "server-error"
		}
	]
}`)
	if b, err := MarshalErrors(nil, e1, e2); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}
