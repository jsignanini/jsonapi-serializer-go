package jsonapi

import (
	"bytes"
	"reflect"
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
	RegisterMarshaler(reflect.TypeOf(&CustomNullableString{}), func(s map[string]interface{}, memberName string, value reflect.Value) {
		if value.IsNil() {
			return
		}
		cns := value.Interface().(*CustomNullableString)
		if !cns.Valid {
			s[memberName] = nil
			return
		}
		s[memberName] = cns.String
	})

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
	b, err := Marshal(&s, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	if bytes.Compare(input, b) != 0 {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", string(input), string(b))
	}
}

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
	if b, err := Marshal(&testTrue, nil); err != nil {
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
	if b, err := Marshal(&testFalse, nil); err != nil {
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
	if b, err := Marshal(&testTrue, nil); err != nil {
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
	if b, err := Marshal(&testFalse, nil); err != nil {
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
	if b, err := Marshal(&testNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedNil, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedNil), string(b))
		}
	}
}

func TestMarshalCustomTypePtr(t *testing.T) {
	type CustomNullableString struct {
		String string
		Valid  bool
	}
	type TestCustomType struct {
		ID  string                `jsonapi:"primary,test_custom_types"`
		Foo *CustomNullableString `jsonapi:"attribute,bar"`
	}
	RegisterMarshaler(reflect.TypeOf(&CustomNullableString{}), func(s map[string]interface{}, memberName string, value reflect.Value) {
		if value.IsNil() {
			return
		}
		cns := value.Interface().(*CustomNullableString)
		if !cns.Valid {
			s[memberName] = nil
			return
		}
		s[memberName] = cns.String
	})

	t1 := TestCustomType{
		ID: "someID",
		Foo: &CustomNullableString{
			String: "hello world!",
			Valid:  true,
		},
	}
	expectedValidString := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_custom_types",
		"attributes": {
			"bar": "hello world!"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&t1, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedValidString, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedValidString), string(b))
		}
	}

	t2 := TestCustomType{
		ID: "someID",
		Foo: &CustomNullableString{
			String: "",
			Valid:  true,
		},
	}
	expectedValidEmptyString := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_custom_types",
		"attributes": {
			"bar": ""
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&t2, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedValidEmptyString, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedValidEmptyString), string(b))
		}
	}

	t3 := TestCustomType{
		ID: "someID",
		Foo: &CustomNullableString{
			String: "something",
			Valid:  false,
		},
	}
	expectedNull := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_custom_types",
		"attributes": {
			"bar": null
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&t3, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedNull, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedNull), string(b))
		}
	}

	t4 := TestCustomType{
		ID: "someID",
	}
	expectedNil := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_custom_types"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&t4, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedNil, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedNil), string(b))
		}
	}
}

func TestMarshalCompound(t *testing.T) {
	type TestCompound struct {
		ID  string `jsonapi:"primary,test_compounds"`
		Foo string `jsonapi:"attribute,bar"`
	}
	tcs := []*TestCompound{
		{
			ID:  "someID1",
			Foo: "hello",
		},
		{
			ID:  "someID2",
			Foo: "world!",
		},
	}
	expected := []byte(`{
	"data": [
		{
			"id": "someID1",
			"type": "test_compounds",
			"attributes": {
				"bar": "hello"
			}
		},
		{
			"id": "someID2",
			"type": "test_compounds",
			"attributes": {
				"bar": "world!"
			}
		}
	],
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&tcs, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestMarshalString(t *testing.T) {
	type TestString struct {
		ID  string `jsonapi:"primary,test_strings"`
		Foo string `jsonapi:"attribute,bar"`
	}
	ts := TestString{
		ID:  "someID",
		Foo: "hello world!",
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_strings",
		"attributes": {
			"bar": "hello world!"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&ts, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestMarshalStringPtr(t *testing.T) {
	type TestStringPtr struct {
		ID  string  `jsonapi:"primary,test_strings"`
		Foo *string `jsonapi:"attribute,bar"`
	}
	s := "hello world!"
	test := TestStringPtr{
		ID:  "someID",
		Foo: &s,
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_strings",
		"attributes": {
			"bar": "hello world!"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&test, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
	testNil := TestStringPtr{
		ID: "someID",
	}
	expectedNil := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_strings"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&testNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expectedNil, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expectedNil), string(b))
		}
	}
}

func TestMarshalInt(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo int    `jsonapi:"attribute,bar"`
	}
	seven := TestInt{
		ID:  "someID",
		Foo: 7,
	}
	want := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints",
		"attributes": {
			"bar": 7
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&seven, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, want) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(want))
		}
	}
}

func TestMarshalIntPtr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *int   `jsonapi:"attribute,bar"`
	}
	seven := 7
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	want := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints",
		"attributes": {
			"bar": 7
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtr, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, want) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(want))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	want := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtr, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, want) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(want))
		}
	}
}
