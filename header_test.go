package lapi

import (
	"testing"
)

var contentType = []string{"application/json", "charset=utf8"}

func TestNewHeader(t *testing.T) {
	h := NewHeader()
	if _, ok := h.(Header); ok == false {
		t.Errorf("Expects an instance of Header. Got %+v", h)
	}
}

func TestFactoryHeader_Get(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.items["Content-Type"] = contentType
	values, ok := h.Get("content-Type")
	if ok == false {
		t.Errorf("Expects content-Type key to be existed")
	}
	if values[0] != "application/json" {
		t.Errorf("Expects values[0] to be application/json")
	}
	if values[1] != "charset=utf8" {
		t.Errorf("Expects values[0] to be application/json")
	}
}

func TestFactoryHeader_Has(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.items["Content-Type"] = contentType
	if !h.Has("content-Type") || !h.Has("content-type") || h.Has("ContentType") {
		t.Errorf("Expects has to return correct")
	}
}

func TestFactoryHeader_Set(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.Set("content-type", contentType...)
	if h.items["Content-Type"][0] != "application/json" || h.items["Content-Type"][1] != "charset=utf8" {
		t.Errorf("Expects Set to be able to set key-value")
	}
}

func TestFactoryHeader_Remove(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.items["Content-Type"] = contentType
	h.Remove("content-TYPE")
	if len(h.items) > 0 {
		t.Errorf("Expects content-TYPE is removed")
	}
}

func TestFactoryHeader_All(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.items["Content-Type"] = contentType
	h.items["Content-Length"] = []string{"1234"}
	if len(h.All()) != 2 {
		t.Errorf("Expects to get 2 items")
	}
}

func TestFactoryHeader_Line(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.items["Content-Type"] = contentType
	if key, line := h.Line("content-TYPE"); key != "Content-Type" || line != "application/json, charset=utf8" {
		t.Errorf("Expects Line to be correct")
	}
}

func TestFactoryHeader_Lines(t *testing.T) {
	h := &FactoryHeader{make(map[string][]string)}
	h.items["Content-Type"] = contentType
	h.items["Content-Length"] = []string{"1234"}

	lines := h.Lines()
	if lines["Content-Type"] != "application/json, charset=utf8" || lines["Content-Length"] != "1234" {
		t.Errorf("Expects Lines to be correct")
	}
}
