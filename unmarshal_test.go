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

func TestUnmarshalBool(t *testing.T) {
	type TestBool struct {
		ID     string `jsonapi:"primary,test_bools"`
		IsTrue bool   `jsonapi:"attribute,is_true"`
	}

	expectedTrue := TestBool{}
	inputTrue := []byte(`{
		"data": {
			"id": "someID",
			"type": "test_bools",
			"attributes": {
				"is_true": true
			}
		}
	}`)
	if err := Unmarshal(inputTrue, &expectedTrue); err != nil {
		t.Errorf(err.Error())
	}
	if !expectedTrue.IsTrue {
		t.Errorf("IsTrue was incorrect, got: %t, want: %t.", expectedTrue.IsTrue, true)
	}

	expectedFalse := TestBool{}
	inputFalse := []byte(`{
		"data": {
			"id": "someID",
			"type": "test_bools",
			"attributes": {
				"is_true": true
			}
		}
	}`)
	if err := Unmarshal(inputFalse, &expectedFalse); err != nil {
		t.Errorf(err.Error())
	}
	if !expectedFalse.IsTrue {
		t.Errorf("IsTrue was incorrect, got: %t, want: %t.", expectedFalse.IsTrue, false)
	}
}

func TestUnmarshalBoolPtr(t *testing.T) {
	type TestBoolPtr struct {
		ID     string `jsonapi:"primary,test_bools"`
		IsTrue *bool  `jsonapi:"attribute,is_true"`
	}

	expectedTrue := TestBoolPtr{}
	inputTrue := []byte(`{
		"data": {
			"id": "someID",
			"type": "test_bools",
			"attributes": {
				"is_true": true
			}
		}
	}`)
	if err := Unmarshal(inputTrue, &expectedTrue); err != nil {
		t.Errorf(err.Error())
	}
	if expectedTrue.IsTrue == nil || !*expectedTrue.IsTrue {
		t.Errorf("IsTrue was incorrect, got: %v, want: %t.", expectedTrue.IsTrue, true)
	}

	expectedFalse := TestBoolPtr{}
	inputFalse := []byte(`{
		"data": {
			"id": "someID",
			"type": "test_bools",
			"attributes": {
				"is_true": false
			}
		}
	}`)
	if err := Unmarshal(inputFalse, &expectedFalse); err != nil {
		t.Errorf(err.Error())
	}
	if expectedFalse.IsTrue == nil || *expectedFalse.IsTrue {
		t.Errorf("IsTrue was incorrect, got: %v, want: %t.", expectedFalse.IsTrue, false)
	}

	expectedNil := TestBoolPtr{}
	inputNil := []byte(`{
		"data": {
			"id": "someID",
			"type": "test_bools"
		}
	}`)
	if err := Unmarshal(inputNil, &expectedNil); err != nil {
		t.Errorf(err.Error())
	}
	if expectedNil.IsTrue != nil {
		t.Errorf("IsTrue was incorrect, got: %v, want: %v.", expectedNil.IsTrue, nil)
	}
}

func TestUnmarshalInt(t *testing.T) {
	type TestUnmarshalInt struct {
		ID  string `jsonapi:"primary,ints"`
		Foo int    `jsonapi:"attribute,bar"`
	}
	valid := []byte(`{
		"data": {
			"id": "someID",
			"type": "ints",
			"attributes": {
				"bar": 99
			}
		}
	}`)
	validTest := &TestUnmarshalInt{}
	if err := Unmarshal(valid, validTest); err != nil {
		t.Errorf(err.Error())
	}
	if validTest.Foo != 99 {
		t.Errorf("expected int: %d, got: %d", 99, validTest.Foo)
	}
	negative := []byte(`{
		"data": {
			"id": "someID",
			"type": "ints",
			"attributes": {
				"bar": -28894
			}
		}
	}`)
	negativeTest := &TestUnmarshalInt{}
	if err := Unmarshal(negative, negativeTest); err != nil {
		t.Errorf(err.Error())
	}
	if negativeTest.Foo != -28894 {
		t.Errorf("expected int: %d, got: %d", -28894, negativeTest.Foo)
	}
}

func TestUnmarshalIntPtr(t *testing.T) {
	type TestUnmarshalInt struct {
		ID  string `jsonapi:"primary,ints"`
		Foo *int   `jsonapi:"attribute,bar"`
	}
	valid := []byte(`{
		"data": {
			"id": "someID",
			"type": "ints",
			"attributes": {
				"bar": 99
			}
		}
	}`)
	validTest := &TestUnmarshalInt{}
	if err := Unmarshal(valid, validTest); err != nil {
		t.Errorf(err.Error())
	}
	if *validTest.Foo != 99 {
		t.Errorf("expected *int: %d, got: %d", 99, *validTest.Foo)
	}
	negative := []byte(`{
		"data": {
			"id": "someID",
			"type": "ints",
			"attributes": {
				"bar": -28894
			}
		}
	}`)
	negativeTest := &TestUnmarshalInt{}
	if err := Unmarshal(negative, negativeTest); err != nil {
		t.Errorf(err.Error())
	}
	if *negativeTest.Foo != -28894 {
		t.Errorf("expected *int: %d, got: %d", -28894, *negativeTest.Foo)
	}
	empty := []byte(`{
		"data": {
			"id": "someID",
			"type": "ints"
		}
	}`)
	emptyTest := &TestUnmarshalInt{}
	if err := Unmarshal(empty, emptyTest); err != nil {
		t.Errorf(err.Error())
	}
	if emptyTest.Foo != nil {
		t.Errorf("expected *int: %v, got: %d", nil, *emptyTest.Foo)
	}
}

func TestUnmarshalInt8(t *testing.T) {
	// TODO
}

func TestUnmarshalInt8Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalInt16(t *testing.T) {
	// TODO
}

func TestUnmarshalInt16Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalInt32(t *testing.T) {
	// TODO
}

func TestUnmarshalInt32Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalInt64(t *testing.T) {
	// TODO
}

func TestUnmarshalInt64Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalUint(t *testing.T) {
	// TODO
}

func TestUnmarshalUintPtr(t *testing.T) {
	// TODO
}

func TestUnmarshalUint8(t *testing.T) {
	// TODO
}

func TestUnmarshalUint8Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalUint16(t *testing.T) {
	// TODO
}

func TestUnmarshalUint16Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalUint32(t *testing.T) {
	// TODO
}

func TestUnmarshalUint32Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalUint64(t *testing.T) {
	// TODO
}

func TestUnmarshalUint64Ptr(t *testing.T) {
	// TODO
}

func TestUnmarshalNestedStruct(t *testing.T) {
	s := Sample{}
	input := []byte(`{
		"data": {
			"id": "someID",
			"type": "samples",
			"attributes": {
				"nested": {
					"nested_string": "hello world!"
				}
			}
		}
	}`)
	if err := Unmarshal(input, &s); err != nil {
		t.Errorf(err.Error())
	}
	if s.ID != "someID" {
		t.Errorf("ID was incorrect, got: %v, want: %v.", s.ID, "someID")
	}
	if s.Nested.NestedString != "hello world!" {
		t.Errorf("Nested.NestedString was incorrect, got: %v, want: %v.", s.Nested.NestedString, "hello world!")
	}
}

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
	RegisterUnmarshaler(reflect.TypeOf(&CustomNullableString{}), func(v interface{}, value reflect.Value) {
		ns := &CustomNullableString{}
		if v != nil {
			ns.Valid = true
			ns.String = v.(string)
		} else {
			ns.Valid = false
		}
		value.Set(reflect.ValueOf(ns))
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

func TestUnmarshalMany(t *testing.T) {
	input := []byte(`{
	"data": [
		{
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
		{
			"id": "someOtherID",
			"type": "samples",
			"attributes": {
				"float64": 9999999.9999,
				"int": 999999,
				"string": "99SomeString"
			},
			"meta": {
				"float64": 12121.2321223,
				"int": 12312313,
				"string": "someString99"
			}
		}
	],
	"jsonapi": {
		"version": "1.0"
	}
}`)
	s := []Sample{}
	if err := Unmarshal(input, &s); err != nil {
		t.Errorf(err.Error())
	}

	if len(s) != 2 {
		t.Errorf("expected slice of len: %d, got: %d", 2, len(s))
	}
	if s[0].ID != "someID" {
		t.Errorf("ID was incorrect, got: %s, want: %s.", s[0].ID, "someID")
	}
	if s[0].Int != 99 {
		t.Errorf("Int was incorrect, got: %d, want: %d.", s[0].Int, 99)
	}
	if s[0].Float64 != 3.14159265359 {
		t.Errorf("Float64 was incorrect, got: %f, want: %f.", s[0].Float64, 3.14159265359)
	}
	if s[0].String != "someString" {
		t.Errorf("String was incorrect, got: %s, want: %s.", s[0].String, "someString")
	}
	if s[0].MetaInt != 5845 {
		t.Errorf("MetaInt was incorrect, got: %d, want: %d.", s[0].MetaInt, 5845)
	}
	if s[0].MetaString != "someString" {
		t.Errorf("MetaString was incorrect, got: %s, want: %s.", s[0].MetaString, "someString")
	}
	if s[0].MetaFloat64 != 99.2486135148 {
		t.Errorf("MetaFloat64 was incorrect, got: %f, want: %f.", s[0].MetaFloat64, 99.2486135148)
	}
	if s[1].ID != "someOtherID" {
		t.Errorf("ID was incorrect, got: %s, want: %s.", s[1].ID, "someOtherID")
	}
	if s[1].Int != 999999 {
		t.Errorf("Int was incorrect, got: %d, want: %d.", s[1].Int, 999999)
	}
	if s[1].Float64 != 9999999.9999 {
		t.Errorf("Float64 was incorrect, got: %f, want: %f.", s[1].Float64, 9999999.9999)
	}
	if s[1].String != "99SomeString" {
		t.Errorf("String was incorrect, got: %s, want: %s.", s[1].String, "99SomeString")
	}
	if s[1].MetaInt != 12312313 {
		t.Errorf("MetaInt was incorrect, got: %d, want: %d.", s[1].MetaInt, 12312313)
	}
	if s[1].MetaString != "someString99" {
		t.Errorf("MetaString was incorrect, got: %s, want: %s.", s[1].MetaString, "someString99")
	}
	if s[1].MetaFloat64 != 12121.2321223 {
		t.Errorf("MetaFloat64 was incorrect, got: %f, want: %f.", s[1].MetaFloat64, 12121.2321223)
	}
}
