package jsonapi

type Links map[string]interface{}

type LinkString string

type LinkObject struct {
	HREF string `json:"href"`
	Meta Meta   `json:"meta,omitempty"`
}

func (l Links) AddLink(name, url string) {
	l[name] = LinkString(url)
}

func (l Links) AddLinkWithMeta(name, url string, meta Meta) {
	l[name] = LinkObject{
		HREF: url,
		Meta: meta,
	}
}
