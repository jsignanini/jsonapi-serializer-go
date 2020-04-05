package jsonapi

import (
	"fmt"
	"math"
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

	// TODO ignored field
	// IgnoredField string `jsonapi:"-"`
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

	// test incorrectly sending a string instead of an int
	wrongTypeOut := TestBool{}
	wrongTypeErr := "invalid value for field bool"
	wrongType := []byte(`{
		"data": {
			"id": "sample-1",
			"type": "floats",
			"attributes": {
				"is_true": "wrong string"
			}
		}
	}`)
	if err := Unmarshal(wrongType, &wrongTypeOut); err == nil {
		t.Errorf("expected error: %s, got no error", wrongTypeErr)
	} else {
		if err.Error() != wrongTypeErr {
			t.Errorf("expected error: %s, got: %s", wrongTypeErr, err.Error())
		}
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

func TestUnmarshalInts(t *testing.T) {
	type Sample struct {
		ID    string `jsonapi:"primary,ints"`
		Int   int    `jsonapi:"attribute,int"`
		Int8  int8   `jsonapi:"attribute,int8"`
		Int16 int16  `jsonapi:"attribute,int16"`
		Int32 int32  `jsonapi:"attribute,int32"`
		Int64 int64  `jsonapi:"attribute,int64"`
	}
	maxOut := Sample{}
	max := []byte(fmt.Sprintf(`{
	"data": {
		"id": "sample-1",
		"type": "ints",
		"attributes": {
			"int": 10,
			"int8": %d,
			"int16": %d,
			"int32": %d,
			"int64": %d
		}
	}
}`, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64))
	if err := Unmarshal(max, &maxOut); err != nil {
		t.Errorf(err.Error())
	}
	if maxOut.ID != "sample-1" {
		t.Errorf("expected id to be: %s, got: %s", "sample-1", maxOut.ID)
	}
	if maxOut.Int != 10 {
		t.Errorf("expected int to be: %d, got: %d", 10, maxOut.Int)
	}
	if maxOut.Int8 != math.MaxInt8 {
		t.Errorf("expected int8 to be: %d, got: %d", math.MaxInt8, maxOut.Int8)
	}
	if maxOut.Int16 != math.MaxInt16 {
		t.Errorf("expected int16 to be: %d, got: %d", math.MaxInt16, maxOut.Int16)
	}
	if maxOut.Int32 != math.MaxInt32 {
		t.Errorf("expected int32 to be: %d, got: %d", math.MaxInt32, maxOut.Int32)
	}
	if maxOut.Int64 != math.MaxInt64 {
		t.Errorf("expected int64 to be: %d, got: %d", math.MaxInt64, maxOut.Int64)
	}
}

func TestUnmarshalUints(t *testing.T) {
	type Sample struct {
		ID     string `jsonapi:"primary,uints"`
		Uint   uint   `jsonapi:"attribute,uint"`
		Uint8  uint8  `jsonapi:"attribute,uint8"`
		Uint16 uint16 `jsonapi:"attribute,uint16"`
		Uint32 uint32 `jsonapi:"attribute,uint32"`
		Uint64 uint64 `jsonapi:"attribute,uint64"`
	}

	// test maximums
	maxOut := Sample{}
	max := []byte(fmt.Sprintf(`{
	"data": {
		"id": "sample-1",
		"type": "uints",
		"attributes": {
			"uint": 10,
			"uint8": %d,
			"uint16": %d,
			"uint32": %d,
			"uint64": %d
		}
	}
}`, math.MaxUint8, math.MaxUint16, math.MaxUint32, uint64(math.MaxUint64)))
	if err := Unmarshal(max, &maxOut); err != nil {
		t.Errorf(err.Error())
	}
	if maxOut.ID != "sample-1" {
		t.Errorf("expected id to be: %s, got: %s", "sample-1", maxOut.ID)
	}
	if maxOut.Uint != 10 {
		t.Errorf("expected uint to be: %d, got: %d", 10, maxOut.Uint)
	}
	if maxOut.Uint8 != math.MaxUint8 {
		t.Errorf("expected uint8 to be: %d, got: %d", math.MaxUint8, maxOut.Uint8)
	}
	if maxOut.Uint16 != math.MaxUint16 {
		t.Errorf("expected uint16 to be: %d, got: %d", math.MaxUint16, maxOut.Uint16)
	}
	if maxOut.Uint32 != math.MaxUint32 {
		t.Errorf("expected uint32 to be: %d, got: %d", math.MaxUint32, maxOut.Uint32)
	}
	if maxOut.Uint64 != math.MaxUint64 {
		t.Errorf("expected uint64 to be: %d, got: %d", uint64(math.MaxUint64), maxOut.Uint64)
	}

	// test incorrectly sending a string instead of an int
	wrongTypeOut := Sample{}
	wrongTypeErr := "number has no digits"
	wrongType := []byte(`{
	"data": {
		"id": "sample-1",
		"type": "uints",
		"attributes": {
			"uint": "wrong string"
		}
	}
}`)
	if err := Unmarshal(wrongType, &wrongTypeOut); err == nil {
		t.Errorf("expected error: %s, got no error", wrongTypeErr)
	} else {
		if err.Error() != wrongTypeErr {
			t.Errorf("expected error: %s, got: %s", wrongTypeErr, err.Error())
		}
	}
}

func TestUnmarshalFloats(t *testing.T) {
	type Sample struct {
		ID      string  `jsonapi:"primary,uints"`
		Float32 float32 `jsonapi:"attribute,float32"`
		Float64 float64 `jsonapi:"attribute,float64"`
	}

	// test maximums
	maxOut := Sample{}
	max := []byte(fmt.Sprintf(`{
	"data": {
		"id": "sample-1",
		"type": "floats",
		"attributes": {
			"float32": %f,
			"float64": %f
		}
	}
}`, math.MaxFloat32, math.MaxFloat64))
	if err := Unmarshal(max, &maxOut); err != nil {
		t.Errorf(err.Error())
	}
	if maxOut.ID != "sample-1" {
		t.Errorf("expected id to be: %s, got: %s", "sample-1", maxOut.ID)
	}
	if maxOut.Float32 != math.MaxFloat32 {
		t.Errorf("expected float32 to be: %f, got: %f", math.MaxFloat32, maxOut.Float32)
	}
	if maxOut.Float64 != math.MaxFloat64 {
		t.Errorf("expected float64 to be: %f, got: %f", math.MaxFloat64, maxOut.Float64)
	}

	// test incorrectly sending a string instead of a float
	wrongTypeOut := Sample{}
	wrongTypeErr := "number has no digits"
	wrongType := []byte(`{
	"data": {
		"id": "sample-1",
		"type": "floats",
		"attributes": {
			"float32": "wrong string"
		}
	}
}`)
	if err := Unmarshal(wrongType, &wrongTypeOut); err == nil {
		t.Errorf("expected error: %s, got no error", wrongTypeErr)
	} else {
		if err.Error() != wrongTypeErr {
			t.Errorf("expected error: %s, got: %s", wrongTypeErr, err.Error())
		}
	}
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

func TestUnmarshalString(t *testing.T) {
	type Sample struct {
		ID     string `jsonapi:"primary,strings"`
		String string `jsonapi:"attribute,string"`
	}
	wrongTypeErrMsg := "invalid value for field string"
	wrongTypeOut := Sample{}
	wrongType := []byte(`{
	"data": {
		"id": "string-id",
		"type": "strings",
		"attributes": {
			"string": {
				"foo": "bar"
			}
		}
	}
}`)
	wrongTypeErr := Unmarshal(wrongType, &wrongTypeOut)
	switch {
	case wrongTypeErr == nil:
		t.Errorf("expected error: %s, but got no error", wrongTypeErrMsg)
	case wrongTypeErr.Error() != wrongTypeErrMsg:
		t.Errorf("expected error: %s, got: %s", wrongTypeErrMsg, wrongTypeErr.Error())
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

func TestUnmarshalErrors(t *testing.T) {
	// unmarshal non-pointer
	nonPointerErrMsg := "v must be pointer"
	nonPointer := Sample{}
	nonPointerErr := Unmarshal([]byte(""), nonPointer)
	switch {
	case nonPointerErr == nil:
		t.Errorf("expected error: %s, but got no error", nonPointerErrMsg)
	case nonPointerErr.Error() != nonPointerErrMsg:
		t.Errorf("expected error: %s, got: %s", nonPointerErrMsg, nonPointerErr.Error())
	}

	// unmarsshal nil pointer
	nilPointerErrMsg := "v must not be nil"
	var nilPointer *Sample
	nilPointerErr := Unmarshal([]byte(""), nilPointer)
	switch {
	case nilPointerErr == nil:
		t.Errorf("expected error: %s, but got no error", nilPointerErrMsg)
	case nilPointerErr.Error() != nilPointerErrMsg:
		t.Errorf("expected error: %s, got: %s", nilPointerErrMsg, nilPointerErr.Error())
	}

	// malformed document json
	malformedJSON := Sample{}
	malformedJSONErr := Unmarshal([]byte("malformed"), &malformedJSON)
	switch {
	case malformedJSONErr == nil:
		t.Error("expected malformed JSON to error out but got no error")
	}

	// malformed compound document json
	malformedCompoundJSON := []*Sample{}
	malformedCompoundJSONErr := Unmarshal([]byte("malformed compound"), &malformedCompoundJSON)
	switch {
	case malformedCompoundJSONErr == nil:
		t.Error("expected malformed JSON to error out but got no error")
	}
}
