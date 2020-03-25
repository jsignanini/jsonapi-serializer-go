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
	expected := "{ID:SAMPLE_ERROR Links:map[] Status:500 Code: Title: Detail: Source:map[] Meta:map[]}"
	if errorString != expected {
		t.Errorf("expected error:\n%s\ngot error:\n%s\n", expected, errorString)
	}
}
