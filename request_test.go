package lapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Request", func() {
	It("NewRequest should return an instance of Request", func() {
		req, _ := http.NewRequest("GET", "http://domain.com:8888/test/user?k1=v1&k2=v2#next", nil)
		Expect(NewRequest(req)).NotTo(BeNil())
	})
})

var _ = Describe("FactoryRequest", func() {
	It("Ancestor should return an instance of http.Request", func() {
		req, _ := http.NewRequest("GET", "/test", nil)
		r := &FactoryRequest{}
		r.ancestor = req
		Expect(r.Ancestor().Method).To(Equal("GET"))
	})

	It("Id should return string representing for request's id", func() {
		r := &FactoryRequest{}
		r.id = "0011"
		Expect(r.Id()).To(Equal("0011"))
	})

	It("WithId should set request's id", func() {
		r := &FactoryRequest{}
		r.WithId("0011")
		Expect(r.id).To(Equal("0011"))
	})

	It("Route should return instance of Route", func() {
		r := &FactoryRequest{}
		Expect(r.Route()).To(BeNil())

		route := NewRoute("get", "/test", nil)
		r.route = route
		Expect(r.Route()).NotTo(BeNil())
		Expect(r.Route().Method()).To(Equal("GET"))
	})

	It("WithRoute should set instance of Route", func() {
		r := &FactoryRequest{}
		route := NewRoute("get", "/test", nil)
		r.WithRoute(route)
		Expect(r.Route()).NotTo(BeNil())
		Expect(r.Route().Method()).To(Equal("GET"))
	})

	It("Header should contain headers", func() {
		a, _ := http.NewRequest("GET", "/test", nil)
		a.Header.Set("content-type", "application/json")
		r := &FactoryRequest{header: NewHeader(), body: NewBody(nil, nil)}
		r.ancestor = a
		r.parseRequest()
		v, ok := r.Header().Get("Content-Type")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal("application/json"))
	})

	It("WithHeader should set headers", func() {
		h := NewHeader()
		h.Set("content-Type", "application/json")
		r := &FactoryRequest{header: NewHeader()}
		r.WithHeader(h)
		v, ok := r.Header().Get("Content-Type")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal("application/json"))
	})

	It("Cookie should return cookie", func() {
		r := &FactoryRequest{}
		r.WithCookie(&http.Cookie{Name: "k1", Value: "v1"})
		v, ok := r.Cookie("k1")
		Expect(ok).To(BeTrue())
		Expect(v.Value).To(Equal("v1"))
	})

	It("WithCookie should sets cookie", func() {
		cookie := &http.Cookie{Name: "k1", Value: "v1"}
		r := &FactoryRequest{}
		r.WithCookie(cookie)
		Expect(r.cookies["k1"].Value).To(Equal("v1"))
	})

	It("Cookies should returns cookies", func() {
		r := &FactoryRequest{}
		r.WithCookie(&http.Cookie{Name: "k1", Value: "v1"})
		r.WithCookie(&http.Cookie{Name: "k2", Value: "v2"})
		Expect(len(r.Cookies())).To(Equal(2))
		Expect(r.Cookies()["k1"].Value).To(Equal("v1"))
		Expect(r.Cookies()["k2"].Value).To(Equal("v2"))
	})

	It("WithCookies should sets cookies", func() {
		cookies := []*http.Cookie{
			&http.Cookie{Name: "k1", Value: "v1"},
			&http.Cookie{Name: "k2", Value: "v2"},
		}
		r := &FactoryRequest{}
		r.WithCookies(cookies)
		Expect(len(r.Cookies())).To(Equal(2))
		Expect(r.Cookies()["k1"].Value).To(Equal("v1"))
		Expect(r.Cookies()["k2"].Value).To(Equal("v2"))
	})

	It("Param should return parameter", func() {
		r := &FactoryRequest{}
		r.params = NewBag()
		r.params.Set("found", true)
		v, ok := r.Param("found")
		Expect(ok).To(Equal(true))
		Expect(v).To(Equal(true))
	})

	It("WithParam should set parameter", func() {
		r := &FactoryRequest{}
		r.params = NewBag()
		r.WithParam("found", true)
		v, ok := r.Param("found")
		Expect(ok).To(Equal(true))
		Expect(v).To(Equal(true))
	})

	It("Scheme should return http scheme", func() {
		r := &FactoryRequest{}
		r.scheme = SCHEME_HTTPS
		Expect(r.Scheme()).To(Equal(SCHEME_HTTPS))
	})

	It("WithScheme should set http scheme", func() {
		r := &FactoryRequest{}
		r.WithScheme(SCHEME_HTTPS)
		Expect(r.scheme).To(Equal(SCHEME_HTTPS))
	})

	It("Host should return host", func() {
		r := &FactoryRequest{}
		r.host = "domain.com:888"
		Expect(r.Host()).To(Equal("domain.com:888"))
	})

	It("WithHost should set host", func() {
		r := &FactoryRequest{}
		r.WithHost("domain.com:888")
		Expect(r.host).To(Equal("domain.com:888"))
	})

	It("Port should return port", func() {
		r := &FactoryRequest{}
		r.port = 888
		Expect(r.Port()).To(Equal(888))
	})

	It("WithPort should set port", func() {
		r := &FactoryRequest{}
		r.WithPort(888)
		Expect(r.port).To(Equal(888))
	})

	It("Uri should return uri", func() {
		r := &FactoryRequest{}
		r.uri = "/test/user"
		Expect(r.Uri()).To(Equal("/test/user"))
	})

	It("WithUri should set uri", func() {
		r := &FactoryRequest{}
		r.WithUri("/test/user")
		Expect(r.uri).To(Equal("/test/user"))
	})

	It("parseContentType should parse header content-type", func() {
		r := &FactoryRequest{header: NewHeader(), body: NewBody(nil, nil)}
		r.header.Set(HEADER_CONTENT_TYPE, CONTENT_TYPE_XML)
		r.parseContentType()
		Expect(r.Body().ContentType()).To(Equal(CONTENT_TYPE_XML))
		Expect(r.Body().Charset()).To(Equal(CONTENT_CHARSET_DEFAULT))

		r.header.Set(HEADER_CONTENT_TYPE, "application/json; charset=UTF-8")
		r.parseContentType()
		Expect(r.Body().ContentType()).To(Equal("application/json"))
		Expect(r.Body().Charset()).To(Equal("UTF-8"))

		r.header.Set(HEADER_CONTENT_TYPE, "*invalid_content_type")
		r.parseContentType()
		Expect(r.Body().ContentType()).To(Equal(CONTENT_TYPE_DEFAULT))
		Expect(r.Body().Charset()).To(Equal(CONTENT_CHARSET_DEFAULT))
	})
})
