package jsonapi

import "testing"

func TestSetJSONPrefix(t *testing.T) {
	def := ""
	if def != jsonPrefix {
		t.Errorf("default jsonPrefix was incorrect, got: %s, want: %s.", jsonPrefix, "")
	}
	space := " "
	SetJSONPrefix(space)
	if space != jsonPrefix {
		t.Errorf("space jsonPrefix was incorrect, got: %s, want: %s.", jsonPrefix, " ")
	}
	SetJSONPrefix(def)
}

func TestSetJSONIndent(t *testing.T) {
	def := "\t"
	if def != jsonIndent {
		t.Errorf("default jsonIndent was incorrect, got: %s, want: %s.", jsonIndent, "\t")
	}
	space := " "
	SetJSONIndent(space)
	if space != jsonIndent {
		t.Errorf("space jsonIndent was incorrect, got: %s, want: %s.", jsonIndent, " ")
	}
	SetJSONIndent(def)
}

func TestSetTagKey(t *testing.T) {
	// TODO
}
