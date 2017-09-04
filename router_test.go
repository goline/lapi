package lapi

import (
	"net/http"
	"testing"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter()
	if r == nil {
		t.Errorf("Expects r is not nil")
	}
}

func TestNewGroupRouter(t *testing.T) {
	r := NewRouter()
	g := NewGroupRouter(r, "/account")
	if g == nil {
		t.Errorf("Expects g is not nil")
	}
}

func TestFactoryRouter_ByName(t *testing.T) {
	r := NewRouter()
	r.Get("/test", nil)
	r.Get("/test/2", nil).WithName("abc")
	if route, ok := r.ByName("GET__test"); ok == false || route.Uri() != "/test" {
		t.Errorf("Expects uri is /test. Got %v", route)
	}
	if route, ok := r.ByName("abc"); ok == false || route.Uri() != "/test/2" {
		t.Errorf("Expects uri is /test/2. Got %v", route)
	}
}

func TestFactoryRouter_Remove(t *testing.T) {
	r := NewRouter()
	r.Get("/test", nil)
	r.Get("/test/2", nil).WithName("abc")
	if len(r.Routes()) != 2 {
		t.Errorf("Expects router has 2 routes. Got %d", len(r.Routes()))
	}
	r.Remove("abc")
	if len(r.Routes()) != 1 {
		t.Errorf("Expects router has 1 routes. Got %d", len(r.Routes()))
	}
	_, ok := r.ByName("abc")
	if ok == true {
		t.Errorf("Expects route's name abc is not removed")
	}
}

func TestFactoryRouter_Set(t *testing.T) {
	r := NewRouter()
	r.Get("/test", nil)
	r.Get("/test/2", nil).WithName("abc")
	r.Set("abc", NewRoute("GET", "/test/abc", nil))
	if route, ok := r.ByName("abc"); ok == false || route.Uri() != "/test/abc" || route.Name() != "abc" {
		t.Errorf("Expects route is set correctly. Got %v", route)
	}
}

func TestFactoryRouter_Group(t *testing.T) {
	r := NewRouter()
	g := r.Group("/auth")
	route := g.Get("/user", nil)
	if route.Uri() != "/auth/user" {
		t.Errorf("Expects group router is working as expected. Got %s", route.Uri())
	}
}

func TestFactoryRouter_Get(t *testing.T) {
	r := NewRouter()
	r.Get("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodGet {
		t.Errorf("Expects router could register GET method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Post(t *testing.T) {
	r := NewRouter()
	r.Post("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodPost {
		t.Errorf("Expects router could register POST method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Put(t *testing.T) {
	r := NewRouter()
	r.Put("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodPut {
		t.Errorf("Expects router could register PUT method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Patch(t *testing.T) {
	r := NewRouter()
	r.Patch("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodPatch {
		t.Errorf("Expects router could register PATCH method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Delete(t *testing.T) {
	r := NewRouter()
	r.Delete("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodDelete {
		t.Errorf("Expects router could register DELETE method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Head(t *testing.T) {
	r := NewRouter()
	r.Head("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodHead {
		t.Errorf("Expects router could register HEAD method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Options(t *testing.T) {
	r := NewRouter()
	r.Options("/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodOptions {
		t.Errorf("Expects router could register OPTIONS method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Register(t *testing.T) {
	r := NewRouter()
	r.Register(http.MethodPost, "/test", nil).WithName("m")
	route, ok := r.ByName("m")
	if ok == false || route.Method() != http.MethodPost {
		t.Errorf("Expects router could register POST method. Got %s", route.Method())
	}
}

func TestFactoryRouter_Routes(t *testing.T) {
	r := NewRouter()
	r.Get("/test", nil)
	r.Register(http.MethodPost, "/test", nil)
	if len(r.Routes()) != 2 {
		t.Errorf("Expects router has 2 routes. Got %d", len(r.Routes()))
	}
}

func TestFactoryRouter_Duplicate_RouteName(t *testing.T) {
	r := NewRouter()
	defer func(router Router) {
		if r := recover(); r != nil {
			if len(router.Routes()) != 1 {
				t.Errorf("Expects only one route is registered")
			}
		}
	}(r)
	r.Get("/test", nil)
	r.Get("/test", nil)
}

func TestFactoryRouter_Route(t *testing.T) {
	r := &FactoryRouter{routes: make([]Route, 0)}
	r.Get("/test", nil).WithName("Get_Test")
	r.Post("/test", nil).WithName("Post_Test")
	req, _ := NewRequest(nil)
	req.WithMethod(http.MethodPost).WithUri("/test")
	err := r.Route(req)
	if err != nil {
		t.Errorf("Expects err to be nil. Got %v", err)
	}
	if req.Route().Name() != "Post_Test" {
		t.Errorf("Expects matched route's name to be Post_Test. Got %s", req.Route().Name())
	}
}

func TestFactoryRouter_Route_NotFound(t *testing.T) {
	r := &FactoryRouter{routes: make([]Route, 0)}
	r.Get("/test", nil).WithName("Get_Test")
	req, _ := NewRequest(nil)
	req.WithMethod(http.MethodPost).WithUri("/test")
	err := r.Route(req)
	if err == nil {
		t.Errorf("Expects err to be not nil")
	}
	if e, ok := err.(SystemError); ok == false || e.Code() != ERROR_HTTP_NOT_FOUND {
		t.Errorf("Expects err code is ERROR_HTTP_NOT_FOUND. Got %v", err)
	}
}
