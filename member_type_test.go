package jsonapi

import (
	"fmt"
	"testing"
)

func TestNewMemberType(t *testing.T) {
	valids := map[string]MemberType{
		"attribute":    MemberTypeAttribute,
		"links":        MemberTypeLinks,
		"primary":      MemberTypePrimary,
		"meta":         MemberTypeMeta,
		"relationship": MemberTypeRelationship,
	}
	invalids := []string{"attributes", "link", "foo", ""}
	for s, v := range valids {
		m, err := NewMemberType(s)
		if err != nil {
			t.Fatal(err)
		}
		if m != v {
			t.Errorf("expected string: %s, to be member type %s, got: %s", s, v, m)
		}
	}
	for _, s := range invalids {
		if _, err := NewMemberType(s); err == nil {
			t.Errorf("expected string: %s, to error out with message: %s", s, fmt.Errorf("MemberType '%s' not found.", s))
		}
	}
}
