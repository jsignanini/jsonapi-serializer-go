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
	rType := reflect.TypeOf(v)

	// only allow pointer or slice kind
	kind := rType.Kind()
	if kind != reflect.Ptr && kind != reflect.Slice {
		return nil, fmt.Errorf("v should be pointer or slice")
	}

	// check if it's a slice
	isSlice := false
	if rType.Elem().Kind() == reflect.Slice {
		isSlice = true
	}

	if isSlice {
		// handle optional params
		ncdp := &NewCompoundDocumentParams{}
		if p != nil {
			ncdp.Links = p.Links
			ncdp.Meta = p.Meta
		}
		document := NewCompoundDocument(ncdp)
		return marshalCompoundDocument(v, document)
	} else {
		// handle optional params
		ndp := &NewDocumentParams{}
		if p != nil {
			ndp.Links = p.Links
			ndp.Meta = p.Meta
		}
		document := NewDocument(ndp)
		return marshalDocument(v, document)
	}
}

func RegisterMarshaler(t reflect.Type, u marshalerFunc) {
	customMarshalers[t] = u
}

type marshalerFunc = func(map[string]interface{}, string, reflect.Value)

var customMarshalers = make(map[reflect.Type]marshalerFunc)

func marshalDocument(v interface{}, d *Document) ([]byte, error) {
	d.Data = NewResource()
	if err := iterateStruct(v, func(value reflect.Value, memberType MemberType, memberNames ...string) error {
		switch memberType {
		case MemberTypePrimary:
			return d.Data.SetIDAndType(value, memberNames[0])
		case MemberTypeLinks:
			return d.Data.SetLinks(value)
		case MemberTypeRelationship:
			if d.Data.Relationships == nil {
				d.Data.Relationships = Relationships{}
			}
			relIsSlice := false
			if value.Kind() == reflect.Slice {
				relIsSlice = true
			}
			if !relIsSlice {
				if err := marshalRelationship(value, d, memberNames); err != nil {
					return err
				}
			} else {
				if err := marshalCompoundRelationship(value, d, memberNames); err != nil {
					return err
				}
			}
			return nil
		default:
			return marshal(d.Data, memberType, memberNames, value)
		}
	}); err != nil {
		return nil, err
	}
	return json.MarshalIndent(&d, jsonPrefix, jsonIndent)
}

func marshalCompoundDocument(v interface{}, cd *CompoundDocument) ([]byte, error) {
	rValue := reflect.ValueOf(v)
	values := rValue.Elem()
	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)
		if value.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("v should be pointer or slice of pointers")
		}
		r := NewResource()
		if err := iterateStruct(value.Interface(), func(value reflect.Value, memberType MemberType, memberNames ...string) error {
			switch memberType {
			case MemberTypePrimary:
				return r.SetIDAndType(value, memberNames[0])
			case MemberTypeLinks:
				return r.SetLinks(value)
			default:
				return marshal(r, memberType, memberNames, value)
			}
		}); err != nil {
			return nil, err
		}
		cd.Data = append(cd.Data, r)
	}
	return json.MarshalIndent(&cd, jsonPrefix, jsonIndent)
}

func marshalRelationship(value reflect.Value, d *Document, memberNames []string) error {
	rel := NewRelationship()
	d.Data.Relationships[memberNames[0]] = rel
	if value.IsNil() {
		return nil
	}
	newIncl := NewResource()
	rel.AddResource(NewResource())
	if err := iterateStruct(value.Interface(), func(v2 reflect.Value, memberType MemberType, memberNames ...string) error {
		switch memberType {
		case MemberTypePrimary:
			if err := rel.Data.SetIDAndType(v2, memberNames[0]); err != nil {
				return err
			}
			return newIncl.SetIDAndType(v2, memberNames[0])
		case MemberTypeLinks:
			return newIncl.SetLinks(v2)
		default:
			return marshal(newIncl, memberType, memberNames, v2)
		}
	}); err != nil {
		return err
	}
	// make sure it's only added once
	for _, incl := range d.Included {
		if incl.Type == newIncl.Type && incl.ID == newIncl.ID {
			return nil
		}
	}
	d.Included = append(d.Included, newIncl)
	return nil
}

func marshalCompoundRelationship(value reflect.Value, d *Document, memberNames []string) error {
	rels := NewCompoundRelationship()
	d.Data.Relationships[memberNames[0]] = rels
	for i := 0; i < value.Len(); i++ {
		sValue := value.Index(i)
		if sValue.Kind() != reflect.Ptr {
			return fmt.Errorf("v should be pointer or slice of pointers")
		}
		newIncl := NewResource()
		newRel := NewRelationship()
		newRel.AddResource(NewResource())
		if err := iterateStruct(sValue.Interface(), func(v2 reflect.Value, memberType MemberType, memberNames ...string) error {
			switch memberType {
			case MemberTypePrimary:
				if err := newRel.Data.SetIDAndType(v2, memberNames[0]); err != nil {
					return err
				}
				return newIncl.SetIDAndType(v2, memberNames[0])
			case MemberTypeLinks:
				return newIncl.SetLinks(v2)
			default:
				return marshal(newIncl, memberType, memberNames, v2)
			}
		}); err != nil {
			return err
		}
		// make sure it's only added once
		for _, incl := range d.Included {
			if incl.Type == newIncl.Type && incl.ID == newIncl.ID {
				return nil
			}
		}
		d.Included = append(d.Included, newIncl)
		rels.Data = append(rels.Data, newRel.Data)
	}
	return nil
}

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

	// ignore nil pointers
	if isPtr && value.IsNil() {
		return nil
	}

	// use custom marshaller if exists
	cm, hasCustomMarshaller := customMarshalers[value.Type()]
	if hasCustomMarshaller {
		cm(search, memberName, value)
		return nil
	}

	// set value
	switch kind {
	case
		reflect.Bool,
		reflect.Complex64, reflect.Complex128,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.String,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		search[memberName] = value.Interface()
	default:
		return fmt.Errorf("type: %+v, not supported, must implement custom marshaller", value.Type())
	}
	return nil
}
