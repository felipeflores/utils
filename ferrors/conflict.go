package ferrors

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrConflict struct {
	error
	message string
	fields  map[string]error
}

func NewConflict(err error) *ErrConflict {
	var fields map[string]error
	if e, ok := err.(validation.Errors); ok {
		fields = e
	}
	return &ErrConflict{
		error:   err,
		message: http.StatusText(http.StatusConflict),
		fields:  fields,
	}
}

func (*ErrConflict) Conflict() bool {
	return true
}

func (e *ErrConflict) GetFields() map[string]error {
	return e.fields
}
