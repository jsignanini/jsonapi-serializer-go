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

Outputs:
```json
{
	"data": {
		"id": "0-394-50294-9",
		"type": "books",
		"attributes": {
			"author": {
				"ID": "c3a6ddb6-7e5e-4264-bd03-ef6e41d76365",
				"FirstName": "Carl",
				"LastName": "Sagan"
			},
			"bindings": [
				"Hardcover",
				"Paperback"
			],
			"publication_date": 1980,
			"subject": "Cosmology",
			"title": "Cosmos"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}
```


## TODOs

- Optionally validate jsonapi spec
- Optionally set jsonapi settings (e.g.: spec version, error/warning on document validation, etc.)
- Support omitempty tag `jsonapi:"attribute,name,omitempty"`
- Standardize internal errors
- Support non-string resource IDs
- Show error or warning when parsing an unsupported builtin type (e.g.: `complex128`)
