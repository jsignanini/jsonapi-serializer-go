package jsonapi

var (
	jsonPrefix = ""
	jsonIndent = "\t"
	tagKey     = "jsonapi"
)

// SetJSONPrefix sets the prefix value for json.MarshalIndent.
// See https://golang.org/pkg/encoding/json/#MarshalIndent.
func SetJSONPrefix(prefix string) {
	jsonPrefix = prefix
}

// SetJSONIndent sets the indent value for json.MarshalIndent.
// See // See https://golang.org/pkg/encoding/json/#MarshalIndent.
func SetJSONIndent(indent string) {
	jsonIndent = indent
}

// SetTagKey sets a custom value for the JSON:API tag key.
func SetTagKey(key string) {
	tagKey = key
}
