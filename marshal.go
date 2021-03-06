package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// MarshalParams are the optional parameters to add links and meta objects to a top-level document.
type MarshalParams struct {
	Links *Links
	Meta  *Meta
}

// Marshal returns the JSON:API encoding of v.
func Marshal(v interface{}, p *MarshalParams) ([]byte, error) {
	rType := reflect.TypeOf(v)

	// only allow pointer or slice kind
	kind := rType.Kind()
	if kind != reflect.Ptr && kind != reflect.Slice {
		return nil, fmt.Errorf("v must be pointer or slice")
	}

	// determine if v is a slice
	isSlice := false
	if rType.Elem().Kind() == reflect.Slice {
		isSlice = true
	}

	// handle compound document
	if isSlice {
		ncdp := &NewCompoundDocumentParams{}
		if p != nil {
			ncdp.Links = p.Links
			ncdp.Meta = p.Meta
		}
		document := NewCompoundDocument(ncdp)
		return marshalCompoundDocument(v, document)
	}

	// handle single document
	ndp := &NewDocumentParams{}
	if p != nil {
		ndp.Links = p.Links
		ndp.Meta = p.Meta
	}
	document := NewDocument(ndp)
	return marshalDocument(v, document)
}

// RegisterMarshaler register a custom marshaller function for a t type.
func RegisterMarshaler(t reflect.Type, u marshalerFunc) {
	customMarshalers[t] = u
}

type marshalerFunc = func(map[string]interface{}, string, reflect.Value)

var customMarshalers = make(map[reflect.Type]marshalerFunc)

func marshalDocument(v interface{}, d *Document) ([]byte, error) {
	d.Data = NewResource()
	if err := iterateStruct(v, func(value reflect.Value, memberType memberType, memberNames ...string) error {
		switch memberType {
		case memberTypePrimary:
			return d.Data.SetIDAndType(value, memberNames[0])
		case memberTypeLinks:
			return d.Data.SetLinks(value)
		case memberTypeRelationship:
			if d.Data.Relationships == nil {
				d.Data.Relationships = Relationships{}
			}
			relIsSlice := false
			if value.Kind() == reflect.Slice {
				relIsSlice = true
			}
			if !relIsSlice {
				if err := marshalRelationship(value, &d.document, d.Data, memberNames); err != nil {
					return err
				}
			} else {
				if err := marshalCompoundRelationship(value, &d.document, d.Data, memberNames); err != nil {
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
			return nil, fmt.Errorf("document must be pointer or slice of pointers")
		}
		r := NewResource()
		if err := iterateStruct(value.Interface(), func(value reflect.Value, memberType memberType, memberNames ...string) error {
			switch memberType {
			case memberTypePrimary:
				return r.SetIDAndType(value, memberNames[0])
			case memberTypeLinks:
				return r.SetLinks(value)
			case memberTypeRelationship:
				if r.Relationships == nil {
					r.Relationships = Relationships{}
				}
				relIsSlice := false
				if value.Kind() == reflect.Slice {
					relIsSlice = true
				}
				if !relIsSlice {
					if err := marshalRelationship(value, &cd.document, r, memberNames); err != nil {
						return err
					}
				} else {
					if err := marshalCompoundRelationship(value, &cd.document, r, memberNames); err != nil {
						return err
					}
				}
				return nil
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

func marshalRelationship(value reflect.Value, d *document, r *Resource, memberNames []string) error {
	rel := NewRelationship()
	r.Relationships[memberNames[0]] = rel
	if value.IsNil() {
		return nil
	}
	newIncl := NewResource()
	rel.AddResource(NewResource())
	if err := iterateStruct(value.Interface(), func(v2 reflect.Value, memberType memberType, memberNames ...string) error {
		switch memberType {
		case memberTypePrimary:
			if err := rel.Data.SetIDAndType(v2, memberNames[0]); err != nil {
				return err
			}
			return newIncl.SetIDAndType(v2, memberNames[0])
		case memberTypeLinks:
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

func marshalCompoundRelationship(value reflect.Value, d *document, r *Resource, memberNames []string) error {
	rels := NewCompoundRelationship()
	r.Relationships[memberNames[0]] = rels
	for i := 0; i < value.Len(); i++ {
		sValue := value.Index(i)
		if sValue.Kind() != reflect.Ptr {
			return fmt.Errorf("relationship must be pointer or slice of pointers")
		}
		newIncl := NewResource()
		newRel := NewRelationship()
		newRel.AddResource(NewResource())
		if err := iterateStruct(sValue.Interface(), func(v2 reflect.Value, memberType memberType, memberNames ...string) error {
			switch memberType {
			case memberTypePrimary:
				if err := newRel.Data.SetIDAndType(v2, memberNames[0]); err != nil {
					return err
				}
				return newIncl.SetIDAndType(v2, memberNames[0])
			case memberTypeLinks:
				return newIncl.SetLinks(v2)
			default:
				return marshal(newIncl, memberType, memberNames, v2)
			}
		}); err != nil {
			return err
		}
		// make sure it's only added once to included
		existsInIncluded := false
		for _, incl := range d.Included {
			if incl.Type == newIncl.Type && incl.ID == newIncl.ID {
				existsInIncluded = true
			}
		}
		if !existsInIncluded {
			d.Included = append(d.Included, newIncl)
		}
		rels.Data = append(rels.Data, newRel.Data)
	}
	return nil
}

func marshal(resource *Resource, memberType memberType, memberNames []string, value reflect.Value) error {
	// figure out search
	var search map[string]interface{}
	switch memberType {
	case memberTypeAttribute:
		search = resource.Attributes
	case memberTypeMeta:
		search = resource.Meta
	case memberTypeRelationship:
		// TODO we should not need to skip relationships, this function should not be called for relationships
		return nil
	}

	// iterate memberNames
	var memberName string
	for i, name := range memberNames {
		memberName = name
		if i == len(memberNames)-1 {
			break
		}
		// TODO do we still need this
		// if _, ok := search[name]; !ok {
		// 	search[name] = make(map[string]interface{})
		// }
		// search = search[name].(map[string]interface{})
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
		reflect.Slice,
		reflect.String,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		search[memberName] = value.Interface()
		// TODO handle error/warning for unsupported types
		// default:
		// 	return fmt.Errorf("type: %+v, not supported, must implement custom marshaller", value.Type())
	}
	return nil
}
