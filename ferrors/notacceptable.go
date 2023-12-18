package ferrors

import "net/http"

type ErrNotAcceptable struct {
	error
	message string
}

func NewNotAcceptable(err error) *ErrNotAcceptable {
	return &ErrNotAcceptable{
		error:   err,
		message: http.StatusText(http.StatusNotAcceptable),
	}
}

func (*ErrNotAcceptable) NotAcceptable() bool {
	return true
}
