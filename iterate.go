package jsonapi

import (
	"reflect"
)

type iterFunc func(reflect.Value, MemberType, ...string) error

func iterateStruct(iface interface{}, iter iterFunc, memberNames ...string) error {
	// TODO check iface is a struct
	fields := reflect.TypeOf(iface)
	values := reflect.ValueOf(iface)

	var numField int
	if fields.Kind() == reflect.Ptr && values.IsNil() {
		return nil
	} else if fields.Kind() == reflect.Ptr {
		numField = values.Elem().NumField()
	} else {
		numField = values.NumField()
	}

	for i := 0; i < numField; i++ {
		field := fields.Elem().Field(i)
		value := values.Elem().Field(i)
		kind := value.Kind()

		// if struct and embedded (anonymus), restart loop
		if kind == reflect.Struct && field.Anonymous {
			iterateStruct(value.Addr().Interface(), iter, memberNames...)
			continue
		}

		// get member info, continue otherwise
		memberType, memberName, err := getMember(field)
		if err != nil {
			continue
		}

		// handle nested structs
		if kind == reflect.Struct {
			iterateStruct(value.Addr().Interface(), iter, append(memberNames, memberName)...)
			continue
		}

		if err := iter(value, memberType, append(memberNames, memberName)...); err != nil {
			return err
		}
	}
	return nil
}
