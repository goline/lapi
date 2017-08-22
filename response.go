package lapi

// Response is a application's response
type Response interface {
	ResponseDescriber
	ResponseInformer
	ResponseSender
}

type ResponseInformer interface {
	// Status gets HTTP status code
	Status() int

	// WithStatus sets HTTP status code
	WithStatus(status int)

	// Message returns HTTP status message
	Message() string

	// WithMessage sets HTTP status message
	WithMessage(message string)
}

// ResponseDescriber handles content
type ResponseDescriber interface {
	// Content gets response's content
	Content() interface{}

	// WithContent sets response's content
	WithContent(content interface{})
}

type ResponseSender interface {
	// Send flushes response out
	Send() error
}
