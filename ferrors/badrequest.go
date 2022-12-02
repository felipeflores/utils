package ferrors

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrBadRequest struct {
	error
	message string
	fields  map[string]error
}

func NewBadRequest(err error) *ErrBadRequest {
	var fields map[string]error
	if e, ok := err.(validation.Errors); ok {
		fields = e
	}
	return &ErrBadRequest{
		error:   err,
		message: http.StatusText(http.StatusBadRequest),
		fields:  fields,
	}
}

func (*ErrBadRequest) BadRequest() bool {
	return true
}

func (e *ErrBadRequest) GetFields() map[string]error {
	return e.fields
}
