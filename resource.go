package jsonapi

import (
	"fmt"
	"reflect"
)

// Resource is a JSON:API resource object.
// See https://jsonapi.org/format/#document-resource-objects.
type Resource struct {
	// Exception: The id member is not required when the resource object originates at the client and represents a new resource to be created on the server.
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`

	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
	Links         Links         `json:"links,omitempty"`
	Meta          Meta          `json:"meta,omitempty"`
}

// NewResource generates a new JSON:API resource object.
func NewResource() *Resource {
	return &Resource{
		Attributes: Attributes{},
		Links:      Links{},
		Meta:       Meta{},
	}
}

// SetIDAndType sets the id and type of a JSON:API resource object.
// TODO warn or error out when ID isn't plural?
func (r *Resource) SetIDAndType(idValue reflect.Value, resourceType string) error {
	kind := idValue.Kind()
	if kind != reflect.String {
		return fmt.Errorf("ID must be a string, got %s", kind)
	}
	id, _ := idValue.Interface().(string)
	if id == "" {
		return fmt.Errorf("ID must be set")
	}
	if resourceType == "" {
		return fmt.Errorf("type must be set")
	}
	r.ID = id
	r.Type = resourceType
	return nil
}

// SetLinks sets the links of a JSON:API resource object.
func (r *Resource) SetLinks(linksValue reflect.Value) error {
	links, ok := linksValue.Interface().(Links)
	if !ok {
		return fmt.Errorf("field tagged as link needs to be of jsonapi.Links type")
	}
	r.Links = links
	return nil
}
