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

type errorStackResponse struct {
	Errors []errorItemResponse `json:"errors"`
}

type errorItemResponse struct {
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
			WithContent(h.getResponseContentForErrors(NewError("ERROR_HTTP_NOT_FOUND", http.StatusText(http.StatusNotFound), nil)))
	case ERROR_HTTP_BAD_REQUEST:
		c.Response().WithStatus(http.StatusBadRequest).
			WithContent(h.getResponseContentForErrors(NewError("ERROR_HTTP_BAD_REQUEST", http.StatusText(http.StatusBadRequest), nil)))
	default:
		c.Response().WithStatus(http.StatusInternalServerError).
			WithContent(h.getResponseContentForErrors(NewError("", err.Error(), nil)))
	}
}

func (h *FactoryRescuer) handleStackError(c Connection, err StackError) {
	c.Response().WithStatus(err.Status()).WithContent(h.getResponseContentForErrors(err.Errors()...))
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
	c.Response().WithContent(h.getResponseContentForErrors(NewError(code, err.Error(), err)))
}

func (h *FactoryRescuer) getResponseContentForErrors(errors ...Error) *errorStackResponse {
	ei := make([]errorItemResponse, len(errors))
	for i, e := range errors {
		ei[i] = errorItemResponse{e.Code(), e.Message()}
	}
	return &errorStackResponse{ei}
}
