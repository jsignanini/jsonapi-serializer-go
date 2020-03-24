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
