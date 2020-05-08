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
	body, err := MarshalErrors(p, errs...)
	if err != nil {
		return err
	}
	return respond(w, r, statusCode, body)
}

func respond(w http.ResponseWriter, r *http.Request, statusCode int, body []byte) (err error) {
	w.WriteHeader(statusCode)
	_, err = w.Write(body)
	return err
}
