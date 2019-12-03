package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	document := NewDocument()
	if err := iterateStruct(document, v, func(field reflect.StructField, value reflect.Value, memberNames ...string) error {
		// get member info, continue otherwise
		memberType, memberName, err := getMember(field)
		if err != nil {
			return nil
		}

		kind := value.Kind()

		if memberType == MemberTypePrimary {
			if kind != reflect.String {
				return fmt.Errorf("ID must be a string")
			}
			id, _ := value.Interface().(string)
			if id == "" {
				return nil
			}
			document.Data.ID = id
			document.Data.Type = memberName
			return nil
		}

		var search map[string]interface{}
		switch memberType {
		case MemberTypeAttribute:
			search = document.Data.Attributes
		case MemberTypeMeta:
			search = document.Data.Meta
		}

		switch kind {
		case reflect.String:
			search[memberName] = value.Interface().(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			search[memberName] = value.Interface().(int)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			search[memberName] = value.Interface().(uint)
		case reflect.Float32, reflect.Float64:
			search[memberName] = value.Interface().(float64)
		default:
			// TODO
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return json.MarshalIndent(&document, "", "\t")
}
