# jsonapi-serializer-go

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/jsignanini/jsonapi-serializer-go)
[![Build Status](https://travis-ci.org/jsignanini/jsonapi-serializer-go.svg?branch=master)](https://travis-ci.org/jsignanini/jsonapi-serializer-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsignanini/jsonapi-serializer-go)](https://goreportcard.com/report/github.com/jsignanini/jsonapi-serializer-go)
[![Coverage Status](https://coveralls.io/repos/github/jsignanini/jsonapi-serializer-go/badge.svg?branch=master)](https://coveralls.io/github/jsignanini/jsonapi-serializer-go?branch=master)

### TODOs

- many payload
- many payload of pointers
- inferred tags (e.g. no `jsonapi:"..."` tag, assume it is an attribute and infer name from field name)
- ignored fields (e.g. `jsonapi:"-"`)
- errors
- support all native types
- support structs
- support embedded structs
- support nested structs
- support custom types with custom un/marshallers
- jsonapi spec validation
- jsonapi settings (e.g.: spec version, error/warning on document validation, etc.)
- support omitempty tag
- add overflow check and tests for int, uint and float (both value and pointers)
