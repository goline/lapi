package lapi

import (
	"testing"
)

func TestNewRoute(t *testing.T) {
	r := NewRoute("GET", "v1", "/test//example", nil)
	if _, ok := r.(Route); ok == false {
		t.Errorf("Expect an instance of Route. Got %+v", r)
	}
	if r.Name() != "GET_v1__test__example" {
		t.Errorf("Expects route's name is GET_v1__test__example. Got %v", r.Name())
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
	r := &FactoryRoute{}
	r.WithHost("abc.com:8888")
	if r.host != "abc.com:8888" {
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
	r := &FactoryRoute{}
	r.WithUri("/test/api")
	if r.uri != "/test/api" {
		t.Errorf("Expects route's uri is /test/api. Got %v", r.uri)
	}
}

func TestFactoryRoute_Version(t *testing.T) {
	r := &FactoryRoute{}
	r.version = "V1"
	if r.Version() != "V1" {
		t.Errorf("Expects route's version is V1. Got %v", r.Version())
	}
}

func TestFactoryRoute_WithVersion(t *testing.T) {
	r := &FactoryRoute{}
	r.WithVersion("V1.1")
	if r.version != "V1.1" {
		t.Errorf("Expects route's version is V1.1. Got %v", r.version)
	}
}

type routeHandler struct{}

func (h *routeHandler) Handle(req Request, res Response) (interface{}, error) {
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

func (h *routeHook) SetUp(req Request, res Response) bool                        { return false }
func (h *routeHook) TearDown(req Request, res Response, result interface{}) bool { return false }

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
