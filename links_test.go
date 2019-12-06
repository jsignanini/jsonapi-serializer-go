package jsonapi

import (
	"bytes"
	"testing"
)

func TestDocumentLinks(t *testing.T) {
	type TestLinks struct {
		ID  string `jsonapi:"primary,test_links"`
		Foo string `jsonapi:"attribute,bar"`
	}
	links := Links{}
	links.AddLink("self", "https://example.com")
	links.AddLinkWithMeta("related", "https://example.com", Meta{"foo": "bar"})

	t1 := TestLinks{
		ID:  "someID",
		Foo: "hello world!",
	}
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_links",
		"attributes": {
			"bar": "hello world!"
		}
	},
	"jsonapi": {
		"version": "1.0"
	},
	"links": {
		"related": {
			"href": "https://example.com",
			"meta": {
				"foo": "bar"
			}
		},
		"self": "https://example.com"
	}
}`)
	if b, err := Marshal(&t1, &MarshalParams{
		Links: &links,
	}); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestResourceLinks(t *testing.T) {
	type TestLinks struct {
		ID      string `jsonapi:"primary,test_links"`
		Foo     string `jsonapi:"attribute,bar"`
		MyLinks Links  `jsonapi:"links,"`
	}
	t1 := TestLinks{
		ID:      "someID",
		Foo:     "hello world!",
		MyLinks: Links{},
	}
	t1.MyLinks.AddLink("self", "https://resource.com")
	t1.MyLinks.AddLinkWithMeta("related", "https://resource.com", Meta{"foo": "bar"})
	expected := []byte(`{
	"data": {
		"id": "someID",
		"type": "test_links",
		"attributes": {
			"bar": "hello world!"
		},
		"links": {
			"related": {
				"href": "https://resource.com",
				"meta": {
					"foo": "bar"
				}
			},
			"self": "https://resource.com"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&t1, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestResourceWithEmbeddedLinks(t *testing.T) {
	type Car struct {
		ID   string `jsonapi:"primary,cars"`
		Make string `jsonapi:"attribute,make"`
	}
	type CarWithLinks struct {
		Car
		Links `jsonapi:"links,"`
	}
	c := CarWithLinks{
		Car: Car{
			ID:   "VIN1192392348",
			Make: "Honda",
		},
		Links: Links{},
	}
	c.Links.AddLink("self", "https://honda.com/VIN1192392348")
	c.Links.AddLinkWithMeta("related", "https://lexus.com", Meta{"luxury": true})
	expected := []byte(`{
	"data": {
		"id": "VIN1192392348",
		"type": "cars",
		"attributes": {
			"make": "Honda"
		},
		"links": {
			"related": {
				"href": "https://lexus.com",
				"meta": {
					"luxury": true
				}
			},
			"self": "https://honda.com/VIN1192392348"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`)
	if b, err := Marshal(&c, nil); err != nil {
		t.Errorf(err.Error())
	} else {
		if bytes.Compare(expected, b) != 0 {
			t.Errorf("Expected:\n%s\nGot:\n%s\n", string(expected), string(b))
		}
	}
}

func TestResourceLinksWrongType(t *testing.T) {
	type TestLinks struct {
		ID      string `jsonapi:"primary,test_links"`
		Foo     string `jsonapi:"attribute,bar"`
		MyLinks string `jsonapi:"links,"`
	}
	t1 := TestLinks{
		ID:      "someID",
		Foo:     "hello world!",
		MyLinks: "should have been a Links type",
	}
	if _, err := Marshal(&t1, nil); err == nil {
		t.Errorf("Expected wrong type error when links tag is set to a non Links type")
	}
}
