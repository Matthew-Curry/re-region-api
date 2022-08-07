package apperrors

import (
	"fmt"
)

/* Implements a custom error struct with public constructors for all app errors. Errors are categorized
into "ErrorKinds" that callers at higher layers use to determine how to handle the error. */

type ErrorKind int

type AppError struct {
	message string
	kind ErrorKind
	source error
}

func (e *AppError) Error() string {
	if e.source != nil {
		return fmt.Sprintf("%s:%s", e.message, e.source.Error())
	}
	
	return fmt.Sprintf(e.message)
}

func (e *AppError) Message() string {
	return e.message
}

func Test() error {
	return &AppError{}
  }

func (e *AppError) IsKind(kind ErrorKind) bool {
	if e.kind == kind {
		return true
	} else {
		return false
	}
}
