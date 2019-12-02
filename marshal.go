package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	data := Resource{
		Attributes: Attributes{},
		Meta:       Meta{},
	}

	rType := reflect.TypeOf(v)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	for i := 0; i < rType.NumField(); i++ {
		// get member info, continue otherwise
		memberType, memberName, err := getMember(rType.Field(i))
		if err != nil {
			continue
		}

		resourceValue := reflect.ValueOf(v).Elem().Field(i)
		resourceKind := resourceValue.Kind()

		if memberType == MemberTypePrimary {
			if resourceKind != reflect.String {
				return nil, fmt.Errorf("ID must be a string")
			}
			data.ID = resourceValue.Interface().(string)
			data.Type = memberName
			continue
		}

		var search map[string]interface{}
		switch memberType {
		case MemberTypeAttribute:
			search = data.Attributes
		case MemberTypeMeta:
			search = data.Meta
		}

		switch resourceValue.Kind() {
		case reflect.String:
			search[memberName] = resourceValue.Interface().(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			search[memberName] = resourceValue.Interface().(int)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			search[memberName] = resourceValue.Interface().(uint)
		case reflect.Float32, reflect.Float64:
			search[memberName] = resourceValue.Interface().(float64)
		default:
			// TODO
		}
	}

	document := Document{
		JSONAPI: JSONAPI{
			Version: "1.0",
		},
		Data: data,
	}
	return json.MarshalIndent(&document, "", "\t")
}
