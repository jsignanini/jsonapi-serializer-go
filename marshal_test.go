package jsonapi

import (
	"bytes"
	"testing"
)

func TestMarshal(t *testing.T) {
	s := Sample{
		ID:         "someID",
		Int:        99,
		Float64:    3.14159265359,
		String:     "someString",
		MetaString: "bar",
	}
	input := []byte(`{
	"data": {
		"id": "someID",
		"type": "samples",
		"attributes": {
			"float64": 3.14159265359,
			"int": 99,
			"string": "someString"
		},
		"meta": {
			"float64": 0,
			"int": 0,
			"string": "bar"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	b, err := Marshal(&s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if bytes.Compare(input, b) != 0 {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", string(input), string(b))
	}
}
