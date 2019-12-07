package jsonapi

var (
	jsonPrefix = ""
	jsonIndent = "\t"
	tagKey     = "jsonapi"
)

func SetJSONPrefix(prefix string) {
	jsonPrefix = prefix
}

func SetJSONIndent(indent string) {
	jsonIndent = indent
}

func SetTagKey(key string) {
	tagKey = key
}
