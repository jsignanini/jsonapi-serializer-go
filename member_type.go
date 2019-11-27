package jsonapi

import "fmt"

type MemberType string

const (
	MemberTypeAttribute MemberType = "attribute"
	MemberTypeMeta      MemberType = "meta"
	MemberTypePrimary   MemberType = "primary"
)

func NewMemberType(s string) (MemberType, error) {
	switch s {
	case "attribute":
		return MemberTypeAttribute, nil
	case "meta":
		return MemberTypeMeta, nil
	case "primary":
		return MemberTypePrimary, nil
	default:
		return "", fmt.Errorf("MemberType '%s' not found.", s)
	}
}
