package jsonapi

import (
	"database/sql"
	"reflect"
	"testing"
)

// TODO
func TestUnmarshalEmbeddedStruct(t *testing.T) {
	// 	type Embedded struct {
	// 		ID  string `jsonapi:"primary,examples"`
	// 		Foo string `jsonapi:"attribute,foo"`
	// 	}
	// 	type ExampleWithEmbedded struct {
	// 		ID       string `jsonapi:"primary,examples"`
	// 		Bar      string `jsonapi:"attribute,bar"`
	// 		Embedded `jsonapi:"embedded,test"`
	// 	}
	// 	input := []byte(`{
	// 	"data": {
	// 		"id": "someID",
	// 		"type": "examples",
	// 		"attributes": {
	// 			"foo": "hello",
	// 			"bar": "world!"
	// 		}
	// 	}
	// }`)
	// 	e := ExampleWithEmbedded{}
	// 	if err := Unmarshal(input, &e); err != nil {
	// 		t.Errorf(err.Error())
	// 	}
	// 	fmt.Printf("%+v\n", e)
}

func TestUnmarshalCustomType(t *testing.T) {
	// TODO
}

func TestUnmarshalCustomTypePtr(t *testing.T) {
	type CustomNullableString struct {
		String string
		Valid  bool
	}
	type Example struct {
		ID     string                `jsonapi:"primary,examples"`
		Custom *CustomNullableString `jsonapi:"attribute,custom"`
	}
	RegisterUnmarshaler(reflect.TypeOf(&CustomNullableString{}), func(v interface{}, rv *reflect.Value) {
		ns := &CustomNullableString{}
		if v != nil {
			ns.Valid = true
			ns.String = v.(string)
		} else {
			ns.Valid = false
		}
		rv.Set(reflect.ValueOf(ns))
	})

	inputWithValue := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"custom": "hello world!"
		}
	}
}`)
	e1 := Example{}
	if err := Unmarshal(inputWithValue, &e1); err != nil {
		t.Errorf(err.Error())
	}
	if !e1.Custom.Valid {
		t.Errorf("Custom.Valid was incorrect, got: %v, want: %v.", e1.Custom.Valid, true)
	}
	if e1.Custom.String != "hello world!" {
		t.Errorf("Custom.String was incorrect, got: %v, want: %v.", e1.Custom.String, "hello world!")
	}

	inputWithEmptyValue := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"custom": ""
		}
	}
}`)
	e2 := Example{}
	if err := Unmarshal(inputWithEmptyValue, &e2); err != nil {
		t.Errorf(err.Error())
	}
	if !e2.Custom.Valid {
		t.Errorf("Custom.Valid was incorrect, got: %v, want: %v.", e2.Custom.Valid, true)
	}
	if e2.Custom.String != "" {
		t.Errorf("Custom.String was incorrect, got: %v, want: %v.", e2.Custom.String, "")
	}

	inputWithNullValue := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"custom": null
		}
	}
}`)
	e3 := Example{}
	if err := Unmarshal(inputWithNullValue, &e3); err != nil {
		t.Errorf(err.Error())
	}
	if e3.Custom.Valid {
		t.Errorf("Custom.Valid was incorrect, got: %v, want: %v.", e3.Custom.Valid, false)
	}
	if e3.Custom.String != "" {
		t.Errorf("Custom.String was incorrect, got: %v, want: %v.", e3.Custom.String, "")
	}

	inputWithWithoutValue := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"foo": "bar"
		}
	}
}`)
	e4 := Example{}
	if err := Unmarshal(inputWithWithoutValue, &e4); err != nil {
		t.Errorf(err.Error())
	}
	if e4.Custom != nil {
		t.Errorf("Custom was incorrect, got: %+v, want: %v.", e4.Custom, nil)
	}
}

func TestUnmarshal(t *testing.T) {
	type NullString struct {
		sql.NullString
	}
	RegisterUnmarshaler(reflect.TypeOf(&NullString{}), func(v interface{}, rv *reflect.Value) {
		ns := &NullString{}
		if v != nil {
			ns.Valid = true
			ns.String = v.(string)
		} else {
			ns.Valid = false
		}
		rv.Set(reflect.ValueOf(ns))
	})
	type Example struct {
		ID string `jsonapi:"primary,examples"`

		FooInt     int     `jsonapi:"attribute,foo_int"`
		FooFloat64 float64 `jsonapi:"attribute,foo_float"`
		FooString  string  `jsonapi:"attribute,foo_string"`

		FooMeta string `jsonapi:"meta,foo"`

		FooCustom *NullString `jsonapi:"attribute,foo_custom"`
	}
	input := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"foo_float": 3.14159265359,
			"foo_int": 99,
			"foo_string": "someString",
			"foo_custom": "hello world!"
		},
		"meta": {
			"foo": "bar"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	e := Example{}
	if err := Unmarshal(input, &e); err != nil {
		t.Errorf(err.Error())
	}
	if e.ID != "someID" {
		t.Errorf("ID was incorrect, got: %s, want: %s.", e.ID, "someID")
	}
	if e.FooInt != 99 {
		t.Errorf("Int was incorrect, got: %d, want: %d.", e.FooInt, 99)
	}
	if e.FooFloat64 != 3.14159265359 {
		t.Errorf("Float was incorrect, got: %f, want: %f.", e.FooFloat64, 3.14159265359)
	}
	if e.FooString != "someString" {
		t.Errorf("String was incorrect, got: %s, want: %s.", e.FooString, "someString")
	}
	if e.FooMeta != "bar" {
		t.Errorf("MetaString was incorrect, got: %s, want: %s.", e.FooMeta, "bar")
	}
	if e.FooCustom.String != "hello world!" {
		t.Errorf("MetaString was incorrect, got: %s, want: %s.", e.FooCustom.String, "hello world!")
	}
}
