package lapi

import (
	"encoding/json"
	"github.com/goline/lapi/parser"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	r, _ := NewResponse(nil)
	if _, ok := r.(Response); ok == false {
		t.Errorf("Expects an instance of Response. Got %+v", r)
	}
}

func TestFactoryResponse_DefaultStatus(t *testing.T) {
	r, _ := NewResponse(nil)
	if r.Status() != http.StatusOK {
		t.Errorf("Expects default status code must be http.StatusOk. Got %v", r.Status())
	}
}

func TestFactoryResponse_Status(t *testing.T) {
	r := &FactoryResponse{}
	r.status = http.StatusBadRequest
	if r.Status() != http.StatusBadRequest {
		t.Errorf("Expects status code must be http.StatusBadRequest. Got %v", r.Status())
	}
}

func TestFactoryResponse_WithStatus(t *testing.T) {
	r := &FactoryResponse{}
	r.WithStatus(http.StatusBadRequest)
	if r.status != http.StatusBadRequest {
		t.Errorf("Expects status code must be http.StatusBadRequest. Got %v", r.Status())
	}
}

func TestFactoryResponse_Message(t *testing.T) {
	r := &FactoryResponse{}
	r.message = "my_own_message"
	if r.Message() != "my_own_message" {
		t.Errorf("Expects status message must be my_own_message. Got %v", r.Message())
	}
}

func TestFactoryResponse_WithMessage(t *testing.T) {
	r := &FactoryResponse{}
	r.WithMessage("my_own_message")
	if r.message != "my_own_message" {
		t.Errorf("Expects status message must be my_own_message. Got %v", r.Message())
	}
}

func TestFactoryResponse_Content(t *testing.T) {
	r := &FactoryResponse{}
	c := make(map[string]interface{})
	c["my_key"] = "my_value"
	r.content = c
	if reflect.DeepEqual(r.Content(), c) == false {
		t.Errorf("Expects content has been set. Got %v", r.Content())
	}
}

func TestFactoryResponse_WithContent(t *testing.T) {
	r := &FactoryResponse{}
	c := make(map[string]interface{})
	c["my_key"] = "my_value"
	if reflect.DeepEqual(r.WithContent(c).Content(), c) == false {
		t.Errorf("Expects content has been set. Got %v", r.Content())
	}
}

func TestFactoryResponse_Header(t *testing.T) {
	r := &FactoryResponse{}
	h := NewHeader()
	h.Set("content-type", "application/json")
	r.header = h
	if rh, ok := r.Header().Get("content-type"); ok == false || rh != "application/json" {
		t.Errorf("Expects header must be set. Got %v", r.Header())
	}
}

func TestFactoryResponse_WithHeader(t *testing.T) {
	r := &FactoryResponse{}
	h := NewHeader()
	h.Set("content-type", "application/json")
	if rh, ok := r.WithHeader(h).Header().Get("content-type"); ok == false || rh != "application/json" {
		t.Errorf("Expects header must be set. Got %v", r.Header())
	}
}

func TestFactoryResponse_Cookies(t *testing.T) {
	r := &FactoryResponse{}
	c1 := &http.Cookie{Name: "my_c1", Value: "val_c1"}
	c2 := &http.Cookie{Name: "my_c2", Value: "val_c2"}
	r.cookies = make([]*http.Cookie, 2)
	r.cookies[0] = c1
	r.cookies[1] = c2
	if len(r.Cookies()) != 2 || r.Cookies()[0].Name != "my_c1" || r.Cookies()[1].Value != "val_c2" {
		t.Errorf("Expects cookies must be set. Got %v", r.Cookies())
	}
}

func TestFactoryResponse_WithCookies(t *testing.T) {
	r := &FactoryResponse{}
	c1 := &http.Cookie{Name: "my_c1", Value: "val_c1"}
	c2 := &http.Cookie{Name: "my_c2", Value: "val_c2"}
	c := make([]*http.Cookie, 2)
	c[0] = c1
	c[1] = c2
	if len(r.WithCookies(c).Cookies()) != 2 || r.Cookies()[0].Name != "my_c1" || r.Cookies()[1].Value != "val_c2" {
		t.Errorf("Expects cookies must be set. Got %v", r.Cookies())
	}
}

func TestFactoryResponse_Send(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		rs := &FactoryResponse{w: w, ParserManager: NewParserManager()}
		rs.WithStatus(http.StatusInternalServerError)
		rs.WithHeader(NewHeader()).Header().Set("content-type", "application/json")
		rs.Header().Set("x-custom-flag", "1234")
		c := make(map[string]interface{})
		c["status"] = true
		c["errors"] = false
		rs.WithContent(c)
		rs.WithCookies([]*http.Cookie{
			&http.Cookie{Name: "my_c1", Value: "my_val1"},
			&http.Cookie{Name: "my_c2", Value: "my_val2"},
		})
		rs.WithParser(&parser.JsonParser{})
		rs.Send()
	}
	rq := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	h(w, rq)

	r := w.Result()
	if r.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expects http status to be http.StatusInternalServerError. Got %v", r.StatusCode)
	}
	if r.Header.Get("x-custom-flag") != "1234" || r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expects headers must be set precisely. Got %v - %v", r.Header.Get("x-custom-flag"), r.Header.Get("Content-Type"))
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Expects err to be nil")
	}

	var v map[string]interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		t.Errorf("Expects err to be nil. Got %v", err)
	}
	for key, value := range v {
		if key == "errors" && value.(bool) != false {
			t.Errorf("Expects errors to be false")
		}
		if key == "status" && value.(bool) != true {
			t.Errorf("Expects status to be true")
		}
	}

	if r.Cookies()[0].Name != "my_c1" || r.Cookies()[0].Value != "my_val1" || r.Cookies()[1].Name != "my_c2" || r.Cookies()[1].Value != "my_val2" {
		t.Errorf("Expects cookies have been set. Got %v", r.Cookies())
	}
}
