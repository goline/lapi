package lapi

import (
	"testing"
)

func TestNewHeader(t *testing.T) {
	h := NewHeader()
	if _, ok := h.(Header); ok == false {
		t.Errorf("Expects an instance of Header. Got %+v", h)
	}
}

func TestFactoryHeader_Get(t *testing.T) {
	h := &FactoryHeader{make(map[string]string)}
	h.items["Content-Type"] = "application/json"
	values, ok := h.Get("content-Type")
	if ok == false {
		t.Errorf("Expects content-Type key to be existed")
	}
	if values != "application/json" {
		t.Errorf("Expects values[0] to be application/json")
	}
}

func TestFactoryHeader_Has(t *testing.T) {
	h := &FactoryHeader{make(map[string]string)}
	h.items["Content-Type"] = "application/json"
	if !h.Has("content-Type") || !h.Has("content-type") || h.Has("ContentType") {
		t.Errorf("Expects has to return correct")
	}
}

func TestFactoryHeader_Set(t *testing.T) {
	h := &FactoryHeader{make(map[string]string)}
	h.Set("content-type", "application/json")
	if h.items["Content-Type"] != "application/json" {
		t.Errorf("Expects Set to be able to set key-value")
	}
}

func TestFactoryHeader_Remove(t *testing.T) {
	h := &FactoryHeader{make(map[string]string)}
	h.items["Content-Type"] = "application/json"
	h.Remove("content-TYPE")
	if len(h.items) > 0 {
		t.Errorf("Expects content-TYPE is removed")
	}
}

func TestFactoryHeader_All(t *testing.T) {
	h := &FactoryHeader{make(map[string]string)}
	h.items["Content-Type"] = "application/json"
	h.items["Content-Length"] = "1234"
	if len(h.All()) != 2 {
		t.Errorf("Expects to get 2 items")
	}
}
