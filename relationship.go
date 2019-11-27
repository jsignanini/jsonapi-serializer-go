package jsonapi

type Relationships map[string]Relationship

type Relationship struct {
	Data  Resource         `json:"data,omitempty"`
	Links RelationshipLink `json:"links,omitempty"`
	Meta  Meta             `json:"meta,omitempty"`
}

type RelationshipLink struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}
