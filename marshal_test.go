package jsonapi

import (
	"bytes"
	"testing"
)

func TestMarshal(t *testing.T) {
	type Example struct {
		ID string `jsonapi:"primary,examples"`

		FooInt     int     `jsonapi:"attribute,foo_int"`
		FooFloat64 float64 `jsonapi:"attribute,foo_float"`
		FooString  string  `jsonapi:"attribute,foo_string"`

		FooMeta string `jsonapi:"meta,foo"`
	}
	e := Example{
		ID:         "someID",
		FooInt:     99,
		FooFloat64: 3.14159265359,
		FooString:  "someString",
		FooMeta:    "bar",
	}

	input := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"foo_float": 3.14159265359,
			"foo_int": 99,
			"foo_string": "someString"
		},
		"meta": {
			"foo": "bar"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)

	b, err := Marshal(&e)
	if err != nil {
		t.Errorf(err.Error())
	}
	if bytes.Compare(input, b) != 0 {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", string(input), string(b))
	}
}
