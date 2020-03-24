package jsonapi

import (
	"testing"
)

func TestError(t *testing.T) {
	sampleError := Error{
		ID:     "SAMPLE_ERROR",
		Status: "500",
	}
	errorString := sampleError.Error()
	expected := `{
	"id": "SAMPLE_ERROR",
	"status": "500"
}`
	if errorString != expected {
		t.Errorf("expected error:\n%s\ngot error:\n%s\n", expected, errorString)
	}
}
