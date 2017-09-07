package lapi

import (
	"net/http"
	"testing"
)

func TestNewErrorHandler(t *testing.T) {
	h := NewErrorHandler()
	if h == nil {
		t.Errorf("Expects h is not nil")
	}
}

func TestFactoryErrorHandler_HandleError_HandleSystemError_HttpNotFound(t *testing.T) {
	c := NewConnection(nil, &FactoryResponse{})
	e := NewSystemError(ERROR_HTTP_NOT_FOUND, "")
	h := &FactoryErrorHandler{}
	h.HandleError(c, e)
	if c.Response().Status() != http.StatusNotFound {
		t.Errorf("Expects http status is StatusNotFound. Got %d", c.Response().Status())
	}
}

func TestFactoryErrorHandler_HandleError_HandleSystemError_HttpBadRequest(t *testing.T) {
	c := NewConnection(nil, &FactoryResponse{})
	e := NewSystemError(ERROR_HTTP_BAD_REQUEST, "")
	h := &FactoryErrorHandler{}
	h.HandleError(c, e)
	if c.Response().Status() != http.StatusBadRequest {
		t.Errorf("Expects http status is StatusBadRequest. Got %d", c.Response().Status())
	}
}

func TestFactoryErrorHandler_HandleError_HandleSystemError_Default(t *testing.T) {
	c := NewConnection(nil, &FactoryResponse{})
	e := NewSystemError(9999, "")
	h := &FactoryErrorHandler{}
	h.HandleError(c, e)
	if c.Response().Status() != http.StatusInternalServerError {
		t.Errorf("Expects http status is StatusInternalServerError. Got %d", c.Response().Status())
	}
}

func TestFactoryErrorHandler_HandleError_HandleStackError(t *testing.T) {
	c := NewConnection(nil, &FactoryResponse{})
	e := NewStackError(http.StatusInternalServerError, []Error{NewError("11", "err1", nil)})
	h := &FactoryErrorHandler{}
	h.HandleError(c, e)
	if c.Response().Status() != http.StatusInternalServerError {
		t.Errorf("Expects http status is StatusInternalServerError. Got %d", c.Response().Status())
	}
}

func TestFactoryErrorHandler_HandleError_UnknownError(t *testing.T) {
	c := NewConnection(nil, &FactoryResponse{})
	e := NewError("11", "err1", nil)
	h := &FactoryErrorHandler{}
	h.HandleError(c, e)
	if c.Response().Status() != http.StatusInternalServerError {
		t.Errorf("Expects http status is StatusInternalServerError. Got %d", c.Response().Status())
	}
}

type myUnknownError struct{}

func (e *myUnknownError) Error() string { return "" }
func (e *myUnknownError) Status() int   { return http.StatusInternalServerError }

func TestFactoryErrorHandler_HandleError_UnknownErrorWithStatus(t *testing.T) {
	c := NewConnection(nil, &FactoryResponse{})
	e := &myUnknownError{}
	h := &FactoryErrorHandler{}
	h.HandleError(c, e)
	if c.Response().Status() != http.StatusInternalServerError {
		t.Errorf("Expects http status is StatusInternalServerError. Got %d", c.Response().Status())
	}
}

func TestFactoryErrorHandler_HandleError_NotHandle(t *testing.T) {
	e := &myUnknownError{}
	h := &FactoryErrorHandler{}
	err := h.HandleError(nil, e)
	if err == nil {
		t.Errorf("Expects err is nil")
	}
}
