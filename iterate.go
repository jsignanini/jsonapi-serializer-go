package jsonapi

import (
	"fmt"
	"reflect"
)

type iterFunc func(reflect.Value, memberType, ...string) error

func iterateStruct(v interface{}, iter iterFunc, memberNames ...string) error {
	rType := reflect.TypeOf(v)
	rValue := reflect.ValueOf(v)

	// check v is a pointer
	if rType.Kind() != reflect.Ptr {
		return fmt.Errorf("v must be a pointer")
	}

	// skip nil pointers
	if rValue.IsNil() {
		return nil
	}

	// check *v is a struct
	if rValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}

	// iterate struct fields
	numFields := rValue.Elem().NumField()
	for i := 0; i < numFields; i++ {
		fType := rType.Elem().Field(i)
		fValue := rValue.Elem().Field(i)
		kind := fValue.Kind()

		// if struct and embedded (anonymus), restart loop
		if kind == reflect.Struct && fType.Anonymous {
			iterateStruct(fValue.Addr().Interface(), iter, memberNames...)
			continue
		}

		// if tag exists, get member info, continue otherwise
		if _, ok := fType.Tag.Lookup(tagKey); !ok {
			continue
		}
		memberType, memberName, err := getMember(fType)
		if err != nil {
			return err
		}

		// handle nested structs
		if kind == reflect.Struct {
			iterateStruct(fValue.Addr().Interface(), iter, append(memberNames, memberName)...)
			continue
		}

		if err := iter(fValue, memberType, append(memberNames, memberName)...); err != nil {
			return err
		}
	}
	return nil
}
