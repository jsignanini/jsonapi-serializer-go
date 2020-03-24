package jsonapi

import (
	"reflect"
	"testing"
)

func TestIterateStruct(t *testing.T) {
	shimIterFunc := iterFunc(func(value reflect.Value, memberType MemberType, memberNames ...string) error { return nil })
	shimMemberNames := []string{}

	// test incorrectly passing a non pointer
	notAPointer := 10
	notAPointerError := "v must be a pointer"
	if err := iterateStruct(notAPointer, shimIterFunc, shimMemberNames...); err == nil {
		t.Errorf("iterateStruct must error out if v is not a pointer")
	} else {
		if err.Error() != notAPointerError {
			t.Errorf("passing a non pointer to iterateStruct should error out with message: %s, but got: %s", notAPointerError, err.Error())
		}
	}

	// test passing a nil pointer
	var nilPointer *int
	if err := iterateStruct(nilPointer, shimIterFunc, shimMemberNames...); err != nil {
		t.Errorf("iterateStruct must not error out if passed a nil pointer, got error: %s", err.Error())
	}

	// test incorrectly passing a pointer to a non struct
	notAPointerToAStruct := 10
	notAPointerToAStructError := "v must be a pointer to a struct"
	if err := iterateStruct(&notAPointerToAStruct, shimIterFunc, shimMemberNames...); err != nil {
		if err.Error() != notAPointerToAStructError {
			t.Errorf("passing a pointer to a non struct to iterateStruct should error out with message: %s, but got: %s", notAPointerToAStructError, err.Error())
		}
	} else {
		t.Errorf("iterateStruct must error out if not passed a pointer to a struct, got no error")
	}
}
