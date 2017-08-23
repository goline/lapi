package lapi

import (
	"fmt"
	"strings"
)

// Error represents for a common error
type Error interface {
	ErrorCoder
	ErrorMessager
	ErrorTracer
	error
}

// HttpError is an Error with HTTP status
type HttpError interface {
	ErrorStatus
	Error
}

// StackError contains multiple Errors
type StackError interface {
	ErrorStatus
	Errors() []Error
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

type ErrorStatus interface {
	Status() int
}

func NewError(code string, message string, err error) Error {
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
	return fmt.Sprint(e.code, e.message)
}

func NewHttpError(status int, error Error) HttpError {
	return &FactoryHttpError{status, error}
}

type FactoryHttpError struct {
	status int
	Error
}

func (e *FactoryHttpError) Status() int {
	return e.status
}

func NewStackError(status int, errors ...Error) StackError {
	return &FactoryStackError{status, errors}
}

type FactoryStackError struct {
	status int
	errors []Error
}

func (e *FactoryStackError) Status() int {
	return e.status
}

func (e *FactoryStackError) Errors() []Error {
	return e.errors
}

func (e *FactoryStackError) Error() string {
	total := len(e.errors)
	if total == 0 {
		return ""
	}
	messages := make([]string, total)
	for i, err := range e.errors {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "\n")
}
