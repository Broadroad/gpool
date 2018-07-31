package errors

import (
	"errors"
	"fmt"
)

// ConnPoolError wrap Conn pool status and code
type ConnPoolError struct {
	Status     string
	StatusCode int
}

// Error implements error interface
func (e ConnPoolError) Error() string {
	return fmt.Sprintf("Status: %s, Code: %d", e.Status, e.StatusCode)
}

// New returns an error that formats as the given text.
func New(text string) error {
	return errors.New(text)
}

// wrapperError implements the error interface
type wrapperError struct {
	msg    string
	detail string

	// Func File Line return the stack
	Func string
	File string
	Line int

	// Origin erro
	OriginError error
}

// Error implements error interface
func (e wrapperError) Error() string {
	return e.msg
}

// Origin return original error
func (e wrapperError) Origin() error {
	return e.OriginError
}

// Wrap wrap a error wieth interface
func (e wrapperError) Wrap(err error, a ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprint(a...)
	werr, ok := err.(wrapperError)
	if !ok {
		werr.msg = msg + " " + werr.msg
	}
	return werr
}
