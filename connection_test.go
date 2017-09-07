package lapi

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	c := NewConnection(nil, nil)
	if c == nil {
		t.Errorf("Expects c is not nil")
	}
}

func TestFactoryConnection_Request(t *testing.T) {
	r := &FactoryRequest{}
	c := &FactoryConnection{}
	c.request = r
	if c.Request() == nil {
		t.Errorf("Expects request is not nil")
	}
}

func TestFactoryConnection_WithRequest(t *testing.T) {
	r := &FactoryRequest{}
	c := &FactoryConnection{}
	if c.WithRequest(r).Request() == nil {
		t.Errorf("Expects request is not nil")
	}
}

func TestFactoryConnection_Response(t *testing.T) {
	r := &FactoryResponse{}
	c := &FactoryConnection{}
	c.response = r
	if c.Response() == nil {
		t.Errorf("Expects response is not nil")
	}
}

func TestFactoryConnection_WithResponse(t *testing.T) {
	r := &FactoryResponse{}
	c := &FactoryConnection{}
	if c.WithResponse(r).Response() == nil {
		t.Errorf("Expects response is not nil")
	}
}
