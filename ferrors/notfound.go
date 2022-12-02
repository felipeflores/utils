package ferrors

import "net/http"

type ErrNotFound struct {
	error
	message string
}

func NewNotFound(err error) *ErrNotFound {
	return &ErrNotFound{
		error:   err,
		message: http.StatusText(http.StatusNotFound),
	}
}

func (*ErrNotFound) NotFound() bool {
	return true
}
