package lapi

import "net/http"

// Responser is a application's response
type Responser interface {
	// Status sets HTTP status code
	Status(status int)

	// Message sets HTTP status message
	Message(message string)

	// Send flushes response out
	Send() error

	// SendHTTP allows to send internal HTTP Response instead
	SendHTTP(res http.Response) error
}