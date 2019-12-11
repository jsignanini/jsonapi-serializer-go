package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type MarshalParams struct {
	Links *Links
	Meta  *Meta
}

func Marshal(v interface{}, p *MarshalParams) ([]byte, error) {
	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Ptr && kind != reflect.Slice {
		return nil, fmt.Errorf("v should be pointer or slice")
	}

	isSlice := false
	if reflect.TypeOf(v).Elem().Kind() == reflect.Slice {
		isSlice = true
	}

	if !isSlice {
		// handle optional params
		document := NewDocument()
		if p != nil && p.Links != nil {
			document.Links = p.Links
		}
		if p != nil && p.Meta != nil {
			document.Meta = p.Meta
		}

		document.Data = NewResource()
		if err := iterateStruct(v, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
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
			if memberType == MemberTypeLinks {
				links, ok := value.Interface().(Links)
				if !ok {
					return fmt.Errorf("field tagged as link needs to be of Links type")
				}
				document.Data.Links = links
				return nil
			}

			return marshal(document.Data, memberType, memberNames, value)
		}); err != nil {
			return nil, err
		}

		return json.MarshalIndent(&document, jsonPrefix, jsonIndent)
	} else {
		document := NewCompoundDocument()
		if p != nil && p.Links != nil {
			document.Links = p.Links
		}
		if p != nil && p.Meta != nil {
			document.Meta = p.Meta
		}
		document.Data = []*Resource{}

		values := reflect.ValueOf(v).Elem()
		for i := 0; i < values.Len(); i++ {
			value := values.Index(i)
			if value.Kind() != reflect.Ptr {
				return nil, fmt.Errorf("v should be pointer or slice of pointers")
			}

			r := NewResource()
			if err := iterateStruct(value.Interface(), func(value reflect.Value, memberType MemberType, memberNames ...string) error {
				kind := value.Kind()

				if memberType == MemberTypePrimary {
					if kind != reflect.String {
						return fmt.Errorf("ID must be a string")
					}
					id, _ := value.Interface().(string)
					if id == "" {
						return nil
					}
					r.ID = id
					r.Type = memberNames[0]
					return nil
				}
				if memberType == MemberTypeLinks {
					links, ok := value.Interface().(Links)
					if !ok {
						return fmt.Errorf("field tagged as link needs to be of Links type")
					}
					r.Links = links
					return nil
				}

				return marshal(r, memberType, memberNames, value)
			}); err != nil {
				return nil, err
			}
			document.Data = append(document.Data, r)
		}

		return json.MarshalIndent(&document, jsonPrefix, jsonIndent)
	}
}

func RegisterMarshaler(t reflect.Type, u marshalerFunc) {
	customMarshalers[t] = u
}

type marshalerFunc = func(map[string]interface{}, string, reflect.Value)

var customMarshalers = make(map[reflect.Type]marshalerFunc)

func marshal(resource *Resource, memberType MemberType, memberNames []string, value reflect.Value) error {
	// figure out search
	var search map[string]interface{}
	switch memberType {
	case MemberTypeAttribute:
		search = resource.Attributes
	case MemberTypeMeta:
		search = resource.Meta
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

	// if pointer, get non-pointer kind
	isPtr := false
	kind := value.Kind()
	if kind == reflect.Ptr {
		isPtr = true
		kind = reflect.New(value.Type().Elem()).Elem().Kind()
	}

	// set value
	switch kind {
	case reflect.Bool:
		// TODO handle pointers in a more generic way
		if isPtr && !value.IsNil() {
			search[memberName] = value.Interface().(*bool)
		} else if !isPtr {
			search[memberName] = value.Interface().(bool)
		}
	case reflect.String:
		if isPtr && !value.IsNil() {
			search[memberName] = value.Interface().(*string)
		} else if !isPtr {
			search[memberName] = value.Interface().(string)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		search[memberName] = value.Interface().(int)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		search[memberName] = value.Interface().(uint)
	case reflect.Float32, reflect.Float64:
		search[memberName] = value.Interface().(float64)
	default:
		cm, ok := customMarshalers[value.Type()]
		if !ok {
			return fmt.Errorf("Type: %+v, not supported, must implement custom marshaller", value.Type())
		}
		cm(search, memberName, value)
	}
	return nil
}
