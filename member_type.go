package jsonapi

import (
	"fmt"
	"reflect"
	"strings"
)

type MemberType string

const (
	MemberTypeAttribute MemberType = "attribute"
	MemberTypeLinks     MemberType = "links"
	MemberTypeMeta      MemberType = "meta"
	MemberTypePrimary   MemberType = "primary"
)

func NewMemberType(s string) (MemberType, error) {
	switch s {
	case "attribute":
		return MemberTypeAttribute, nil
	case "links":
		return MemberTypeLinks, nil
	case "meta":
		return MemberTypeMeta, nil
	case "primary":
		return MemberTypePrimary, nil
	default:
		return "", fmt.Errorf("MemberType '%s' not found.", s)
	}
}
func getMember(field reflect.StructField) (MemberType, string, error) {
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
	memberType, err := NewMemberType(tagParts[0])
	if err != nil {
		return "", "", err
	}
	return memberType, tagParts[1], nil
}
