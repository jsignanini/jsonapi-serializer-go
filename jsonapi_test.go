package jsonapi

import (
	"bytes"
	"testing"
)

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
	SetTagKey("customKey")
	type Article struct {
		ID    string `customKey:"primary,articles"`
		Title string `customKey:"attribute,title"`
	}
	article := Article{
		ID:    "article-id",
		Title: "Hello World!",
	}
	articleExpected := []byte(`{
	"data": {
		"id": "article-id",
		"type": "articles",
		"attributes": {
			"title": "Hello World!"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if got, err := Marshal(&article, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(got, articleExpected) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(articleExpected), string(got))
		}
	}
	// TODO fix this so we don't have to reset the key
	SetTagKey("jsonapi")
}
