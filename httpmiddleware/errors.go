package httpmiddleware

import (
	"fmt"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type (
	badrequest interface {
		BadRequest() bool
	}
)

func httpStatusCode(err error) int {
	fmt.Println(err)
	switch err.(type) {
	case badrequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
