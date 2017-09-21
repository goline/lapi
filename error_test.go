package lapi

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError("err_code", "err_msg", nil)
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
}

func TestFactoryError_Error(t *testing.T) {
	e := &FactoryError{
		code:    ERR_HTTP_NOT_FOUND,
		message: "A message",
		err:     errors.New("original error"),
	}
	if e.Error() != "original error" {
		t.Errorf("Expects Error is 'original error'. Got %s", e.Error())
	}
}
