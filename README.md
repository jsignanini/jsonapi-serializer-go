# jsonapi-serializer-go

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/jsignanini/jsonapi-serializer-go)
[![Build Status](https://travis-ci.org/jsignanini/jsonapi-serializer-go.svg?branch=master)](https://travis-ci.org/jsignanini/jsonapi-serializer-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsignanini/jsonapi-serializer-go)](https://goreportcard.com/report/github.com/jsignanini/jsonapi-serializer-go)
[![Coverage Status](https://coveralls.io/repos/github/jsignanini/jsonapi-serializer-go/badge.svg?branch=master)](https://coveralls.io/github/jsignanini/jsonapi-serializer-go?branch=master)

## Installation
Install jsonapi-serializer-go with:
```sh
go get -u github.com/jsignanini/jsonapi-serializer-go
```

Then, import it using:
```go
import "github.com/jsignanini/jsonapi-serializer-go"
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/jsignanini/jsonapi-serializer-go"
)

func main() {
	// sample data
	type (
		BookBinding string
		BookSubject string
		Author      struct {
			ID        string `jsonapi:"primary,authors"`
			FirstName string `jsonapi:"attribute,first_name"`
			LastName  string `jsonapi:"attribute,last_name"`
		}
		Book struct {
			ISBN            string        `jsonapi:"primary,books"`
			Bindings        []BookBinding `jsonapi:"attribute,bindings"`
			PublicationYear int           `jsonapi:"attribute,publication_date"`
			Subject         BookSubject   `jsonapi:"attribute,subject"`
			Title           string        `jsonapi:"attribute,title"`
			Author          *Author       `jsonapi:"attribute,author"`
		}
	)
	const (
		BookBindingHardcover BookBinding = "Hardcover"
		BookBindingPaperback BookBinding = "Paperback"
	)
	cosmos := Book{
		ISBN:            "0-394-50294-9",
		Bindings:        []BookBinding{BookBindingHardcover, BookBindingPaperback},
		PublicationYear: 1980,
		Subject:         "Cosmology",
		Title:           "Cosmos",
		Author: &Author{
			ID:        "c3a6ddb6-7e5e-4264-bd03-ef6e41d76365",
			FirstName: "Carl",
			LastName:  "Sagan",
		},
	}

	// marshal
	jsonBytes, err := jsonapi.Marshal(&cosmos, nil)
	if err != nil {
		panic(err)
	}

	// print output
	fmt.Println(string(jsonBytes))
}
```


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
