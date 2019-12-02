package jsonapi

import (
	"reflect"
	"testing"
)

type Sample struct {
	// resource ID
	ID string `jsonapi:"primary,samples"`

	// basic types on attributes
	Float64 float64 `jsonapi:"attribute,float64"`
	Int     int     `jsonapi:"attribute,int"`
	String  string  `jsonapi:"attribute,string"`

	// basic types on meta
	MetaString  string  `jsonapi:"meta,string"`
	MetaFloat64 float64 `jsonapi:"meta,float64"`
	MetaInt     int     `jsonapi:"meta,int"`

	// custom type
	CustomStructPtr *CustomNullableString `jsonapi:"attribute,custom_struct_ptr"`

	// nested struct
	Nested SampleNested `jsonapi:"attribute,nested"`

	// embedded struct
	Embedded

	// inferred tags
	Default         string
	DefaultWithName string

	// ignored field
	IgnoredField string `jsonapi:"-"`
}

type CustomNullableString struct {
	String string
	Valid  bool
}

type Embedded struct {
	ID             string `jsonapi:"primary,embeddeds"`
	EmbeddedString string `jsonapi:"attribute,embedded_string"`
}

type SampleNested struct {
	ID           string `jsonapi:"primary,sample_nesteds"`
	NestedString string `jsonapi:"attribute,nested_string"`
}

// func TestUnmarshalStruct(t *testing.T) {
// 	s := Sample{}
// 	input := []byte(`{
// 		"data": {
// 			"id": "someID",
// 			"type": "samples",
// 			"attributes": {
// 				"nested": {
// 					"nested_string": "hello world!"
// 				}
// 			}
// 		}
// 	}`)
// 	if err := Unmarshal(input, &s); err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	if s.ID != "someID" {
// 		t.Errorf("ID was incorrect, got: %v, want: %v.", s.ID, "someID")
// 	}
// 	if s.Nested.NestedString != "hello world!" {
// 		t.Errorf("Nested.NestedString was incorrect, got: %v, want: %v.", s.Nested.NestedString, "hello world!")
// 	}
// }

func TestUnmarshalEmbeddedStruct(t *testing.T) {
	input := []byte(`{
		"data": {
			"id": "someID",
			"type": "samples",
			"attributes": {
				"string": "hello",
				"embedded_string": "world!"
			}
		}
	}`)
	s := Sample{}
	if err := Unmarshal(input, &s); err != nil {
		t.Errorf(err.Error())
	}
	if s.String != "hello" {
		t.Errorf("String was incorrect, got: %v, want: %v.", s.String, "hello")
	}
	if s.EmbeddedString != "world!" {
		t.Errorf("EmbeddedString was incorrect, got: %v, want: %v.", s.EmbeddedString, "world!")
	}
}

func TestUnmarshalCustomType(t *testing.T) {
	// TODO
}

func TestUnmarshalCustomTypePtr(t *testing.T) {
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
			"custom_struct_ptr": "hello world!"
		}
	}
}`)
	s1 := Sample{}
	if err := Unmarshal(inputWithValue, &s1); err != nil {
		t.Errorf(err.Error())
	}
	if !s1.CustomStructPtr.Valid {
		t.Errorf("CustomStructPtr.Valid was incorrect, got: %v, want: %v.", s1.CustomStructPtr.Valid, true)
	}
	if s1.CustomStructPtr.String != "hello world!" {
		t.Errorf("CustomStructPtr.String was incorrect, got: %v, want: %v.", s1.CustomStructPtr.String, "hello world!")
	}

	inputWithEmptyValue := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"custom_struct_ptr": ""
		}
	}
}`)
	s2 := Sample{}
	if err := Unmarshal(inputWithEmptyValue, &s2); err != nil {
		t.Errorf(err.Error())
	}
	if !s2.CustomStructPtr.Valid {
		t.Errorf("CustomStructPtr.Valid was incorrect, got: %v, want: %v.", s2.CustomStructPtr.Valid, true)
	}
	if s2.CustomStructPtr.String != "" {
		t.Errorf("CustomStructPtr.String was incorrect, got: %v, want: %v.", s2.CustomStructPtr.String, "")
	}

	inputWithNullValue := []byte(`{
	"data": {
		"id": "someID",
		"type": "examples",
		"attributes": {
			"custom_struct_ptr": null
		}
	}
}`)
	s3 := Sample{}
	if err := Unmarshal(inputWithNullValue, &s3); err != nil {
		t.Errorf(err.Error())
	}
	if s3.CustomStructPtr.Valid {
		t.Errorf("CustomStructPtr.Valid was incorrect, got: %v, want: %v.", s3.CustomStructPtr.Valid, false)
	}
	if s3.CustomStructPtr.String != "" {
		t.Errorf("CustomStructPtr.String was incorrect, got: %v, want: %v.", s3.CustomStructPtr.String, "")
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
	s4 := Sample{}
	if err := Unmarshal(inputWithWithoutValue, &s4); err != nil {
		t.Errorf(err.Error())
	}
	if s4.CustomStructPtr != nil {
		t.Errorf("CustomStructPtr was incorrect, got: %+v, want: %v.", s4.CustomStructPtr, nil)
	}
}

func TestUnmarshal(t *testing.T) {
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
			"float64": 99.2486135148,
			"int": 5845,
			"string": "someString"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	s := Sample{}
	if err := Unmarshal(input, &s); err != nil {
		t.Errorf(err.Error())
	}
	if s.ID != "someID" {
		t.Errorf("ID was incorrect, got: %s, want: %s.", s.ID, "someID")
	}
	if s.Int != 99 {
		t.Errorf("Int was incorrect, got: %d, want: %d.", s.Int, 99)
	}
	if s.Float64 != 3.14159265359 {
		t.Errorf("Float64 was incorrect, got: %f, want: %f.", s.Float64, 3.14159265359)
	}
	if s.String != "someString" {
		t.Errorf("String was incorrect, got: %s, want: %s.", s.String, "someString")
	}
	if s.MetaInt != 5845 {
		t.Errorf("MetaInt was incorrect, got: %d, want: %d.", s.MetaInt, 5845)
	}
	if s.MetaString != "someString" {
		t.Errorf("MetaString was incorrect, got: %s, want: %s.", s.MetaString, "someString")
	}
	if s.MetaFloat64 != 99.2486135148 {
		t.Errorf("MetaFloat64 was incorrect, got: %f, want: %f.", s.MetaFloat64, 99.2486135148)
	}
}
