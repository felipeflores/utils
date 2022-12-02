package ferrors

import "net/http"

type ErrInternalServer struct {
	error
	message string
}

func NewInternalServer(err error) *ErrInternalServer {
	return &ErrInternalServer{
		error:   err,
		message: http.StatusText(http.StatusInternalServerError),
	}
}

func (*ErrInternalServer) InternalServer() bool {
	return true
}
