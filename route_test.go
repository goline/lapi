package lapi

import (
	"testing"
)

func TestNewRoute(t *testing.T) {
	r := NewRoute("GET", "/v1/test//example", nil)
	if _, ok := r.(Route); ok == false {
		t.Errorf("Expect an instance of Route. Got %+v", r)
	}
	if r.Name() != "GET__v1_test__example" {
		t.Errorf("Expects route's name is GET__v1_test__example. Got %v", r.Name())
	}
}

func TestFactoryRoute_Name(t *testing.T) {
	r := &FactoryRoute{}
	r.name = "my_name"
	if r.Name() != "my_name" {
		t.Errorf("Expects route's name is my_name. Got %v", r.Name())
	}
}

func TestFactoryRoute_WithName(t *testing.T) {
	r := &FactoryRoute{}
	r.WithName("my_name")
	if r.name != "my_name" {
		t.Errorf("Expects route's name is my_name. Got %v", r.name)
	}
}

func TestFactoryRoute_Host(t *testing.T) {
	r := &FactoryRoute{}
	r.host = "abc.com:8888"
	if r.Host() != "abc.com:8888" {
		t.Errorf("Expects route's host is abc.com:8888. Got %v", r.Host())
	}
}

func TestFactoryRoute_WithHost(t *testing.T) {
	r := &FactoryRoute{pvHost: &patternVerifier{}}
	p := "({locale:[a-z]{2}}).abc.com:8888"
	r.WithHost(p)
	if r.host != p {
		t.Errorf("Expects route's host is abc.com:8888. Got %v", r.host)
	}
}

func TestFactoryRoute_Method(t *testing.T) {
	r := &FactoryRoute{}
	r.method = "PUT"
	if r.Method() != "PUT" {
		t.Errorf("Expects route's method is PUT. Got %v", r.Method())
	}
}

func TestFactoryRoute_WithMethod(t *testing.T) {
	r := &FactoryRoute{}
	r.WithMethod("PUT")
	if r.method != "PUT" {
		t.Errorf("Expects route's method is PUT. Got %v", r.method)
	}
}

func TestFactoryRoute_Uri(t *testing.T) {
	r := &FactoryRoute{}
	r.uri = "/test/api"
	if r.Uri() != "/test/api" {
		t.Errorf("Expects route's uri is /test/api. Got %v", r.Uri())
	}
}

func TestFactoryRoute_WithUri(t *testing.T) {
	t.SkipNow()
	r := &FactoryRoute{}
	r.WithUri("/test/api")
	if r.uri != "/test/api" {
		t.Errorf("Expects route's uri is /test/api. Got %v", r.uri)
	}
}

type routeHandler struct{}

func (h *routeHandler) Handle(connection Connection) (interface{}, error) {
	return nil, nil
}

func TestFactoryRoute_Handler(t *testing.T) {
	r := &FactoryRoute{}
	r.handler = &routeHandler{}
	if r.Handler() == nil {
		t.Errorf("Expects route's handler is not nil. Got %v", r.Handler())
	}
}

func TestFactoryRoute_WithHandler(t *testing.T) {
	r := &FactoryRoute{}
	r.WithHandler(&routeHandler{})
	if r.handler == nil {
		t.Errorf("Expects route's handler is not nil. Got %v", r.handler)
	}
}

type routeHook struct{}

func (h *routeHook) SetUp(connection Connection) bool { return false }
func (h *routeHook) TearDown(connection Connection, result interface{}, err error) bool {
	return false
}

func TestFactoryRoute_Hooks(t *testing.T) {
	r := &FactoryRoute{}
	r.hooks = make([]Hook, 2)
	r.hooks[0] = &routeHook{}
	r.hooks[1] = &routeHook{}

	if len(r.Hooks()) != 2 {
		t.Errorf("Expects number of route's hooks is 2. Got %v", len(r.Hooks()))
	}
}

func TestFactoryRoute_WithHooks(t *testing.T) {
	r := &FactoryRoute{}
	if len(r.WithHooks(&routeHook{}, &routeHook{}).Hooks()) != 2 {
		t.Errorf("Expects number of route's hooks is 2. Got %v", len(r.Hooks()))
	}
}

func TestFactoryRoute_Match_HostEmpty(t *testing.T) {
	req, _ := NewRequest(nil)
	req.WithHost("domain.com").
		WithUri("/test/abc")
	r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
	r.WithUri("/test/({id:\\d+})")

	_, ok := r.Match(req)
	if ok == true {
		t.Errorf("Expects ok to be false")
	}
}

func TestFactoryRoute_Match_HostNotEmpty(t *testing.T) {
	req, _ := NewRequest(nil)
	req.WithHost("en.domain.com").
		WithUri("/test/55")
	r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
	r.WithUri("/test/({id:\\d+})").WithHost("({locale:[a-z]{2}}).domain.com")

	_, ok := r.Match(req)
	if ok == false {
		t.Errorf("Expects ok to be true")
	}

	locale, ok := req.Param("locale")
	if ok == false || locale != "en" {
		t.Errorf("Expects ok to be true and locale is en. Got %s", locale)
	}
}

func TestFactoryRoute_Match_UriEmpty(t *testing.T) {
	req, _ := NewRequest(nil)
	req.WithHost("domain.com").
		WithUri("/test/abc")
	r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}

	_, ok := r.Match(req)
	if ok == true {
		t.Errorf("Expects ok to be false")
	}
}

func TestFactoryRoute_Match_ZeroKeys(t *testing.T) {
	req := &FactoryRequest{params: NewBag()}
	req.WithHost("domain.com").
		WithUri("/test/abc")
	r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
	r.WithUri("/test/abc")

	_, ok := r.Match(req)
	if ok == false {
		t.Errorf("Expects ok to be true")
	}
	if len(req.params.All()) > 0 {
		t.Errorf("Expects no params in request")
	}
}
