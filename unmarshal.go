package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type iterFunc func(reflect.StructField, reflect.Value, ...string) error

func iterateStruct(document *Document, iface interface{}, iter iterFunc, memberNames ...string) error {
	// TODO check iface is a struct
	fields := reflect.TypeOf(iface)
	values := reflect.ValueOf(iface)

	var numField int
	if fields.Kind() == reflect.Ptr {
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
			iterateStruct(document, value.Addr().Interface(), iter, memberNames...)
			continue
		}

		// get member info, continue otherwise
		_, memberName, err := getMember(field)
		if err != nil {
			continue
		}

		// handle nested structs
		if kind == reflect.Struct {
			iterateStruct(document, value.Addr().Interface(), iter, append(memberNames, memberName)...)
			continue
		}

		if err := iter(field, value, append(memberNames, memberName)...); err != nil {
			return err
		}
	}
	return nil
}

func getValueForMember(document *Document, memberType MemberType, memberNames ...string) (interface{}, error) {
	var search map[string]interface{}
	switch memberType {
	case MemberTypeAttribute:
		search = document.Data.Attributes
	case MemberTypeMeta:
		search = document.Data.Meta
	}
	for i, name := range memberNames {
		value, ok := search[name]
		if !ok {
			return "", fmt.Errorf("not ok")
		}
		if i == len(memberNames)-1 {
			return value, nil
		}
		search = search[name].(map[string]interface{})
	}
	return "", fmt.Errorf("memberNames was empty")
}

func Unmarshal(data []byte, v interface{}) error {
	document := Document{}
	if err := json.Unmarshal(data, &document); err != nil {
		return err
	}

	if err := iterateStruct(&document, v, func(field reflect.StructField, value reflect.Value, memberNames ...string) error {
		fieldKind := value.Kind()

		// get member info, continue otherwise
		memberType, _, err := getMember(field)
		if err != nil {
			return err
		}

		// TODO this sets ID for all nexted primary tag fields
		if memberType == MemberTypePrimary {
			if fieldKind != reflect.String {
				return fmt.Errorf("ID must be a string")
			}
			value.SetString(document.Data.ID)
			return nil
		}

		// get raw value
		v, err := getValueForMember(&document, memberType, memberNames...)
		if err != nil {
			return nil
		}

		// set raw value
		if err := unmarshal(v, &value); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func RegisterUnmarshaler(t reflect.Type, u unmarshalerFunc) {
	customUnmarshalers[t] = u
}

type unmarshalerFunc = func(interface{}, *reflect.Value)

var customUnmarshalers = make(map[reflect.Type]unmarshalerFunc)

func unmarshal(v interface{}, rv *reflect.Value) error {
	switch rv.Kind() {
	case reflect.String:
		if val, ok := v.(string); ok {
			rv.SetString(val)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, ok := v.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			rv.SetInt(int64(val))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if val, ok := v.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			rv.SetUint(uint64(val))
		}
	case reflect.Float32, reflect.Float64:
		if val, ok := v.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			rv.SetFloat(val)
		}
	default:
		cu, ok := customUnmarshalers[rv.Type()]
		if !ok {
			return fmt.Errorf("Type not supported, must implement custom unmarshaller")
		}
		cu(v, rv)
	}
	return nil
}

func getMember(field reflect.StructField) (MemberType, string, error) {
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		return "", "", fmt.Errorf("tag: %s, not specified", tagKey)
	}
	if tag == "" {
		return "", "", fmt.Errorf("tag: %s, was empty", tagKey)
	}
	tagParts := strings.Split(tag, ",")
	if len(tagParts) != 2 {
		return "", "", fmt.Errorf("tag: %s, was not formatted properly", tagKey)
	}
	memberType, err := NewMemberType(tagParts[0])
	if err != nil {
		return "", "", err
	}
	return memberType, tagParts[1], nil
}
