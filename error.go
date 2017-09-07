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

// SystemError uses for system error
type SystemError interface {
	ErrorNo
	ErrorMessager
	error
}

// HttpError is an Error with HTTP status
type HttpError interface {
	ErrorStatus
	ErrorCoder
	ErrorMessager
	ErrorTracer
	error
}

// StackError contains multiple Errors
type StackError interface {
	ErrorStatus
	Errors() []Error
	error
}

// ErrorNo contains error's code
type ErrorNo interface {
	Code() uint
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
	return getErrorString(e)
}

func NewHttpError(status int, code string, message string, err error) HttpError {
	return &FactoryHttpError{status, code, message, err}
}

type FactoryHttpError struct {
	status  int
	code    string
	message string
	err     error
}

func (e *FactoryHttpError) Status() int {
	return e.status
}

func (e *FactoryHttpError) Code() string {
	return e.code
}

func (e *FactoryHttpError) Message() string {
	return e.message
}

func (e *FactoryHttpError) Trace() error {
	return e.err
}

func (e *FactoryHttpError) Error() string {
	return getErrorString(e)
}

func NewStackError(status int, errors []Error) StackError {
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

func NewSystemError(code uint, message string) SystemError {
	return &FactorySystemError{code, message}
}

type FactorySystemError struct {
	code    uint
	message string
}

func (e *FactorySystemError) Code() uint {
	return e.code
}

func (e *FactorySystemError) Message() string {
	return e.message
}

func (e *FactorySystemError) Error() string {
	return getErrorString(e)
}

func getErrorString(err interface{}) string {
	var code, message string
	if e, ok := err.(ErrorCoder); ok == true {
		code = e.Code()
	} else if e, ok := err.(ErrorNo); ok == true {
		code = fmt.Sprintf("%v", e.Code())
	}
	if e, ok := err.(ErrorMessager); ok == true {
		message = e.Message()
	}

	return fmt.Sprintf("[%v] %v", code, message)
}
