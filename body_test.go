package lapi

import (
	"github.com/goline/lapi/parser"
	"testing"
)

func TestNewBody(t *testing.T) {
	b := NewBody()
	if b == nil {
		t.Errorf("Expects b is not nil")
	}
}

func TestFactoryBody_Charset(t *testing.T) {
	b := &FactoryBody{}
	b.charset = "UTF-8"
	if b.Charset() != "UTF-8" {
		t.Errorf("Expects charset is UTF-8. Got %s", b.Charset())
	}
}

func TestFactoryBody_WithCharset(t *testing.T) {
	b := &FactoryBody{}
	b.WithCharset("UTF-8")
	if b.charset != "UTF-8" {
		t.Errorf("Expects charset is UTF-8. Got %s", b.charset)
	}
}

func TestFactoryBody_Content(t *testing.T) {
	b := &FactoryBody{}
	if b.content != nil || b.err != nil {
		t.Errorf("Expects content and err is nil. Got %v, %v", b.content, b.err)
	}

	b.content = "a_string"
	if c, err := b.Content(); err != nil || c != "a_string" {
		t.Errorf("Expects content and err are correct. Got %v, %v", c, err)
	}
}

type sampleItem struct {
	Price float64 `json:"price"`
}

func TestFactoryBody_WithContent(t *testing.T) {
	b := &FactoryBody{ParserManager: NewParserManager()}
	b.WithContent("a_string")
	if b.err == nil {
		t.Errorf("Expects err is not nil")
	}
	if e, ok := b.err.(SystemError); ok == false || e.Code() != ERROR_NO_PARSER_FOUND {
		t.Errorf("Expects err is SystemError. Got %v", b.err)
	}

	b.WithParser(&parser.JsonParser{})
	b.WithContentType(CONTENT_TYPE_JSON)
	v := &sampleItem{Price: 5.67}
	b.WithContent(v)
	if b.err != nil {
		t.Errorf("Expects err is nil")
	}
	if string(b.contentBytes) != "{\"price\":5.67}" {
		t.Errorf("Expects content is set correctly. Got %s", string(b.contentBytes))
	}
}

func TestFactoryBody_ContentBytes(t *testing.T) {
	b := &FactoryBody{}
	if len(b.contentBytes) != 0 || b.err != nil {
		t.Errorf("Expects contentBytes is empty and err is nil. Got %v, %v", b.contentBytes, b.err)
	}

	b.contentBytes = []byte("a_string")
	if bb, err := b.ContentBytes(); len(bb) == 0 || err != nil {
		t.Errorf("Expects bb is not empty and err is nil. Got %v, %v", bb, err)
	}
}

func TestFactoryBody_WithContentBytes(t *testing.T) {
	b := &FactoryBody{ParserManager: NewParserManager()}
	b.WithContentBytes([]byte("a_string"), nil)
	if b.err == nil {
		t.Errorf("Expects err is not nil")
	}
	if e, ok := b.err.(SystemError); ok == false || e.Code() != ERROR_NO_PARSER_FOUND {
		t.Errorf("Expects err is SystemError. Got %v", b.err)
	}

	v := &sampleItem{}
	b.WithParser(&parser.JsonParser{})
	b.WithContentType(CONTENT_TYPE_JSON)
	b.WithContentBytes([]byte("{\"price\":5.67}"), v)
	if b.err != nil {
		t.Errorf("Expects err is nil. Got %v", b.err)
	}
	if v.Price != 5.67 {
		t.Errorf("Expects content is set correctly. Got %v", v)
	}
}
