package lapi

// Error represents for a common error
type Error interface {
	ErrorCoder
	ErrorMessager
	ErrorTracer
}

// HttpError is an Error with HTTP status
type HttpError interface {
	ErrorStatus
	Error

	error // implements error
}

// StackError contains multiple Errors
type StackError interface {
	ErrorStatus
	Errors() []Error

	error // implements error
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