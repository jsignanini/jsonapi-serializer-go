package jsonapi

import (
	"net/http"
)

// Respond encodes v in to a JSON:API object and writes it to the body of response w. It also sets
// statusCode as the response status code.
func Respond(w http.ResponseWriter, r *http.Request, statusCode int, v interface{}) error {
	body, err := Marshal(v, nil)
	if err != nil {
		return err
	}
	return respond(w, r, statusCode, body)
}

// RespondError encodes v in to a JSON:API error object and writes it to the body of response w. It
// also sets statusCode as the response status code.
func RespondError(w http.ResponseWriter, r *http.Request, statusCode int, p *MarshalParams, errs ...Error) error {
	// TODO figure out how to trigger this error for test coverage
	body, _ := MarshalErrors(p, errs...)
	return respond(w, r, statusCode, body)
}

func respond(w http.ResponseWriter, r *http.Request, statusCode int, body []byte) (err error) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(statusCode)
	// TODO figure out how to trigger this error for test coverage
	_, err = w.Write(body)
	return err
}
