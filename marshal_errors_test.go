package jsonapi

import (
	"bytes"
	"testing"
)

func TestMarshalErrors(t *testing.T) {
	e := Error{
		ID:     "someID",
		Status: "404",
		Title:  "not-found",
	}
	expected := []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"id": "someID",
			"status": "404",
			"title": "not-found"
		}
	]
}`)
	if b, err := MarshalErrors(nil, e); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
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
