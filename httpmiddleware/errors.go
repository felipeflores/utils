package httpmiddleware

import (
	"fmt"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Fields    []Field   `json:"fields,omitempty"`
}

type Field struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type (
	badrequest interface {
		BadRequest() bool
		GetFields() map[string]error
	}
	notfound interface {
		NotFound() bool
	}
	unauthorized interface {
		Unauthorized() bool
	}
	notacceptable interface {
		NotAcceptable() bool
	}
)

func httpStatusCode(err error) int {
	switch e := err.(type) {
	case badrequest:
		fmt.Println(e.GetFields())
		return http.StatusBadRequest
	case notfound:
		return http.StatusNotFound
	case unauthorized:
		return http.StatusUnauthorized
	case notacceptable:
		return http.StatusNotAcceptable
	default:
		return http.StatusInternalServerError
	}
}
