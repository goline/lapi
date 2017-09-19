package lapi

import "net/http"

// Rescuer handles error
type Rescuer interface {
	// Rescue handles error, it returns nil if error is handled,
	// and error itself if could not be handled properly
	// Server should panic if an error is returned
	Rescue(connection Connection, err error) error
}

func NewRescuer() Rescuer {
	return &FactoryRescuer{}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FactoryRescuer struct{}

func (h *FactoryRescuer) Rescue(connection Connection, err error) error {
	if connection == nil {
		return err
	}
	switch e := err.(type) {
	case SystemError:
		h.handleSystemError(connection, e)
	case StackError:
		h.handleStackError(connection, e)
	default:
		h.handleUnknownError(connection, e)
	}

	return nil
}

func (h *FactoryRescuer) handleSystemError(c Connection, err SystemError) {
	switch err.Code() {
	case ERROR_HTTP_NOT_FOUND:
		c.Response().WithStatus(http.StatusNotFound).
			WithContent(&ErrorResponse{"ERROR_HTTP_NOT_FOUND", http.StatusText(http.StatusNotFound)})
	case ERROR_HTTP_BAD_REQUEST:
		c.Response().WithStatus(http.StatusBadRequest).
			WithContent(&ErrorResponse{"ERROR_HTTP_BAD_REQUEST", http.StatusText(http.StatusBadRequest)})
	default:
		c.Response().WithStatus(http.StatusInternalServerError).
			WithContent(&ErrorResponse{"ERROR_INTERNAL_SERVER_ERROR", err.Error()})
	}
}

func (h *FactoryRescuer) handleStackError(c Connection, err StackError) {
	c.Response().WithStatus(err.Status()).WithContent(&ErrorResponse{"", err.Error()})
}

func (h *FactoryRescuer) handleUnknownError(c Connection, err error) {
	if e, ok := err.(ErrorStatus); ok == true {
		c.Response().WithStatus(e.Status())
	} else {
		c.Response().WithStatus(http.StatusInternalServerError)
	}
	code := "ERROR_UNKNOWN_ERROR"
	if e, ok := err.(ErrorCoder); ok == true {
		code = e.Code()
	}
	c.Response().WithContent(&ErrorResponse{code, err.Error()})
}
