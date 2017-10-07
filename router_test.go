package lapi

import (
	"net/http"

	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	It("NewRouter should return an instance of Router", func() {
		Expect(NewRouter()).NotTo(BeNil())
	})

	It("NewGroupRouter should return an instance of Router", func() {
		Expect(NewGroupRouter(NewRouter(), "/account")).NotTo(BeNil())
	})
})

var _ = Describe("FactoryRouter", func() {
	It("ByName should return route by name", func() {
		r := NewRouter()
		r.Get("/test", nil)
		r.Get("/test/2", nil).WithName("abc")

		route, ok := r.ByName("GET__test")
		Expect(ok).To(BeTrue())
		Expect(route.Uri()).To(Equal("/test"))

		route, ok = r.ByName("abc")
		Expect(ok).To(BeTrue())
		Expect(route.Uri()).To(Equal("/test/2"))
	})

	It("Remove should delete a route", func() {
		r := NewRouter()
		r.Get("/test", nil)
		r.Get("/test/2", nil).WithName("abc")
		Expect(len(r.Routes())).To(Equal(2))

		r.Remove("abc")
		Expect(len(r.Routes())).To(Equal(1))

		_, ok := r.ByName("abc")
		Expect(ok).To(BeFalse())
	})

	It("Set should set a route", func() {
		r := NewRouter()
		r.Get("/test", nil)
		r.Get("/test/2", nil).WithName("abc")
		r.Set("abc", NewRoute("GET", "/test/abc", nil))
		route, ok := r.ByName("abc")
		Expect(ok).To(BeTrue())
		Expect(route.Uri()).To(Equal("/test/abc"))
		Expect(route.Name()).To(Equal("abc"))
	})

	It("Group should group routes", func() {
		r := NewRouter()
		g := r.Group("/auth")
		route := g.Get("/user", nil)
		Expect(route.Uri()).To(Equal("/auth/user"))
	})

	It("Any should register empty method route", func() {
		r := NewRouter()
		r.Any("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(""))
	})

	It("Get should register GET route", func() {
		r := NewRouter()
		r.Get("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodGet))
	})

	It("Post should register POST route", func() {
		r := NewRouter()
		r.Post("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodPost))
	})

	It("Put should register PUT route", func() {
		r := NewRouter()
		r.Put("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodPut))
	})

	It("Patch should register PATCH route", func() {
		r := NewRouter()
		r.Patch("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodPatch))
	})

	It("Delete should register DELETE route", func() {
		r := NewRouter()
		r.Delete("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodDelete))
	})

	It("Head should register HEAD route", func() {
		r := NewRouter()
		r.Head("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodHead))
	})

	It("Options should register OPTIONS route", func() {
		r := NewRouter()
		r.Options("/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodOptions))
	})

	It("Register should register a route", func() {
		r := NewRouter()
		r.Register(http.MethodOptions, "/test", nil).WithName("m")
		route, ok := r.ByName("m")
		Expect(ok).To(BeTrue())
		Expect(route.Method()).To(Equal(http.MethodOptions))
	})

	It("Routes should return all routes", func() {
		r := NewRouter()
		r.Get("/test", nil)
		r.Register(http.MethodPost, "/test", nil)
		Expect(len(r.Routes())).To(Equal(2))
	})

	It("Routes should panic with error code ERR_ROUTER_DUPLICATE_ROUTE_NAME", func() {
		r := NewRouter()
		defer func(router Router) {
			if r := recover(); r != nil {
				Expect(len(router.Routes())).To(Equal(1))
				Expect(r.(errors.Error).Code()).To(Equal(ERR_ROUTER_DUPLICATE_ROUTE_NAME))
			}
		}(r)
		r.Get("/test", nil)
		r.Get("/test", nil)
	})

	It("Route should route request", func() {
		r := &FactoryRouter{routes: make([]Route, 0)}
		r.Get("/test", nil).WithName("Get_Test")
		r.Post("/test", nil).WithName("Post_Test")
		req := NewRequest(nil)
		req.WithMethod(http.MethodPost).WithUri("/test")
		err := r.Route(req)
		Expect(err).To(BeNil())
		Expect(req.Route().Name()).To(Equal("Post_Test"))
	})

	It("Route should return error code ERR_HTTP_NOT_FOUND", func() {
		r := &FactoryRouter{routes: make([]Route, 0)}
		r.Get("/test", nil).WithName("Get_Test")
		req := NewRequest(nil)
		req.WithMethod(http.MethodPost).WithUri("/test")
		err := r.Route(req)
		Expect(err).NotTo(BeNil())
		Expect(err.Code()).To(Equal(ERR_HTTP_NOT_FOUND))
	})

	It("WithHook should register hook for all routes", func() {
		r := &FactoryRouter{routes: make([]Route, 0)}
		r.Register("GET", "/test", nil).WithName("my_route")
		r.WithHook(&routeHook{})
		route, ok := r.ByName("my_route")
		Expect(ok).To(BeTrue())
		Expect(len(route.Hooks())).To(Equal(1))
	})

	It("WithTag should register hook for all routes", func() {
		r := &FactoryRouter{routes: make([]Route, 0)}
		r.Register("GET", "/test", nil).WithName("my_route")
		r.WithTag("V2")
		route, ok := r.ByName("my_route")
		Expect(ok).To(BeTrue())
		Expect(len(route.Tags())).To(Equal(1))
	})

	It("WithRoute should register a route", func() {
		r := &FactoryRouter{routes: make([]Route, 0)}
		r.WithRoute(NewRoute("GET", "/test", nil))
		Expect(len(r.routes)).To(Equal(1))
	})

	It("routeIndex should return route position", func() {
		r := &FactoryRouter{routes: make([]Route, 0)}
		r.Register("GET", "/test", nil).WithName("my_route")
		i, ok := r.routeIndex("my_other_route")
		Expect(ok).To(BeFalse())
		Expect(i).To(Equal(-1))
	})

	It("Copy should copy all routes to another router", func() {
		r1 := &FactoryRouter{routes: make([]Route, 0)}
		r2 := &FactoryRouter{routes: make([]Route, 0)}

		r1.Register("Get", "/test", nil)
		r2.Register("Post", "/test", nil)
		Expect(len(r1.routes)).To(Equal(1))
		Expect(len(r2.routes)).To(Equal(1))
		Expect(len(r1.Copy(r2).Routes())).To(Equal(2))
	})
})
