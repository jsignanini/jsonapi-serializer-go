package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	document := Document{}
	if err := json.Unmarshal(data, &document); err != nil {
		return err
	}

	rType := reflect.TypeOf(v)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	for i := 0; i < rType.NumField(); i++ {
		memberType, memberName, err := getMember(rType.Field(i))
		if err != nil {
			return err
		}
		resourceValue := reflect.ValueOf(v).Elem().Field(i)
		resourceKind := resourceValue.Kind()

		if memberType == MemberTypePrimary {
			if resourceKind != reflect.String {
				return fmt.Errorf("ID must be a string")
			}
			resourceValue.SetString(document.Data.ID)
			continue
		}

		var search map[string]interface{}
		switch memberType {
		case MemberTypeAttribute:
			search = document.Data.Attributes
		case MemberTypeMeta:
			search = document.Data.Meta
		}

		// skip if member missing from JSON
		if _, ok := search[memberName]; !ok {
			continue
		}

		if err := unmarshal(search[memberName], &resourceValue); err != nil {
			return err
		}
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
	tag := field.Tag.Get(tagKey)
	// TODO error check tagParts length?
	tagParts := strings.Split(tag, ",")
	memberType, err := NewMemberType(tagParts[0])
	if err != nil {
		return "", "", err
	}
	return memberType, tagParts[1], nil
}
