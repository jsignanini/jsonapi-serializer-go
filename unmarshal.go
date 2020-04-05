package jsonapi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	rType := reflect.TypeOf(v)
	rValue := reflect.ValueOf(v)
	kind := rType.Kind()

	// v must be pointer
	if kind != reflect.Ptr {
		return fmt.Errorf("v must be pointer")
	}

	// v must not be nil
	if rValue.IsNil() {
		return fmt.Errorf("v must not be nil")
	}

	// determine if v is a slice
	isSlice := false
	if rType.Elem().Kind() == reflect.Slice {
		isSlice = true
	}

	if isSlice {
		document := NewCompoundDocument(nil)
		if err := json.Unmarshal(data, document); err != nil {
			return err
		}
		return unmarshalCompoundDocument(v, document)
	} else {
		document := NewDocument(nil)
		if err := json.Unmarshal(data, document); err != nil {
			return err
		}
		return unmarshalDocument(v, document)
	}
}

func RegisterUnmarshaler(t reflect.Type, u unmarshalerFunc) {
	customUnmarshalers[t] = u
}

type unmarshalerFunc = func(interface{}, reflect.Value)

var customUnmarshalers = make(map[reflect.Type]unmarshalerFunc)

func unmarshalCompoundDocument(v interface{}, cd *CompoundDocument) error {
	rValue := reflect.ValueOf(v)
	for _, resource := range cd.Data {
		v2 := reflect.New(rValue.Elem().Type().Elem()).Interface()
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
		value := rValue.Elem()
		value.Set(reflect.Append(value, reflect.ValueOf(v2).Elem()))
	}
	return nil
}

func unmarshalDocument(v interface{}, d *Document) error {
	return iterateStruct(v, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
		fieldKind := value.Kind()
		// TODO this sets ID for all nexted primary tag fields
		if memberType == MemberTypePrimary {
			if fieldKind != reflect.String {
				return fmt.Errorf("ID must be a string")
			}
			value.SetString(d.Data.ID)
			return nil
		}

		// set raw value
		return unmarshal(d.Data, memberType, memberNames, value)
	})
}

func unmarshal(resource *Resource, memberType MemberType, memberNames []string, field reflect.Value) error {
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

	if cu, ok := customUnmarshalers[field.Type()]; ok {
		cu(rawValue, field)
		return nil
	}

	// if pointer, get non-pointer kind
	if field.Kind() == reflect.Ptr {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}
	value := reflect.Indirect(reflect.ValueOf(rawValue))

	// set values by kind
	switch field.Kind() {
	case reflect.Bool:
		return setBool(field, value)
	case reflect.String:
		return setString(field, value)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return setInt(field, value)
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		return setUint(field, value)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return setFloat(field, value)
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

func setInt(field, value reflect.Value) error {
	bf := new(big.Float)
	float64Str := float64StringFromValue(value)
	if _, err := fmt.Sscan(float64Str, bf); err != nil {
		return err
	}
	i, _ := bf.Int64()
	field.SetInt(i)
	return nil
}

func setUint(field, value reflect.Value) error {
	bf := new(big.Float)
	float64Str := float64StringFromValue(value)
	if _, err := fmt.Sscan(float64Str, bf); err != nil {
		return err
	}
	ui, _ := bf.Uint64()
	field.SetUint(ui)
	return nil
}

func setFloat(field, value reflect.Value) error {
	if _, ok := value.Interface().(float64); !ok {
		return fmt.Errorf("number has no digits")
	}
	field.SetFloat(value.Float())
	return nil
}

func setBool(field, value reflect.Value) error {
	if _, ok := value.Interface().(bool); !ok {
		return fmt.Errorf("invalid value for field %s", field.Type().Name())
	}
	field.SetBool(value.Bool())
	return nil
}

func setString(field, value reflect.Value) error {
	if _, ok := value.Interface().(string); !ok {
		return fmt.Errorf("invalid value for field %s", field.Type().Name())
	}
	field.SetString(value.String())
	return nil
}

func float64StringFromValue(value reflect.Value) string {
	s := fmt.Sprintf("%s", value)
	s = strings.TrimPrefix(s, "%!s(float64=")
	s = strings.TrimSuffix(s, ")")
	return s
}
