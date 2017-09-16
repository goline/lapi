package lapi

import (
	"errors"
	"testing"
)

func TestPanicOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expects r is not nil")
		}
	}()
	PanicOnError(errors.New("ERROR"))
}

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expects r is not nil")
		}
	}()
	Must(nil, errors.New("ERROR"))
}
