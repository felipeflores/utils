package ferrors

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// ErrUnauthorized implements middleware's unauthorized and message interfaces.
type ErrUnauthorized struct {
	error

	message string
}

// NewUnauthorized returns a struct that implements go's error interface,
// and middleware's unauthorized and message interfaces,
// setting message to err.Error() string.
// If err is nil a panic() will raise.
// To have the stack trace, call this function with
// errors.Wrap(), errors.Wrapf() or errors.WithStack().
// Example: NewUnauthorized(errors.WithStack(err))
func NewUnauthorized(err error) *ErrUnauthorized {
	return &ErrUnauthorized{message: http.StatusText(http.StatusUnauthorized), error: err}
}

// Unauthorized implements middleware's unauthorized interface.
// The return value of this function is irrelevant.
func (*ErrUnauthorized) Unauthorized() bool { return true }

// Msg implements middleware's message interface.
func (e *ErrUnauthorized) Msg() string { return e.message }

// WithMessage returns an unauthorized error after setting a custom message.
func (e *ErrUnauthorized) WithMessage(message string) *ErrUnauthorized {
	e.message = message
	return e
}

// StatusOnly returns an unauthorized error after removing the message.
func (e *ErrUnauthorized) StatusOnly() *ErrUnauthorized {
	e.message = ""
	return e
}

// Cause implements the pkg/errors interface.
func (e *ErrUnauthorized) Cause() error {
	return e.error
}

// Format implements the fmt.Formatter.
// With this method, it's possible to have the stack trace.
// The stack trace is provided by the pkg/errors.
func (e *ErrUnauthorized) Format(st fmt.State, verb rune) {
	if _, ok := e.error.(fmt.Formatter); ok {
		e.error.(fmt.Formatter).Format(st, verb)
		return
	}

	switch verb {
	case 'v':
		if st.Flag('+') {
			for _, f := range e.StackTrace() {
				f.Format(st, verb)
			}
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(st, "%q", e.Error())
	case 'q':
		_, _ = io.WriteString(st, e.Error())
	}
}

// StackTrace implements the pkg/errors stack trace interface.
// With this method, it's possible to have the stack trace.
func (e *ErrUnauthorized) StackTrace() errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	if _, ok := e.error.(stackTracer); ok {
		return e.error.(stackTracer).StackTrace()
	}

	return nil
}
