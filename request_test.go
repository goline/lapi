package lapi

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	r := NewRequest(req)
	if r == nil {
		t.Errorf("Expects NewRequest to return an instance of Request. Got nil")
	}
}

func TestFactoryRequest_Ancestor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	r := &FactoryRequest{}
	r.ancestor = req
	if r.Ancestor().Method != "GET" {
		t.Errorf("Expects ancestor's method to be GET. Got %v", r.Ancestor().Method)
	}
}

func TestFactoryRequest_Route(t *testing.T) {
	r := &FactoryRequest{}
	if r.Route() != nil {
		t.Errorf("Expects request's route to be nil. Got %v", r.Route())
	}

	route := NewRoute("get", "v1", "/test", nil)
	r.route = route
	if r.Route() == nil {
		t.Errorf("Expects request's route to be not nil. Got nil")
	}
	if r.Route().Method() != "GET" {
		t.Errorf("Expects request's route's method to be GET. Got %v", r.Route().Method())
	}
}

func TestFactoryRequest_WithRoute(t *testing.T) {
	r := &FactoryRequest{}
	route := NewRoute("get", "v1", "/test", nil)
	r.WithRoute(route)
	if r.Route() == nil {
		t.Errorf("Expects request's route to be not nil. Got nil")
	}
	if r.Route().Method() != "GET" {
		t.Errorf("Expects request's route's method to be GET. Got %v", r.Route().Method())
	}
}

func TestFactoryRequest_Header(t *testing.T) {
	a, _ := http.NewRequest("GET", "/test", nil)
	a.Header.Set("content-type", "application/json")
	r := &FactoryRequest{}
	r.ancestor = a
	r.parseRequest()
	if v, ok := r.Header().Get("Content-Type"); ok == false || v[0] != "application/json" {
		t.Errorf("Expects header content-type equals application/json")
	}
}

func TestFactoryRequest_WithHeader(t *testing.T) {
	h := NewHeader()
	h.Set("content-Type", "application/json")
	r := &FactoryRequest{}
	r.WithHeader(h)
	if v, ok := r.Header().Get("Content-Type"); ok == false || v[0] != "application/json" {
		t.Errorf("Expects header content-type equals application/json")
	}
}

func TestFactoryRequest_Cookie(t *testing.T) {
	r := &FactoryRequest{}
	r.WithCookie(&http.Cookie{Name: "k1", Value: "v1"})
	if c, ok := r.Cookie("k1"); ok == false || c.Value != "v1" {
		t.Errorf("Expects cookie is set correctly. Got %v", c)
	}
}

func TestFactoryRequest_WithCookie(t *testing.T) {
	cookie := &http.Cookie{Name: "k1", Value: "v1"}
	r := &FactoryRequest{}
	r.WithCookie(cookie)
	if r.cookies["k1"] == nil {
		t.Errorf("Expects WithCookie is correct. Got %v", r.cookies)
	}
}

func TestFactoryRequest_Cookies(t *testing.T) {
	r := &FactoryRequest{}
	r.WithCookie(&http.Cookie{Name: "k1", Value: "v1"})
	r.WithCookie(&http.Cookie{Name: "k2", Value: "v2"})
	if len(r.cookies) != 2 || r.cookies["k1"].Value != "v1" || r.cookies["k2"].Value != "v2" {
		t.Errorf("Expects cookies is set correctly. Got %v", r.cookies)
	}
}

func TestFactoryRequest_WithCookies(t *testing.T) {
	cookies := []*http.Cookie{
		&http.Cookie{Name: "k1", Value: "v1"},
		&http.Cookie{Name: "k2", Value: "v2"},
	}
	r := &FactoryRequest{}
	r.WithCookies(cookies)
	if len(r.cookies) != 2 || r.cookies["k1"].Value != "v1" || r.cookies["k2"].Value != "v2" {
		t.Errorf("Expects cookies is set correctly. Got %v", r.cookies)
	}
}

func TestFactoryRequest_Input(t *testing.T) {
	input := map[string]string{"name": "value"}
	r := &FactoryRequest{}
	r.input = input
	if reflect.DeepEqual(r.Input(), input) == false {
		t.Errorf("Expects input is assigned. Got %v", r.Input())
	}
}

func TestFactoryRequest_WithInput(t *testing.T) {
	input := map[string]string{"name": "value"}
	r := &FactoryRequest{}
	r.WithInput(input)
	if reflect.DeepEqual(r.Input(), input) == false {
		t.Errorf("Expects input is assigned. Got %v", r.Input())
	}
}

func TestFactoryRequest_Param(t *testing.T) {
	r := &FactoryRequest{}
	r.params = NewBag()
	if v, ok := r.Param("not_found"); ok == true {
		t.Errorf("Expects key not_found is not found. Got %v", v)
	}

	r.params.Set("found", true)
	if v, ok := r.Param("found"); ok == false || v.(bool) != true {
		t.Errorf("Expects key found is found. Got %v", v)
	}
}

func TestFactoryRequest_WithParam(t *testing.T) {
	r := &FactoryRequest{}
	r.params = NewBag()
	if v, ok := r.Param("not_found"); ok == true {
		t.Errorf("Expects key not_found is not found. Got %v", v)
	}

	r.WithParam("found", true)
	if v, ok := r.Param("found"); ok == false || v.(bool) != true {
		t.Errorf("Expects key found is found. Got %v", v)
	}
}
