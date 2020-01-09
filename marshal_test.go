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
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalInt8(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo int8   `jsonapi:"attribute,bar"`
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

func TestMarshalInt8Ptr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *int8  `jsonapi:"attribute,bar"`
	}
	seven := int8(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalInt16(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo int16  `jsonapi:"attribute,bar"`
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

func TestMarshalInt16Ptr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *int16 `jsonapi:"attribute,bar"`
	}
	seven := int16(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalInt32(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo int32  `jsonapi:"attribute,bar"`
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

func TestMarshalInt32Ptr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *int32 `jsonapi:"attribute,bar"`
	}
	seven := int32(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalInt64(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo int64  `jsonapi:"attribute,bar"`
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

func TestMarshalInt64Ptr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *int64 `jsonapi:"attribute,bar"`
	}
	seven := int64(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalUint(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo uint   `jsonapi:"attribute,bar"`
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

func TestMarshalUintPtr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *uint  `jsonapi:"attribute,bar"`
	}
	seven := uint(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalUint8(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo uint8  `jsonapi:"attribute,bar"`
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

func TestMarshalUint8Ptr(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo *uint8 `jsonapi:"attribute,bar"`
	}
	seven := uint8(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalUint16(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo uint16 `jsonapi:"attribute,bar"`
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

func TestMarshalUint16Ptr(t *testing.T) {
	type TestInt struct {
		ID  string  `jsonapi:"primary,test_ints"`
		Foo *uint16 `jsonapi:"attribute,bar"`
	}
	seven := uint16(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalUint32(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo uint32 `jsonapi:"attribute,bar"`
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

func TestMarshalUint32Ptr(t *testing.T) {
	type TestInt struct {
		ID  string  `jsonapi:"primary,test_ints"`
		Foo *uint32 `jsonapi:"attribute,bar"`
	}
	seven := uint32(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalUint64(t *testing.T) {
	type TestInt struct {
		ID  string `jsonapi:"primary,test_ints"`
		Foo uint64 `jsonapi:"attribute,bar"`
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

func TestMarshalUint64Ptr(t *testing.T) {
	type TestInt struct {
		ID  string  `jsonapi:"primary,test_ints"`
		Foo *uint64 `jsonapi:"attribute,bar"`
	}
	seven := uint64(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalUintptr(t *testing.T) {
	type TestInt struct {
		ID  string  `jsonapi:"primary,test_ints"`
		Foo uintptr `jsonapi:"attribute,bar"`
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

func TestMarshalUintptrPtr(t *testing.T) {
	type TestInt struct {
		ID  string   `jsonapi:"primary,test_ints"`
		Foo *uintptr `jsonapi:"attribute,bar"`
	}
	seven := uintptr(7)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
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
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalFloat32(t *testing.T) {
	type TestInt struct {
		ID  string  `jsonapi:"primary,test_ints"`
		Foo float32 `jsonapi:"attribute,bar"`
	}
	seven := TestInt{
		ID:  "someID",
		Foo: 7.99,
	}
	want := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints",
		"attributes": {
			"bar": 7.99
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

func TestMarshalFloat32Ptr(t *testing.T) {
	type TestInt struct {
		ID  string   `jsonapi:"primary,test_ints"`
		Foo *float32 `jsonapi:"attribute,bar"`
	}
	seven := float32(7.99)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints",
		"attributes": {
			"bar": 7.99
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtr, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}

func TestMarshalFloat64(t *testing.T) {
	type TestInt struct {
		ID  string  `jsonapi:"primary,test_ints"`
		Foo float64 `jsonapi:"attribute,bar"`
	}
	seven := TestInt{
		ID:  "someID",
		Foo: 7.99,
	}
	want := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints",
		"attributes": {
			"bar": 7.99
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

func TestMarshalFloat64Ptr(t *testing.T) {
	type TestInt struct {
		ID  string   `jsonapi:"primary,test_ints"`
		Foo *float64 `jsonapi:"attribute,bar"`
	}
	seven := float64(7.99)
	sevenPtr := TestInt{
		ID:  "someID",
		Foo: &seven,
	}
	expectedSeven := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints",
		"attributes": {
			"bar": 7.99
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtr, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedSeven) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedSeven))
		}
	}
	sevenPtrNil := TestInt{
		ID: "someID",
	}
	expectedMissing := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_ints"
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&sevenPtrNil, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expectedMissing) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(got), string(expectedMissing))
		}
	}
}
