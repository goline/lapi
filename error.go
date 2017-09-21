package lapi

import (
	"fmt"
)

// Error represents for a common error
type Error interface {
	ErrorCoder
	ErrorMessager
	ErrorTracer
	error
}

type ErrorCoder interface {
	// Code returns error's code
	Code() string
}

type ErrorMessager interface {
	// Message returns error's message
	Message() string
}

type ErrorTracer interface {
	// Trace returns original system error
	Trace() error
}

type ErrorStatuser interface {
	// Status returns http status code
	Status() int
}

func NewError(code string, message string, err error) Error {
	if err == nil {
		err = fmt.Errorf("[%s] %s", code, message)
	}
	return &FactoryError{code, message, err}
}

type FactoryError struct {
	code    string
	message string
	err     error
}

func (e *FactoryError) Code() string {
	return e.code
}

func (e *FactoryError) Message() string {
	return e.message
}

func (e *FactoryError) Trace() error {
	return e.err
}

func (e *FactoryError) Error() string {
	return e.err.Error()
}
