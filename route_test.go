package lapi

import (
	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route", func() {
	It("NewRoute should return an instance of Route", func() {
		r := NewRoute("GET", "/v1/test//example", nil)
		Expect(r).NotTo(BeNil())
		Expect(r.Name()).To(Equal("GET__v1_test__example"))
	})
})

type routeHandler struct{}

func (h *routeHandler) Handle(connection Connection) (interface{}, errors.Error) {
	return nil, nil
}

type routeHook struct{}

func (h *routeHook) SetUp(connection Connection) errors.Error { return nil }
func (h *routeHook) TearDown(connection Connection, result interface{}, err errors.Error) errors.Error {
	return nil
}

type sampleRequestInput struct{}
type sampleResponseOutput struct{}

var _ = Describe("FactoryRoute", func() {
	It("Name should return route's name", func() {
		r := &FactoryRoute{}
		r.name = "my_name"
		Expect(r.Name()).To(Equal("my_name"))
	})

	It("WithName should set route's name", func() {
		r := &FactoryRoute{}
		r.WithName("my_name")
		Expect(r.name).To(Equal("my_name"))
	})

	It("Host should return route's host", func() {
		r := &FactoryRoute{}
		r.host = "abc.com:8888"
		Expect(r.Host()).To(Equal("abc.com:8888"))
	})

	It("WithHost should set route's host", func() {
		r := &FactoryRoute{pvHost: &patternVerifier{}}
		p := "<locale:[a-z]{2}>.abc.com:8888"
		r.WithHost(p)
		Expect(r.host).To(Equal(p))
	})

	It("Method should return route's method", func() {
		r := &FactoryRoute{}
		r.method = "PUT"
		Expect(r.Method()).To(Equal("PUT"))
	})

	It("WithMethod should set route's method", func() {
		r := &FactoryRoute{}
		r.WithMethod("PUT")
		Expect(r.method).To(Equal("PUT"))
	})

	It("Uri should return route's uri", func() {
		r := &FactoryRoute{}
		r.uri = "/test/api"
		Expect(r.Uri()).To(Equal("/test/api"))
	})

	It("WithUri should set route's uri", func() {
		r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
		r.WithUri("/test/api")
		Expect(r.uri).To(Equal("/test/api"))
	})

	It("Handler should return route's handler", func() {
		r := &FactoryRoute{}
		r.handler = &routeHandler{}
		Expect(r.Handler()).NotTo(BeNil())
	})

	It("WithHandler should set route's handler", func() {
		r := &FactoryRoute{}
		r.WithHandler(&routeHandler{})
		Expect(r.handler).NotTo(BeNil())
	})

	It("Hooks should return route's hooks", func() {
		r := &FactoryRoute{}
		r.hooks = make(map[int]*Slice, 2)
		r.WithHook(&routeHook{})
		r.WithHook(&routeHook{})

		hooks, ok := r.Hooks()[PRIORITY_DEFAULT]
		Expect(ok).To(BeTrue())
		Expect(len(hooks.All())).To(Equal(2))
	})

	It("WithHooks should set route's hooks", func() {
		r := &FactoryRoute{}
		hooks, ok := r.WithHooks(&routeHook{}, &routeHook{}).Hooks()[0]
		Expect(ok).To(BeTrue())
		Expect(len(hooks.All())).To(Equal(2))
	})

	It("Hook should return route's hook", func() {
		r := &FactoryRoute{hooks: make(map[int]*Slice)}
		Expect(len(r.WithHook(&routeHook{}).Hooks())).To(Equal(1))
	})

	It("Match should match empty host", func() {
		req := NewRequest(nil)
		req.WithHost("domain.com").
			WithUri("/test/abc")
		r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
		r.WithUri("/test/<id:\\d+>")

		_, ok := r.Match(req)
		Expect(ok).To(BeFalse())
	})

	It("Match should verify host not empty", func() {
		req := NewRequest(nil)
		req.WithHost("en.domain.com").
			WithUri("/test/55")
		r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
		r.WithUri("/test/<id:\\d+>").
			WithHost("<locale:[a-z]{2}>.domain.com")

		_, ok := r.Match(req)
		Expect(ok).To(BeTrue())

		locale, ok := req.Param("locale")
		Expect(ok).To(BeTrue())
		Expect(locale).To(Equal("en"))
	})

	It("Match should verify uri empty", func() {
		req := NewRequest(nil)
		req.WithHost("domain.com").
			WithUri("/test/abc")
		r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}

		_, ok := r.Match(req)
		Expect(ok).To(BeFalse())
	})

	It("Match should verify no keys", func() {
		req := &FactoryRequest{params: NewBag()}
		req.WithHost("domain.com").
			WithUri("/test/abc")
		r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
		r.WithUri("/test/abc")

		_, ok := r.Match(req)
		Expect(ok).To(BeTrue())
		Expect(len(req.params.All())).To(BeZero())
	})

	It("Match should request without parameters", func() {
		req := &FactoryRequest{params: NewBag()}
		req.WithUri("/v1/user")
		r := &FactoryRoute{pvHost: &patternVerifier{}, pvUri: &patternVerifier{}}
		r.WithUri("/.*")

		_, ok := r.Match(req)
		Expect(ok).To(BeTrue())
		Expect(len(req.params.All())).To(BeZero())
	})

	It("Tags should return route's tags", func() {
		r := &FactoryRoute{tags: make([]string, 0)}
		Expect(len(r.tags)).To(BeZero())

		r.tags = append(r.tags, "a_tag")
		Expect(len(r.tags)).To(Equal(1))
	})

	It("WithTag should add route tag", func() {
		r := &FactoryRoute{tags: make([]string, 0)}
		Expect(len(r.WithTag("a_tag").Tags())).To(Equal(1))
	})

	It("WithTags should add route tags", func() {
		r := &FactoryRoute{tags: make([]string, 0)}
		Expect(len(r.WithTags("a_tag", "another_tag").Tags())).To(Equal(2))
	})
})
