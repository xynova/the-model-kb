package errors

import (
	"errors"
	"fmt"
)

// Code is a stable machine-readable error class.
type Code string

const (
	CodeInvalidArgument Code = "invalid_argument"
	CodeNotFound        Code = "not_found"
	CodeUnavailable     Code = "unavailable"
	CodeFailed          Code = "failed"
	CodeParse           Code = "parse"
)

// Error is a typed domain error with optional fields and unwrap support.
type Error struct {
	Code    Code
	Op      string
	Message string
	Fields  map[string]any
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause == nil {
		return fmt.Sprintf("%s: %s", e.Op, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Cause)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// With attaches a field for operators and callers.
func (e *Error) With(key string, value any) *Error {
	if e.Fields == nil {
		e.Fields = map[string]any{}
	}
	e.Fields[key] = value
	return e
}

// Wrap builds a typed error around a cause.
func Wrap(cause error, code Code, op, message string) *Error {
	return &Error{Code: code, Op: op, Message: message, Cause: cause}
}

// New builds a typed error without a cause.
func New(code Code, op, message string) *Error {
	return &Error{Code: code, Op: op, Message: message}
}

// IsCode reports whether err (or any wrapped cause) carries code.
func IsCode(err error, code Code) bool {
	var de *Error
	if errors.As(err, &de) {
		return de.Code == code
	}
	return false
}
