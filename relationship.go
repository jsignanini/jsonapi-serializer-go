package jsonapi

type Relationships map[string]interface{}

// type Relationship struct {
// 	Data  *Resource         `json:"data,omitempty"`
// 	Links *RelationshipLink `json:"links,omitempty"`
// 	Meta  *Meta             `json:"meta,omitempty"`
// }

type RelationshipLink struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// func NewRelationship() *Relationship {
// 	r := NewResource()
// 	return &Relationship{
// 		Data: r,
// 		// Meta: &Meta{},
// 	}
// }

type relationship struct {
	Links *RelationshipLink `json:"links,omitempty"`
	Meta  *Meta             `json:"meta,omitempty"`
}

type Relationship struct {
	Data *Resource `json:"data,omitempty"`
	relationship
}

type CompoundRelationship struct {
	Data []*Resource `json:"data,omitempty"`
	relationship
}

func NewRelationship() *Relationship {
	r := NewResource()
	return &Relationship{
		Data:         r,
		relationship: relationship{
			// Links: &RelationshipLink{},
			// Meta:  &Meta{},
		},
	}
}

func NewCompoundRelationship() *CompoundRelationship {
	return &CompoundRelationship{
		Data:         []*Resource{},
		relationship: relationship{
			// Links: &RelationshipLink{},
			// Meta:  &Meta{},
		},
	}
}
