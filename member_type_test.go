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

func TestGetMember(t *testing.T) {
	type Correct struct {
		ID      string `jsonapi:"primary,corrects"`
		Correct string `jsonapi:"attribute,correct"`
	}
	c := Correct{
		ID:      "correct-id",
		Correct: "correct",
	}
	if _, err := Marshal(&c, nil); err != nil {
		t.Errorf("got unexpected error: %s, for correct member: %s", err, "attribute")
	}

	type Incorrect struct {
		ID        string `jsonapi:"primary,incorrects"`
		Incorrect string `jsonapi:"foo,empty"`
	}
	i := Incorrect{
		ID:        "incorrect-id",
		Incorrect: "incorrect",
	}
	if _, err := Marshal(&i, nil); err == nil {
		t.Errorf("expected incorrect member: %s, to error out", "foo")
	}

	type Empty struct {
		ID    string `jsonapi:"primary,empties"`
		Empty string `jsonapi:""`
	}
	e := Empty{
		ID:    "empty-id",
		Empty: "empty",
	}
	if _, err := Marshal(&e, nil); err == nil {
		t.Errorf("expected empty tag error: %s, but got no error", fmt.Errorf("tag: %s, not specified", tagKey))
	}
}
