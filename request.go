package lapi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Request represents for an application's request
type Request interface {
	RequestHeader
	RequestCookies
	RequestInputer
	RequestAncestor
	RequestResolver
	RequestInformer
	RequestParameter
	ParserManager
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

// RequestInformer contains request information
type RequestInformer interface {
	// Method returns request's method
	Method() string

	// WithMethod sets request's method
	WithMethod(method string) Request

	// Scheme returns request's scheme, such as http, https, ftp.
	Scheme() string

	// WithScheme sets request's scheme
	WithScheme(scheme string) Request

	// Host return request's host
	Host() string

	// WithHost sets request's host
	WithHost(host string) Request

	// Port return request's port
	Port() int

	// WithPort sets request's port
	WithPort(port int) Request

	// Uri returns request's uri
	Uri() string

	// WithUri sets request's uri
	WithUri(uri string) Request
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

func NewRequest(req *http.Request) (Request, error) {
	r := &FactoryRequest{
		ancestor:      req,
		cookies:       make(map[string]*http.Cookie),
		params:        NewBag(),
		ParserManager: NewParserManager(),
	}
	if req != nil {
		if err := r.parseRequest(); err != nil {
			return nil, err
		}
	}
	return r, nil
}

type FactoryRequest struct {
	ancestor *http.Request
	header   Header
	input    interface{}
	cookies  map[string]*http.Cookie
	params   Bag
	route    Route
	method   string
	scheme   string
	host     string
	port     int
	uri      string
	ParserManager
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

func (r *FactoryRequest) Method() string {
	return r.method
}

func (r *FactoryRequest) WithMethod(method string) Request {
	r.method = method
	return r
}

func (r *FactoryRequest) Scheme() string {
	return r.scheme
}

func (r *FactoryRequest) WithScheme(scheme string) Request {
	r.scheme = scheme
	return r
}

func (r *FactoryRequest) Host() string {
	return r.host
}

func (r *FactoryRequest) WithHost(host string) Request {
	r.host = host
	return r
}

func (r *FactoryRequest) Port() int {
	return r.port
}

func (r *FactoryRequest) WithPort(port int) Request {
	r.port = port
	return r
}

func (r *FactoryRequest) Uri() string {
	return r.uri
}

func (r *FactoryRequest) WithUri(uri string) Request {
	r.uri = uri
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
	if r.cookies == nil {
		r.cookies = make(map[string]*http.Cookie)
	}
	r.cookies[cookie.Name] = cookie
	return r
}

func (r *FactoryRequest) Cookies() map[string]*http.Cookie {
	return r.cookies
}

func (r *FactoryRequest) WithCookies(cookies []*http.Cookie) Request {
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

func (r *FactoryRequest) parseRequest() error {
	var err error
	err = r.parseRequestAddress()
	if err != nil {
		return err
	}

	err = r.parseRequestHeader()
	if err != nil {
		return err
	}

	err = r.parseRequestBody()
	if err != nil {
		return err
	}

	err = r.parseCookies()
	if err != nil {
		return err
	}

	return nil
}

func (r *FactoryRequest) parseRequestAddress() error {
	r.WithMethod(r.ancestor.Method)
	r.WithHost(r.ancestor.URL.Hostname())
	if p, _ := strconv.Atoi(r.ancestor.URL.Port()); p > 0 {
		r.WithPort(p)
		r.WithScheme(r.ancestor.URL.Scheme)
	} else {
		r.WithPort(PORT_HTTP)
		r.WithScheme(SCHEME_HTTP)
	}
	r.WithUri(r.ancestor.URL.Path)
	q := r.ancestor.URL.Query()
	if len(q) > 0 {
		for k, v := range q {
			r.WithParam(k, v)
		}
	}
	return nil
}

func (r *FactoryRequest) parseRequestHeader() error {
	if r.header == nil {
		r.header = NewHeader()
	}
	for key := range r.ancestor.Header {
		r.Header().Set(key, r.ancestor.Header.Get(key))
	}
	return nil
}

func (r *FactoryRequest) parseRequestBody() error {
	if r.ancestor.Body == nil {
		return nil
	}

	body, err := ioutil.ReadAll(r.ancestor.Body)
	if err != nil {
		return err
	}

	var input interface{}
	ct, ok := r.header.Get("Content-Type")
	if ok == false {
		ct = CONTENT_TYPE_DEFAULT
	}
	p, ok := r.Parser(ct)
	if ok == false {
		return errors.New(fmt.Sprintf("Unable to find an appropriate parser for %s", ct))
	}
	err = p.Decode(body, input)
	if err != nil {
		return err
	}
	r.WithInput(input)

	return nil
}

func (r *FactoryRequest) parseCookies() error {
	r.WithCookies(r.ancestor.Cookies())
	return nil
}
