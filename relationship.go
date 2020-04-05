package jsonapi

// Relationships is a map of JSON:API relationship objects.
type Relationships map[string]interface{}

// RelationshipLink is a JSON:API relationship links object.
// See https://jsonapi.org/format/#document-resource-object-related-resource-links.
type RelationshipLink struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

type relationship struct {
	Links *RelationshipLink `json:"links,omitempty"`
	Meta  *Meta             `json:"meta,omitempty"`
}

// Relationship struct { is a JSON:API relationship object.
// See https://jsonapi.org/format/#document-resource-object-relationships.
type Relationship struct {
	Data *Resource `json:"data"`
	relationship
}

// CompoundRelationship is a JSON:API compound relationship object.
// See https://jsonapi.org/format/#document-resource-object-relationships.
type CompoundRelationship struct {
	Data []*Resource `json:"data"`
	relationship
}

// NewRelationship generates a new JSON:API relationship object.
func NewRelationship() *Relationship {
	return &Relationship{
		relationship: relationship{
			// Links: &RelationshipLink{},
			// Meta:  &Meta{},
		},
	}
}

// AddResource adds a new resource to an existing realationship.
func (r *Relationship) AddResource(resource *Resource) {
	r.Data = resource
}

// NewCompoundRelationship generates a new JSON:API compound relationship object.
func NewCompoundRelationship() *CompoundRelationship {
	return &CompoundRelationship{
		Data:         []*Resource{},
		relationship: relationship{
			// Links: &RelationshipLink{},
			// Meta:  &Meta{},
		},
	}
}
