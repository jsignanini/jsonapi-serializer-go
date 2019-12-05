package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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
	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Ptr && kind != reflect.Slice {
		return fmt.Errorf("v should be pointer or slice")
	}

	isSlice := false
	if reflect.TypeOf(v).Elem().Kind() == reflect.Slice {
		isSlice = true
	}

	if !isSlice {
		document := NewDocument()
		if err := json.Unmarshal(data, document); err != nil {
			return err
		}

		if err := iterateStruct(v, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
			fieldKind := value.Kind()
			// TODO this sets ID for all nexted primary tag fields
			if memberType == MemberTypePrimary {
				if fieldKind != reflect.String {
					return fmt.Errorf("ID must be a string")
				}
				value.SetString(document.Data.ID)
				return nil
			}

			// set raw value
			return unmarshal(document.Data, memberType, memberNames, value)
		}); err != nil {
			return err
		}
	} else {
		document := NewCompoundDocument()
		if err := json.Unmarshal(data, document); err != nil {
			return err
		}

		for _, resource := range document.Data {
			v2 := reflect.New(reflect.ValueOf(v).Elem().Type().Elem()).Interface()
			if err := iterateStruct(v2, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
				fieldKind := value.Kind()
				// TODO this sets ID for all nexted primary tag fields
				if memberType == MemberTypePrimary {
					if fieldKind != reflect.String {
						return fmt.Errorf("ID must be a string")
					}
					value.SetString(resource.ID)
					return nil
				}

				// set raw value
				return unmarshal(resource, memberType, memberNames, value)
			}); err != nil {
				return err
			}
			value := reflect.ValueOf(v).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(v2).Elem()))
		}
	}

	return nil
}

func RegisterUnmarshaler(t reflect.Type, u unmarshalerFunc) {
	customUnmarshalers[t] = u
}

type unmarshalerFunc = func(interface{}, reflect.Value)

var customUnmarshalers = make(map[reflect.Type]unmarshalerFunc)

func unmarshal(resource *Resource, memberType MemberType, memberNames []string, value reflect.Value) error {
	// find raw value if exists
	var search map[string]interface{}
	switch memberType {
	case MemberTypeAttribute:
		search = resource.Attributes
	case MemberTypeMeta:
		search = resource.Meta
	}
	rawValue, found := deepSearch(search, memberNames...)
	if !found {
		return nil
	}

	// if pointer, get non-pointer kind
	isPtr := false
	kind := value.Kind()
	if kind == reflect.Ptr {
		isPtr = true
		kind = reflect.New(value.Type().Elem()).Elem().Kind()
	}

	// set values by kind
	switch kind {
	case reflect.Bool:
		if val, ok := rawValue.(bool); ok {
			if isPtr {
				value.Set(reflect.ValueOf(&val))
			} else {
				value.SetBool(val)
			}
		}
	case reflect.String:
		if val, ok := rawValue.(string); ok {
			if isPtr {
				value.Set(reflect.ValueOf(&val))
			} else {
				value.SetString(val)
			}
		}
	case reflect.Int:
		if val, ok := rawValue.(float64); ok {
			intVal := int(val)
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&intVal))
			} else {
				value.SetInt(int64(intVal))
			}
		}
	case reflect.Int8:
		if val, ok := rawValue.(float64); ok {
			intVal := int8(val)
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&intVal))
			} else {
				value.SetInt(int64(intVal))
			}
		}
	case reflect.Int16:
		if val, ok := rawValue.(float64); ok {
			intVal := int16(val)
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&intVal))
			} else {
				value.SetInt(int64(intVal))
			}
		}
	case reflect.Int32:
		if val, ok := rawValue.(float64); ok {
			intVal := int32(val)
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&intVal))
			} else {
				value.SetInt(int64(intVal))
			}
		}
	case reflect.Int64:
		if val, ok := rawValue.(float64); ok {
			intVal := int64(val)
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&intVal))
			} else {
				value.SetInt(intVal)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if val, ok := rawValue.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&val))
			} else {
				value.SetUint(uint64(val))
			}
		}
	case reflect.Float32, reflect.Float64:
		if val, ok := rawValue.(float64); ok {
			// TODO resourceValue.OverflowInt(val)
			if isPtr {
				value.Set(reflect.ValueOf(&val))
			} else {
				value.SetFloat(val)
			}
		}
	default:
		cu, ok := customUnmarshalers[value.Type()]
		if !ok {
			return fmt.Errorf("Type: %+v, not supported, must implement custom unmarshaller", value.Type())
		}
		cu(rawValue, value)
	}
	return nil
}

func deepSearch(tree map[string]interface{}, keys ...string) (interface{}, bool) {
	key, keys := keys[0], keys[1:]
	value, ok := tree[key]
	if !ok {
		return nil, false
	}
	if len(keys) == 0 {
		return value, true
	}
	return deepSearch(tree[key].(map[string]interface{}), keys...)
}
