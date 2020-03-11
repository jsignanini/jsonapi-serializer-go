package jsonapi

import (
	"fmt"
	"reflect"
)

type Resource struct {
	// Exception: The id member is not required when the resource object originates at the client and represents a new resource to be created on the server.
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`

	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
	Links         Links         `json:"links,omitempty"`
	Meta          Meta          `json:"meta,omitempty"`
}

func NewResource() *Resource {
	return &Resource{
		Attributes: Attributes{},
		Meta:       Meta{},
	}
}

// TODO warn or error out when ID isn't plural?
func (r *Resource) SetIDAndType(idValue reflect.Value, resourceType string) error {
	kind := idValue.Kind()
	if kind != reflect.String {
		fmt.Println("here")
		return fmt.Errorf("ID must be a string, got %s", kind)
	}
	id, _ := idValue.Interface().(string)
	if id == "" {
		return fmt.Errorf("ID must be set")
	}
	r.ID = id
	r.Type = resourceType
	return nil
}

func (r *Resource) SetLinks(linksValue reflect.Value) error {
	links, ok := linksValue.Interface().(Links)
	if !ok {
		return fmt.Errorf("field tagged as link needs to be of jsonapi.Links type")
	}
	r.Links = links
	return nil
}
