package lapi

import "github.com/goline/errors"

// Handler is a request's handler
type Handler interface {
	// Handle performs logic for solving request
	Handle(connection Connection) (interface{}, errors.Error)
}

// IOHandler describes input and output for handler
// This interface aims to support to generate documentation only
type IOHandler interface {
	IO() (input interface{}, output interface{})
}
