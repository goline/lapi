package lapi

import (
	"errors"
	"fmt"
	"net/http"
)

// Handler is a request's handler
type Handler interface {
	// Handle performs logic for solving request
	Handle(connection Connection) (interface{}, error)
}

// ErrorHandler handles error
type ErrorHandler interface {
	// HandleError handles error, it returns nil if error is handled,
	// and error itself if could not be handled properly
	// Server should panic if an error is returned
	HandleError(connection Connection, err error) error
}

func NewErrorHandler() ErrorHandler {
	return &FactoryErrorHandler{}
}

type FactoryErrorHandler struct{}

func (h *FactoryErrorHandler) HandleError(connection Connection, err error) error {
	if e, ok := err.(ErrorStatus); ok == true {
		connection.Response().WithStatus(e.Status())
	} else {
		connection.Response().WithStatus(http.StatusInternalServerError)
	}

	var es []Error
	switch e := err.(type) {
	case Error:
		es = make([]Error, 1)
		es[0] = e
	case SystemError:
		es = make([]Error, 1)
		switch e.Code() {
		case ERROR_HTTP_NOT_FOUND:
			es[0] = NewError("ERROR_HTTP_NOT_FOUND", e.Message(), nil)
		default:
			es[0] = NewError(fmt.Sprintf("%d", e.Code()), e.Message(), nil)
		}
	case StackError:
		es = e.Errors()
	default:
		es[0] = NewError("", "ERROR_HANDLE_INVALID_ERROR", errors.New("Error's type is not supported."))
	}

	ei := make([]errorItemResponse, len(es))
	for i, e := range es {
		ei[i] = errorItemResponse{e.Code(), e.Message()}
	}
	er := &errorStackResponse{ei}
	connection.Response().WithContent(er)
	return nil
}

type errorStackResponse struct {
	Errors []errorItemResponse `json:"errors"`
}

type errorItemResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
