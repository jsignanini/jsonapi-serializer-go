package jsonapi

import (
	"net/http"
)

// Respond
func Respond(w http.ResponseWriter, r *http.Request, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", ContentType)
	body, err := Marshal(v, nil)
	if err != nil {
		return err
	}
	return respond(w, r, statusCode, body)
}

// RespondError
func RespondError(w http.ResponseWriter, r *http.Request, statusCode int, p *MarshalParams, errs ...Error) error {
	w.Header().Set("Content-Type", ContentType)
	// TODO figure out how to trigger this error for test coverage
	body, _ := MarshalErrors(p, errs...)
	return respond(w, r, statusCode, body)
}

func respond(w http.ResponseWriter, r *http.Request, statusCode int, body []byte) (err error) {
	w.WriteHeader(statusCode)
	// TODO figure out how to trigger this error for test coverage
	_, err = w.Write(body)
	return err
}
