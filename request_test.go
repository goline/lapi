package lapi

import (
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://domain.com:8888/test/user?k1=v1&k2=v2#next", nil)
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

	route := NewRoute("get", "/test", nil)
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
	route := NewRoute("get", "/test", nil)
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
	r := &FactoryRequest{header: NewHeader(), Body: NewBody()}
	r.ancestor = a
	r.parseRequest()
	if v, ok := r.Header().Get("Content-Type"); ok == false || v != "application/json" {
		t.Errorf("Expects header content-type equals application/json")
	}
}

func TestFactoryRequest_WithHeader(t *testing.T) {
	h := NewHeader()
	h.Set("content-Type", "application/json")
	r := &FactoryRequest{header: NewHeader()}
	r.WithHeader(h)
	if v, ok := r.Header().Get("Content-Type"); ok == false || v != "application/json" {
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
	if len(r.Cookies()) != 2 || r.Cookies()["k1"].Value != "v1" || r.Cookies()["k2"].Value != "v2" {
		t.Errorf("Expects cookies is set correctly. Got %v", r.Cookies())
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

func TestFactoryRequest_Scheme(t *testing.T) {
	r := &FactoryRequest{}
	r.scheme = SCHEME_HTTPS
	if r.Scheme() != SCHEME_HTTPS {
		t.Errorf("Expects scheme to be https. Got %v", r.Scheme())
	}
}

func TestFactoryRequest_WithScheme(t *testing.T) {
	r := &FactoryRequest{}
	r.WithScheme(SCHEME_HTTPS)
	if r.scheme != SCHEME_HTTPS {
		t.Errorf("Expects scheme to be https. Got %v", r.scheme)
	}
}

func TestFactoryRequest_Host(t *testing.T) {
	r := &FactoryRequest{}
	r.host = "domain.com:888"
	if r.Host() != "domain.com:888" {
		t.Errorf("Expects host to be domain.com:888. Got %v", r.Host())
	}
}

func TestFactoryRequest_WithHost(t *testing.T) {
	r := &FactoryRequest{}
	r.WithHost("domain.com:888")
	if r.host != "domain.com:888" {
		t.Errorf("Expects host to be domain.com:888. Got %v", r.host)
	}
}

func TestFactoryRequest_Port(t *testing.T) {
	r := &FactoryRequest{}
	r.port = 888
	if r.Port() != 888 {
		t.Errorf("Expects port to be 888. Got %d", r.Port())
	}
}

func TestFactoryRequest_WithPort(t *testing.T) {
	r := &FactoryRequest{}
	r.WithPort(888)
	if r.port != 888 {
		t.Errorf("Expects port to be 888. Got %d", r.port)
	}
}

func TestFactoryRequest_Uri(t *testing.T) {
	r := &FactoryRequest{}
	r.uri = "/test/user"
	if r.Uri() != "/test/user" {
		t.Errorf("Expects uri to be /test/user. Got %v", r.Uri())
	}
}

func TestFactoryRequest_WithUri(t *testing.T) {
	r := &FactoryRequest{}
	r.WithUri("/test/user")
	if r.uri != "/test/user" {
		t.Errorf("Expects uri to be /test/user. Got %v", r.uri)
	}
}

func TestFactoryRequest_ParseContentType(t *testing.T) {
	r := &FactoryRequest{header: NewHeader(), Body: NewBody()}
	r.header.Set(HEADER_CONTENT_TYPE, CONTENT_TYPE_XML)
	r.parseContentType()
	if r.ContentType() != CONTENT_TYPE_XML || r.Charset() == "" {
		t.Errorf("Expects content type is CONTENT_TYPE_XML and charset is CONTENT_CHARSET_DEFAULT. Got %s, %s", r.ContentType(), r.Charset())
	}

	r.header.Set(HEADER_CONTENT_TYPE, "application/json; charset=UTF-8")
	r.parseContentType()
	if r.ContentType() != "application/json" || r.Charset() != "UTF-8" {
		t.Errorf("Expects content type is application/json and charset is utf-8. Got %s, %s", r.ContentType(), r.Charset())
	}

	r.header.Set(HEADER_CONTENT_TYPE, "*invalid_content_type")
	r.parseContentType()
	if r.ContentType() != CONTENT_TYPE_DEFAULT || r.Charset() != CONTENT_CHARSET_DEFAULT {
		t.Errorf("Expects content type is CONTENT_TYPE_DEFAULT and charset is CONTENT_CHARSET_DEFAULT. Got %s, %s", r.ContentType(), r.Charset())
	}
}
