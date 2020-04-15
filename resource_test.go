package jsonapi

import (
	"reflect"
	"testing"
)

func TestSetIDAndType(t *testing.T) {
	stringID := reflect.ValueOf("someID")
	intID := reflect.ValueOf(3234324)
	invalidID := reflect.ValueOf(false)

	r1 := NewResource()
	if err := r1.SetIDAndType(stringID, "articles"); err != nil {
		t.Error(err.Error())
	}

	r2 := NewResource()
	if err := r2.SetIDAndType(intID, "articles"); err != nil {
		t.Error(err.Error())
	}

	r3 := NewResource()
	if err := r3.SetIDAndType(stringID, ""); err == nil {
		t.Error("resource type must be set")
	}

	r4 := NewResource()
	if err := r4.SetIDAndType(invalidID, "articles"); err == nil {
		t.Error("resource id should only accept string or int type")
	}
}
