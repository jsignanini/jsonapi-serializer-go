package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	document := NewDocument()
	if err := iterateStruct(document, v, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
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
			document.Data.Type = memberNames[0]
			return nil
		}

		marshal(document, memberType, memberNames, value)

		return nil
	}); err != nil {
		return nil, err
	}

	return json.MarshalIndent(&document, "", "\t")
}

// document *Document, memberType MemberType, memberNames ...string
func marshal(document *Document, memberType MemberType, memberNames []string, value reflect.Value) {
	// figure out search
	var search map[string]interface{}
	switch memberType {
	case MemberTypeAttribute:
		search = document.Data.Attributes
	case MemberTypeMeta:
		search = document.Data.Meta
	}

	// iterate memberNames
	var memberName string
	for i, name := range memberNames {
		memberName = name
		if i == len(memberNames)-1 {
			break
		}
		if _, ok := search[name]; !ok {
			search[name] = make(map[string]interface{})
		}
		search = search[name].(map[string]interface{})
	}

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// set value
	switch value.Kind() {
	case reflect.Bool:
		search[memberName] = value.Interface().(bool)
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
		// cu, ok := customUnmarshalers[rv.Type()]
		// if !ok {
		// 	return fmt.Errorf("Type not supported, must implement custom unmarshaller")
		// }
		// cu(v, rv)
	}
	// return nil
}
