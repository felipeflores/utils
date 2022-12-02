package ferrors

import (
	"net/http"
)

type ErrForbidden struct {
	error
	message string
}

func NewForbidden(err error) *ErrForbidden {
	return &ErrForbidden{
		error:   err,
		message: http.StatusText(http.StatusForbidden),
	}
}

func (*ErrForbidden) Forbidden() bool {
	return true
}
