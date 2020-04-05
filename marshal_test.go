package jsonapi

import (
	"bytes"
	"fmt"
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
	b, err := Marshal(&s, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	if bytes.Compare(input, b) != 0 {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", string(input), string(b))
	}

	// test incorrectly passing a non pointer to a struct
	notPointerOrSliceError := "v must be pointer or slice"
	if b, err := Marshal(s, nil); err == nil {
		fmt.Println(string(b))
		t.Errorf("marshal must error out if v is not a pointer or a slice")
	} else {
		if err.Error() != notPointerOrSliceError {
			t.Errorf("marshal must error out if v is not a pointer or a slice with error: %s, got: %s", notPointerOrSliceError, err.Error())
		}
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

func TestMarshalCustomTypeString(t *testing.T) {
	type CustomStr string
	type TestString struct {
		ID  string    `jsonapi:"primary,test_strings"`
		Foo CustomStr `jsonapi:"attribute,bar"`
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

func TestMarshalCustomTypeStringPtr(t *testing.T) {
	type CustomStr string
	type TestStringPtr struct {
		ID  string     `jsonapi:"primary,test_strings"`
		Foo *CustomStr `jsonapi:"attribute,bar"`
	}
	s := CustomStr("hello world!")
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

func TestMarshalCompound(t *testing.T) {
	type TestCompound struct {
		ID    string `jsonapi:"primary,test_compounds"`
		Foo   string `jsonapi:"attribute,bar"`
		Links Links  `jsonapi:"links,self"`
	}
	tcs := []*TestCompound{
		{
			ID:  "someID1",
			Foo: "hello",
			Links: Links{
				"self": "/self/link",
			},
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
			},
			"links": {
				"self": "/self/link"
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
	},
	"meta": {
		"hello": "world!"
	},
	"links": {
		"self": "/foo/bar"
	}
}`)
	if b, err := Marshal(&tcs, &MarshalParams{
		Meta: &Meta{
			"hello": "world!",
		},
		Links: &Links{
			"self": "/foo/bar",
		},
	}); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestMarshalCompoundWithRelationships(t *testing.T) {
	type Author struct {
		ID   string `jsonapi:"primary,authors"`
		Name string `jsonapi:"attribute,name"`
	}
	type Article struct {
		ID     string  `jsonapi:"primary,articles"`
		Title  string  `jsonapi:"attribute,title"`
		Author *Author `jsonapi:"relationship,author"`
	}
	articles := []*Article{
		{
			ID:    "article-1",
			Title: "Hello world 1!",
			Author: &Author{
				ID:   "author-1",
				Name: "John",
			},
		},
		{
			ID:    "article-2",
			Title: "Hello world 2!",
			Author: &Author{
				ID:   "author-2",
				Name: "Juan",
			},
		},
	}
	expected := []byte(`{
	"data": [
		{
			"id": "article-1",
			"type": "articles",
			"attributes": {
				"title": "Hello world 1!"
			},
			"relationships": {
				"author": {
					"data": {
						"id": "author-1",
						"type": "authors"
					}
				}
			}
		},
		{
			"id": "article-2",
			"type": "articles",
			"attributes": {
				"title": "Hello world 2!"
			},
			"relationships": {
				"author": {
					"data": {
						"id": "author-2",
						"type": "authors"
					}
				}
			}
		}
	],
	"jsonapi": {
		"version": "1.0"
	},
	"included": [
		{
			"id": "author-1",
			"type": "authors",
			"attributes": {
				"name": "John"
			}
		},
		{
			"id": "author-2",
			"type": "authors",
			"attributes": {
				"name": "Juan"
			}
		}
	]
}`)
	if b, err := Marshal(&articles, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestMarshalCompoundWithCompoundRelationships(t *testing.T) {
	type Author struct {
		ID   string `jsonapi:"primary,authors"`
		Name string `jsonapi:"attribute,name"`
	}
	type Article struct {
		ID      string    `jsonapi:"primary,articles"`
		Title   string    `jsonapi:"attribute,title"`
		Authors []*Author `jsonapi:"relationship,authors"`
	}
	articles := []*Article{
		{
			ID:    "article-1",
			Title: "Hello world 1!",
			Authors: []*Author{
				&Author{
					ID:   "author-1",
					Name: "John",
				},
				&Author{
					ID:   "author-3",
					Name: "Fred",
				},
			},
		},
		{
			ID:    "article-2",
			Title: "Hello world 2!",
			Authors: []*Author{
				&Author{
					ID:   "author-2",
					Name: "Juan",
				},
				&Author{
					ID:   "author-3",
					Name: "Fred",
				},
			},
		},
	}
	expected := []byte(`{
	"data": [
		{
			"id": "article-1",
			"type": "articles",
			"attributes": {
				"title": "Hello world 1!"
			},
			"relationships": {
				"authors": {
					"data": [
						{
							"id": "author-1",
							"type": "authors"
						},
						{
							"id": "author-3",
							"type": "authors"
						}
					]
				}
			}
		},
		{
			"id": "article-2",
			"type": "articles",
			"attributes": {
				"title": "Hello world 2!"
			},
			"relationships": {
				"authors": {
					"data": [
						{
							"id": "author-2",
							"type": "authors"
						},
						{
							"id": "author-3",
							"type": "authors"
						}
					]
				}
			}
		}
	],
	"jsonapi": {
		"version": "1.0"
	},
	"included": [
		{
			"id": "author-1",
			"type": "authors",
			"attributes": {
				"name": "John"
			}
		},
		{
			"id": "author-3",
			"type": "authors",
			"attributes": {
				"name": "Fred"
			}
		},
		{
			"id": "author-2",
			"type": "authors",
			"attributes": {
				"name": "Juan"
			}
		}
	]
}`)
	if b, err := Marshal(&articles, nil); err != nil {
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

func TestMarshalRelationship(t *testing.T) {
	type Bar struct {
		ID    string `jsonapi:"primary,bars"`
		Hello string `jsonapi:"attribute,hello"`
	}
	type TestRelationship struct {
		ID      string `jsonapi:"primary,test_relationships"`
		Foo     string `jsonapi:"attribute,foo"`
		Bar     *Bar   `jsonapi:"relationship,bar"`
		Another *Bar   `jsonapi:"relationship,another"`
		Repeat  *Bar   `jsonapi:"relationship,repeat"`
	}
	test := TestRelationship{
		ID:  "someID",
		Foo: "bar",
		Bar: &Bar{
			ID:    "barID",
			Hello: "world!",
		},
		Another: &Bar{
			ID:    "barID2",
			Hello: "world2!",
		},
		Repeat: &Bar{
			ID:    "barID2",
			Hello: "world2!",
		},
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_relationships",
		"attributes": {
			"foo": "bar"
		},
		"relationships": {
			"another": {
				"data": {
					"id": "barID2",
					"type": "bars"
				}
			},
			"bar": {
				"data": {
					"id": "barID",
					"type": "bars"
				}
			},
			"repeat": {
				"data": {
					"id": "barID2",
					"type": "bars"
				}
			}
		}
	},
	"jsonapi": {
		"version": "1.0"
	},
	"included": [
		{
			"id": "barID",
			"type": "bars",
			"attributes": {
				"hello": "world!"
			}
		},
		{
			"id": "barID2",
			"type": "bars",
			"attributes": {
				"hello": "world2!"
			}
		}
	]
}`)
	if got, err := Marshal(&test, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expected) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(got))
		}
	}
}

func TestMarshalRelationshipEmpty(t *testing.T) {
	type Bar struct {
		ID    string `jsonapi:"primary,bars"`
		Hello string `jsonapi:"attribute,hello"`
	}
	type TestRelationship struct {
		ID      string `jsonapi:"primary,test_relationships"`
		Foo     string `jsonapi:"attribute,foo"`
		Bar     *Bar   `jsonapi:"relationship,bar"`
		Another *Bar   `jsonapi:"relationship,another"`
	}
	test := TestRelationship{
		ID:  "someID",
		Foo: "bar",
		Bar: &Bar{
			ID:    "barID",
			Hello: "world!",
		},
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_relationships",
		"attributes": {
			"foo": "bar"
		},
		"relationships": {
			"another": {
				"data": null
			},
			"bar": {
				"data": {
					"id": "barID",
					"type": "bars"
				}
			}
		}
	},
	"jsonapi": {
		"version": "1.0"
	},
	"included": [
		{
			"id": "barID",
			"type": "bars",
			"attributes": {
				"hello": "world!"
			}
		}
	]
}`)
	if got, err := Marshal(&test, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expected) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(got))
		}
	}
}

func TestMarshalRelationshipArray(t *testing.T) {
	type Bar struct {
		ID    string `jsonapi:"primary,bars"`
		Hello string `jsonapi:"attribute,hello"`
	}
	type TestRelationship struct {
		ID      string `jsonapi:"primary,test_relationships"`
		Foo     string `jsonapi:"attribute,foo"`
		Bars    []*Bar `jsonapi:"relationship,bars"`
		Another *Bar   `jsonapi:"relationship,another"`
	}
	test := TestRelationship{
		ID:  "someID",
		Foo: "bar",
		Bars: []*Bar{
			&Bar{
				ID:    "barID1.1",
				Hello: "world1.1!",
			},
			&Bar{
				ID:    "barID1.2",
				Hello: "world1.2!",
			},
		},
		Another: &Bar{
			ID:    "barID2",
			Hello: "world2!",
		},
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_relationships",
		"attributes": {
			"foo": "bar"
		},
		"relationships": {
			"another": {
				"data": {
					"id": "barID2",
					"type": "bars"
				}
			},
			"bars": {
				"data": [
					{
						"id": "barID1.1",
						"type": "bars"
					},
					{
						"id": "barID1.2",
						"type": "bars"
					}
				]
			}
		}
	},
	"jsonapi": {
		"version": "1.0"
	},
	"included": [
		{
			"id": "barID1.1",
			"type": "bars",
			"attributes": {
				"hello": "world1.1!"
			}
		},
		{
			"id": "barID1.2",
			"type": "bars",
			"attributes": {
				"hello": "world1.2!"
			}
		},
		{
			"id": "barID2",
			"type": "bars",
			"attributes": {
				"hello": "world2!"
			}
		}
	]
}`)
	if got, err := Marshal(&test, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expected) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(got))
		}
	}
}

func TestMarshalRelationshipEmptyArray(t *testing.T) {
	type Bar struct {
		ID    string `jsonapi:"primary,bars"`
		Hello string `jsonapi:"attribute,hello"`
	}
	type TestRelationship struct {
		ID      string `jsonapi:"primary,test_relationships"`
		Foo     string `jsonapi:"attribute,foo"`
		Bars    []*Bar `jsonapi:"relationship,bars"`
		Another *Bar   `jsonapi:"relationship,another"`
	}
	test := TestRelationship{
		ID:  "someID",
		Foo: "bar",
		Another: &Bar{
			ID:    "barID2",
			Hello: "world2!",
		},
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_relationships",
		"attributes": {
			"foo": "bar"
		},
		"relationships": {
			"another": {
				"data": {
					"id": "barID2",
					"type": "bars"
				}
			},
			"bars": {
				"data": []
			}
		}
	},
	"jsonapi": {
		"version": "1.0"
	},
	"included": [
		{
			"id": "barID2",
			"type": "bars",
			"attributes": {
				"hello": "world2!"
			}
		}
	]
}`)
	if got, err := Marshal(&test, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, expected) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(got))
		}
	}
}

func TestMarshalErrors2(t *testing.T) {
	// unmarshal non-pointer
	nonPointerOrSliceErrMsg := "v must be pointer or slice"
	nonPointerOrSlice := Sample{}
	_, nonPointerOrSliceErr := Marshal(nonPointerOrSlice, nil)
	switch {
	case nonPointerOrSliceErr == nil:
		t.Errorf("expected error: %s, but got no error", nonPointerOrSliceErrMsg)
	case nonPointerOrSliceErr.Error() != nonPointerOrSliceErrMsg:
		t.Errorf("expected error: %s, got: %s", nonPointerOrSliceErrMsg, nonPointerOrSliceErr.Error())
	}

	// missing type
	type MissingType struct {
		ID  string `jsonapi:"primary,"`
		Foo string `jsonapi:"attribute,foo"`
	}
	missingType := &MissingType{
		ID:  "missing-type-1",
		Foo: "bar",
	}
	missmissingTypeErrMsg := "type must be set"
	_, missmissingTypeErr := Marshal(missingType, nil)
	switch {
	case missmissingTypeErr == nil:
		t.Errorf("expected error: %s, but got no error", missmissingTypeErrMsg)
	case missmissingTypeErr.Error() != missmissingTypeErrMsg:
		t.Errorf("expected error: %s, got: %s", missmissingTypeErrMsg, missmissingTypeErr.Error())
	}

	// wrong id type
	type WrongIDType struct {
		ID  bool   `jsonapi:"primary,wrong_id_types"`
		Foo string `jsonapi:"attribute,foo"`
	}
	wrongIDType := &WrongIDType{
		ID:  true,
		Foo: "bar",
	}
	wrongIDTypeErrMsg := "ID must be a string, got bool"
	_, wrongIDTypeErr := Marshal(wrongIDType, nil)
	switch {
	case wrongIDTypeErr == nil:
		t.Errorf("expected error: %s, but got no error", wrongIDTypeErrMsg)
	case wrongIDTypeErr.Error() != wrongIDTypeErrMsg:
		t.Errorf("expected error: %s, got: %s", wrongIDTypeErrMsg, wrongIDTypeErr.Error())
	}

	// wrong id type in relationship
	type WrongIDTypeInRel struct {
		ID          string       `jsonapi:"primary,wrong_id_type_in_rels"`
		Foo         string       `jsonapi:"attribute,foo"`
		WrongIDType *WrongIDType `jsonapi:"relationship,wrong_id_type"`
	}
	wrongIDTypeInRel := &WrongIDTypeInRel{
		ID:  "wrong-id-type-in-rel-1",
		Foo: "bar",
		WrongIDType: &WrongIDType{
			ID:  false,
			Foo: "bar",
		},
	}
	wrongIDTypeInRelErrMsg := "ID must be a string, got bool"
	_, wrongIDTypeInRelErr := Marshal(wrongIDTypeInRel, nil)
	switch {
	case wrongIDTypeInRelErr == nil:
		t.Errorf("expected error: %s, but got no error", wrongIDTypeInRelErrMsg)
	case wrongIDTypeInRelErr.Error() != wrongIDTypeInRelErrMsg:
		t.Errorf("expected error: %s, got: %s", wrongIDTypeInRelErrMsg, wrongIDTypeInRelErr.Error())
	}

	// wrong id type in compound relationship
	type WrongIDTypeInRels struct {
		ID           string         `jsonapi:"primary,wrong_id_type_in_rels"`
		Foo          string         `jsonapi:"attribute,foo"`
		WrongIDType  *WrongIDType   `jsonapi:"relationship,wrong_id_type"`
		WrongIDTypes []*WrongIDType `jsonapi:"relationship,wrong_id_types"`
	}
	wrongIDTypeInRels := &WrongIDTypeInRels{
		ID:  "wrong-id-type-in-rel-1",
		Foo: "bar",
		WrongIDTypes: []*WrongIDType{
			&WrongIDType{
				ID:  false,
				Foo: "bar",
			},
		},
	}
	wrongIDTypeInRelsErrMsg := "ID must be a string, got bool"
	_, wrongIDTypeInRelsErr := Marshal(wrongIDTypeInRels, nil)
	switch {
	case wrongIDTypeInRelsErr == nil:
		t.Errorf("expected error: %s, but got no error", wrongIDTypeInRelsErrMsg)
	case wrongIDTypeInRelsErr.Error() != wrongIDTypeInRelsErrMsg:
		t.Errorf("expected error: %s, got: %s", wrongIDTypeInRelsErrMsg, wrongIDTypeInRelsErr.Error())
	}

	// wrong id type in relationships in compound document
	wrongIDTypesInRels := &[]*WrongIDTypeInRels{
		&WrongIDTypeInRels{
			ID:  "wrong-id-type-in-rel-1",
			Foo: "bar",
			WrongIDType: &WrongIDType{
				ID:  false,
				Foo: "test",
			},
			WrongIDTypes: []*WrongIDType{
				&WrongIDType{
					ID:  true,
					Foo: "bar",
				},
			},
		},
		&WrongIDTypeInRels{
			ID:  "wrong-id-type-in-rel-2",
			Foo: "bar",
			WrongIDType: &WrongIDType{
				ID:  false,
				Foo: "test",
			},
			WrongIDTypes: []*WrongIDType{
				&WrongIDType{
					ID:  false,
					Foo: "bar",
				},
			},
		},
	}
	wrongIDTypesInRelsErrMsg := "ID must be a string, got bool"
	_, wrongIDTypesInRelsErr := Marshal(wrongIDTypesInRels, nil)
	switch {
	case wrongIDTypesInRelsErr == nil:
		t.Errorf("expected error: %s, but got no error", wrongIDTypesInRelsErrMsg)
	case wrongIDTypesInRelsErr.Error() != wrongIDTypesInRelsErrMsg:
		t.Errorf("expected error: %s, got: %s", wrongIDTypesInRelsErrMsg, wrongIDTypesInRelsErr.Error())
	}
}
