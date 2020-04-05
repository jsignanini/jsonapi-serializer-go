package jsonapi

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewMemberType(t *testing.T) {
	valids := map[string]memberType{
		"attribute":    memberTypeAttribute,
		"links":        memberTypeLinks,
		"primary":      memberTypePrimary,
		"meta":         memberTypeMeta,
		"relationship": memberTypeRelationship,
	}
	invalids := []string{"attributes", "link", "foo", ""}
	for s, v := range valids {
		m, err := newMemberType(s)
		if err != nil {
			t.Fatal(err)
		}
		if m != v {
			t.Errorf("expected string: %s, to be member type %s, got: %s", s, v, m)
		}
	}
	for _, s := range invalids {
		if _, err := newMemberType(s); err == nil {
			t.Errorf("expected string: %s, to error out with message: %s", s, fmt.Errorf("member type '%s' not found", s))
		}
	}
}

func TestGetMember(t *testing.T) {
	type GetMemterTest struct {
		ID        string `jsonapi:"primary,corrects"`
		Correct   string `jsonapi:"attribute,correct"`
		Incorrect string `jsonapi:"foo,empty"`
		Malformed string `jsonapi:"foo,bar,malformed"`
		Empty     string `jsonapi:""`
		NoTag     string
	}
	test := GetMemterTest{
		ID:        "correct-id",
		Correct:   "correct",
		Incorrect: "incorrect",
		Empty:     "empty",
	}

	if c, ok := reflect.TypeOf(test).FieldByName("Correct"); !ok {
		t.Fatal("not ok")
	} else {
		if _, _, err := getMember(c); err != nil {
			t.Errorf("got unexpected error: %s, for correct member: %s", err, "attribute")
		}
	}
	if i, ok := reflect.TypeOf(test).FieldByName("Incorrect"); !ok {
		t.Fatal("not ok")
	} else {
		if _, _, err := getMember(i); err == nil {
			t.Errorf("expected incorrect member: %s, to error out", "foo")
		}
	}
	if e, ok := reflect.TypeOf(test).FieldByName("Malformed"); !ok {
		t.Fatal("not ok")
	} else {
		if _, _, err := getMember(e); err == nil {
			t.Errorf("expected malformed tag error: %s, but got no error", fmt.Errorf("tag: %s, was not formatted properly", tagKey))
		}
	}
	if e, ok := reflect.TypeOf(test).FieldByName("Empty"); !ok {
		t.Fatal("not ok")
	} else {
		if _, _, err := getMember(e); err == nil {
			t.Errorf("expected empty tag error: %s, but got no error", fmt.Errorf("tag: %s, not specified", tagKey))
		}
	}
	if e, ok := reflect.TypeOf(test).FieldByName("NoTag"); !ok {
		t.Fatal("not ok")
	} else {
		if _, _, err := getMember(e); err == nil {
			t.Errorf("expected tag missing error: %s, but got no error", fmt.Errorf("tag: %s, not specified", tagKey))
		}
	}
	if _, err := Marshal(&test, nil); err == nil {
		t.Errorf("expected getMember error but got no error")
	}
}
