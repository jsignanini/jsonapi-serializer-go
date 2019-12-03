package jsonapi

import (
	"bytes"
	"testing"
)

func TestMarshalBool(t *testing.T) {
	type TestBool struct {
		ID     string `jsonapi:"primary,test_bools"`
		IsTrue bool   `jsonapi:"attribute,is_true"`
	}
	testTrue := TestBool{
		ID:     "someID",
		IsTrue: true,
	}
	expectedTrue := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_bools",
		"attributes": {
			"is_true": true
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&testTrue); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedTrue, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedTrue), string(b))
		}
	}

	testFalse := TestBool{
		ID:     "someID",
		IsTrue: false,
	}
	expectedFalse := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_bools",
		"attributes": {
			"is_true": false
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&testFalse); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedFalse, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedFalse), string(b))
		}
	}
}

func TestMarshalBoolPtr(t *testing.T) {
	truthy := true
	falsy := false
	type TestBool struct {
		ID     string `jsonapi:"primary,test_bools"`
		IsTrue *bool  `jsonapi:"attribute,is_true"`
	}
	testTrue := TestBool{
		ID:     "someID",
		IsTrue: &truthy,
	}
	expectedTrue := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_bools",
		"attributes": {
			"is_true": true
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&testTrue); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedTrue, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedTrue), string(b))
		}
	}

	testFalse := TestBool{
		ID:     "someID",
		IsTrue: &falsy,
	}
	expectedFalse := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_bools",
		"attributes": {
			"is_true": false
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&testFalse); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedFalse, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedFalse), string(b))
		}
	}

	testNil := TestBool{
		ID:     "someID",
		IsTrue: nil,
	}
	expectedNil := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_bools"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&testNil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedNil, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedNil), string(b))
		}
	}
}

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
			"embedded_string": "",
			"float64": 3.14159265359,
			"int": 99,
			"nested": {
				"nested_string": ""
			},
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
