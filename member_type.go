package jsonapi

import (
	"fmt"
	"reflect"
	"strings"
)

// memberType is a JSONP:API member type.
type memberType string

const (
	memberTypeAttribute    memberType = "attribute"
	memberTypeLinks        memberType = "links"
	memberTypeMeta         memberType = "meta"
	memberTypePrimary      memberType = "primary"
	memberTypeRelationship memberType = "relationship"
)

func newMemberType(s string) (memberType, error) {
	switch s {
	case "attribute":
		return memberTypeAttribute, nil
	case "links":
		return memberTypeLinks, nil
	case "meta":
		return memberTypeMeta, nil
	case "primary":
		return memberTypePrimary, nil
	case "relationship":
		return memberTypeRelationship, nil
	default:
		return "", fmt.Errorf("member type '%s' not found", s)
	}
}

func getMember(field reflect.StructField) (memberType, string, error) {
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
	memberType, err := newMemberType(tagParts[0])
	if err != nil {
		return "", "", err
	}
	return memberType, tagParts[1], nil
}
