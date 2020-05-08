package jsonapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	type Car struct {
		VIN   string `jsonapi:"primary,cars"`
		Make  string `jsonapi:"attribute,make"`
		Model string `jsonapi:"attribute,model"`
	}

	type respondTest struct {
		ExpectedBody        []byte
		ExpectedContentType string
		ExpectedStatusCode  int
		Handler             http.HandlerFunc
	}

	tests := []respondTest{
		{
			ExpectedBody: []byte(`{
	"data": {
		"id": "5YJSA1DG9DFP14705",
		"type": "cars",
		"attributes": {
			"make": "Honda",
			"model": "CR-V"
		}
	},
	"jsonapi": {
		"version": "1.0"
	}
}`),
			ExpectedContentType: ContentType,
			ExpectedStatusCode:  http.StatusOK,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				car := Car{
					VIN:   "5YJSA1DG9DFP14705",
					Make:  "Honda",
					Model: "CR-V",
				}
				if err := Respond(w, r, http.StatusOK, &car); err != nil {
					t.Errorf("expected no error, got: %s", err.Error())
				}
			}),
		},
		{
			ExpectedContentType: ContentType,
			ExpectedStatusCode:  http.StatusOK,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				car := Car{
					Make:  "Honda",
					Model: "CR-V",
				}
				if err := Respond(w, r, http.StatusOK, &car); err != nil {
					if err.Error() != "ID must be set" {
						t.Errorf("expected error: %s, got: %s", "ID must be set", err.Error())
					}
				}
			}),
		},
	}

	for _, rt := range tests {
		w := httptest.NewRecorder()
		rt.Handler(w, httptest.NewRequest("GET", "http://example.com/foo", nil))
		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != rt.ExpectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", rt.ExpectedStatusCode, res.StatusCode)
		}
		if res.Header.Get("Content-Type") != rt.ExpectedContentType {
			t.Errorf("expected content-type header: %s, got: %s", rt.ExpectedContentType, res.Header.Get("Content-Type"))
		}
		if bytes.Compare(body, rt.ExpectedBody) != 0 {
			t.Errorf("expected body: %s, got: %s", string(rt.ExpectedBody), string(body))
		}
	}

}

func TestRespondError(t *testing.T) {
	type respondErrorTest struct {
		ExpectedBody        []byte
		ExpectedContentType string
		ExpectedStatusCode  int
		Handler             http.HandlerFunc
	}

	tests := []respondErrorTest{
		{
			ExpectedBody: []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"title": "not_found"
		}
	]
}`),
			ExpectedContentType: ContentType,
			ExpectedStatusCode:  http.StatusNotFound,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := RespondError(w, r, http.StatusNotFound, nil, Error{Title: "not_found"}); err != nil {
					t.Errorf("expected no error, got: %s", err.Error())
				}
			}),
		},
		{
			ExpectedBody: []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"title": "internal_server_error"
		}
	]
}`),
			ExpectedContentType: ContentType,
			ExpectedStatusCode:  http.StatusInternalServerError,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := RespondError(w, r, http.StatusInternalServerError, nil, Error{Title: "internal_server_error"}); err != nil {
					t.Errorf("expected no error, got: %s", err.Error())
				}
			}),
		},
		{
			ExpectedBody: []byte(`{
	"jsonapi": {
		"version": "1.0"
	},
	"errors": [
		{
			"title": "internal_server_error"
		}
	]
}`),
			ExpectedContentType: ContentType,
			ExpectedStatusCode:  http.StatusInternalServerError,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := RespondError(w, r, http.StatusInternalServerError, nil, Error{Title: "internal_server_error"}); err != nil {
					t.Errorf("expected no error, got: %s", err.Error())
				}
			}),
		},
	}

	for _, rt := range tests {
		w := httptest.NewRecorder()
		rt.Handler(w, httptest.NewRequest("GET", "http://example.com/foo", nil))
		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != rt.ExpectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", rt.ExpectedStatusCode, res.StatusCode)
		}
		if res.Header.Get("Content-Type") != rt.ExpectedContentType {
			t.Errorf("expected content-type header: %s, got: %s", rt.ExpectedContentType, res.Header.Get("Content-Type"))
		}
		if bytes.Compare(body, rt.ExpectedBody) != 0 {
			t.Errorf("expected body: %s, got: %s", string(rt.ExpectedBody), string(body))
		}
	}

}
