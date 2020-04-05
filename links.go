package jsonapi

// Links is a JSON:API links object.
// See https://jsonapi.org/format/#document-links.
type Links map[string]interface{}

type linkString string

type linkObject struct {
	HREF string `json:"href"`
	Meta Meta   `json:"meta,omitempty"`
}

// AddLink adds a url-only link.
func (l Links) AddLink(name, url string) {
	l[name] = linkString(url)
}

// AddLinkWithMeta adds a link with a meta object.
func (l Links) AddLinkWithMeta(name, url string, meta Meta) {
	l[name] = linkObject{
		HREF: url,
		Meta: meta,
	}
}
