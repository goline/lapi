package lapi

import (
	"github.com/goline/errors"
	"net/http"
	"testing"
)

func TestNewRescuer(t *testing.T) {
	h := NewRescuer()
	if h == nil {
		t.Errorf("Expects h is not nil")
	}
}

func TestFactoryRescuer_Rescue_HandleSystemError_HttpNotFound(t *testing.T) {
	c := NewConnection(nil, getEmptyResponse())
	e := errors.New(ERR_HTTP_NOT_FOUND, "")
	h := &FactoryRescuer{}
	h.Rescue(c, e)
	if c.Response().Status() != http.StatusNotFound {
		t.Errorf("Expects http status is StatusNotFound. Got %d", c.Response().Status())
	}
}

func TestFactoryRescuer_Rescue_HandleSystemError_HttpBadRequest(t *testing.T) {
	c := NewConnection(nil, getEmptyResponse())
	e := errors.New(ERR_HTTP_BAD_REQUEST, "")
	h := &FactoryRescuer{}
	h.Rescue(c, e)
	if c.Response().Status() != http.StatusBadRequest {
		t.Errorf("Expects http status is StatusBadRequest. Got %d", c.Response().Status())
	}
}

func TestFactoryRescuer_Rescue_UnknownError(t *testing.T) {
	c := NewConnection(nil, getEmptyResponse())
	e := errors.New("11", "err1")
	h := &FactoryRescuer{}
	h.Rescue(c, e)
	if c.Response().Status() != http.StatusInternalServerError {
		t.Errorf("Expects http status is StatusInternalServerError. Got %d", c.Response().Status())
	}
}

type myUnknownError struct{}

func (e *myUnknownError) Error() string { return "" }
func (e *myUnknownError) Status() int   { return http.StatusInternalServerError }

func TestFactoryRescuer_Rescue_UnknownErrorWithStatus(t *testing.T) {
	t.SkipNow() // temporary not support error's http status for now
	c := NewConnection(nil, getEmptyResponse())
	e := &myUnknownError{}
	h := &FactoryRescuer{}
	h.Rescue(c, e)
	if c.Response().Status() != http.StatusInternalServerError {
		t.Errorf("Expects http status is StatusInternalServerError. Got %d", c.Response().Status())
	}
}

func TestFactoryRescuer_Rescue_NotHandle(t *testing.T) {
	e := &myUnknownError{}
	h := &FactoryRescuer{}
	err := h.Rescue(nil, e)
	if err == nil {
		t.Errorf("Expects err is nil")
	}
}

func getEmptyResponse() Response {
	return &FactoryResponse{Body: NewBody()}
}
