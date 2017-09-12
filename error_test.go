package lapi

import (
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError("err_code", "err_msg", nil)
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
}
