package lapi

import (
	"net/http"
)

// Request represents for an application's request
type Request interface {
	RequestHeader
	RequestCookies
	RequestInputer
	RequestAncestor
	RequestResolver
	RequestParameter
}

// RequestAncestor keeps original http.Request
type RequestAncestor interface {
	Ancestor() *http.Request
}

// RequestResolver returns routing information
type RequestResolver interface {
	// Route returns matched route for request
	Route() Route

	// WithRoute sets request's routes
	WithRoute(route Route) Request
}

// RequestHeader manages request's header
type RequestHeader interface {
	// Header returns an instance of Header
	Header() Header

	// WithHeader allows to set Header
	WithHeader(header Header) Request
}

// RequestCookies handles request's cookies
type RequestCookies interface {
	// Cookie returns an appropriate cookie by name
	Cookie(name string) (*http.Cookie, bool)

	// WithCookie sets cookie
	WithCookie(cookie *http.Cookie) Request

	// Cookies returns all cookies
	Cookies() map[string]*http.Cookie

	// WithCookies sets request's cookies
	WithCookies(cookies []*http.Cookie) Request
}

// RequestInputer handles request's input (body)
type RequestInputer interface {
	// Input returns request's input
	Input() interface{}

	// WithInput sets request's input
	WithInput(input interface{}) Request
}

// RequestParameter handles request's query parameters and additional parameters
type RequestParameter interface {
	// Param returns value of a proposed key if exists, ok will be false if key is not found
	Param(key string) (value interface{}, ok bool)

	// WithParam sets parameter by key
	WithParam(key string, value interface{}) Request
}

func NewRequest(req *http.Request) Request {
	r := &FactoryRequest{
		ancestor: req,
		cookies:  make(map[string]*http.Cookie),
		params:   NewBag(),
	}
	r.parseRequest()
	return r
}

type FactoryRequest struct {
	ancestor *http.Request
	header   Header
	input    interface{}
	cookies  map[string]*http.Cookie
	params   Bag
	route    Route
}

func (r *FactoryRequest) Ancestor() *http.Request {
	return r.ancestor
}

func (r *FactoryRequest) Route() Route {
	return r.route
}

func (r *FactoryRequest) WithRoute(route Route) Request {
	r.route = route
	return r
}

func (r *FactoryRequest) Header() Header {
	return r.header
}

func (r *FactoryRequest) WithHeader(header Header) Request {
	r.header = header
	return r
}

func (r *FactoryRequest) Cookie(name string) (*http.Cookie, bool) {
	cookie, ok := r.cookies[name]
	return cookie, ok
}

func (r *FactoryRequest) WithCookie(cookie *http.Cookie) Request {
	r.cookies[cookie.Name] = cookie
	return r
}

func (r *FactoryRequest) Cookies() map[string]*http.Cookie {
	return r.cookies
}

func (r *FactoryRequest) WithCookies(cookies []*http.Cookie) Request {
	r.cookies = make(map[string]*http.Cookie)
	for _, cookie := range cookies {
		r.WithCookie(cookie)
	}
	return r
}

func (r *FactoryRequest) Input() interface{} {
	return r.input
}

func (r *FactoryRequest) WithInput(input interface{}) Request {
	r.input = input
	return r
}

func (r *FactoryRequest) Param(key string) (value interface{}, ok bool) {
	return r.params.Get(key)
}

func (r *FactoryRequest) WithParam(key string, value interface{}) Request {
	r.params.Set(key, value)
	return r
}

func (r *FactoryRequest) parseRequest() {
	r.parseRequestHeader()
}

func (r *FactoryRequest) parseRequestHeader() {
	for key := range r.ancestor.Header {
		r.Header().Set(key, r.ancestor.Header.Get(key))
	}
}
