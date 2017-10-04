package lapi

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Route acts a route describer
type Route interface {
	RouteTagger
	RouteHooker
	RouteHandler
	RouteMatcher
	RouteDescriber
	RouteIdentifier
}

// RouteIdentifier identifies route
type RouteIdentifier interface {
	// Name returns route's name
	Name() string

	// WithName applies route's name
	WithName(name string) Route
}

// RouteDescriber describes route's information
type RouteDescriber interface {
	// Host gives HTTP Host
	Host() string

	// WithHost sets route's host
	WithHost(host string) Route

	// Method returns HTTP Method string
	Method() string

	// WithMethod sets route's method
	WithMethod(method string) Route

	// Uri gives HTTP Uri
	Uri() string

	// WithUri sets route's uri
	WithUri(uri string) Route
}

// RouteHandler manages route's handler
type RouteHandler interface {
	// Handler shows Handler of this route
	Handler() Handler

	// WithHandler sets route's handler
	WithHandler(handler Handler) Route
}

// RouteMatcher matches request
type RouteMatcher interface {
	// Match tests and returns matched route for proposed request
	Match(request Request) (Route, bool)
}

// RouteHooker manages route's hooks
type RouteHooker interface {
	// Hooks returns all hooks for route
	Hooks() map[int]*Slice

	// WithHooks allows to add hooks
	WithHooks(hooks ...Hook) Route

	// WithHook add a single hook
	WithHook(hook Hook) Route
}

// RouteTagger lets route become taggable
type RouteTagger interface {
	// Tags returns all tags of route
	Tags() []string

	// WithTag adds a tag to route
	WithTag(tag string) Route

	// WithTags sets route's tags
	WithTags(tags ...string) Route
}

func NewRoute(method string, uri string, handler Handler) Route {
	r := &FactoryRoute{
		pvHost: &patternVerifier{},
		pvUri:  &patternVerifier{},
		hooks:  make(map[int]*Slice),
		tags:   make([]string, 0),
	}
	return r.
		WithMethod(method).
		WithUri(uri).
		WithHandler(handler).
		WithName(r.genRouteName())
}

type FactoryRoute struct {
	name           string
	host           string
	method         string
	uri            string
	handler        Handler
	hooks          map[int]*Slice
	pvHost         *patternVerifier
	pvUri          *patternVerifier
	requestInput   reflect.Type
	responseOutput reflect.Type
	tags           []string
}

type patternVerifier struct {
	pattern string
	keys    []string
	reg     *regexp.Regexp
}

func (r *FactoryRoute) Name() string {
	return r.name
}

func (r *FactoryRoute) WithName(name string) Route {
	r.name = name
	return r
}

func (r *FactoryRoute) Method() string {
	return r.method
}

func (r *FactoryRoute) WithMethod(method string) Route {
	r.method = strings.ToUpper(method)
	return r
}

func (r *FactoryRoute) Host() string {
	return r.host
}

func (r *FactoryRoute) WithHost(host string) Route {
	var err error
	r.host = host
	r.pvHost.pattern, r.pvHost.keys = r.extractKeyPattern(host)
	r.pvHost.reg, err = regexp.Compile(r.pvHost.pattern)
	PanicOnError(err)
	return r
}

func (r *FactoryRoute) Uri() string {
	return r.uri
}

func (r *FactoryRoute) WithUri(uri string) Route {
	var err error
	r.uri = uri
	r.pvUri.pattern, r.pvUri.keys = r.extractKeyPattern(uri)
	r.pvUri.reg, err = regexp.Compile(r.pvUri.pattern)
	PanicOnError(err)
	return r
}

func (r *FactoryRoute) Handler() Handler {
	return r.handler
}

func (r *FactoryRoute) WithHandler(handler Handler) Route {
	r.handler = handler
	return r
}

func (r *FactoryRoute) Hooks() map[int]*Slice {
	return r.hooks
}

func (r *FactoryRoute) WithHooks(hooks ...Hook) Route {
	r.hooks = make(map[int]*Slice)
	for _, hook := range hooks {
		r.WithHook(hook)
	}
	return r
}

func (r *FactoryRoute) WithHook(hook Hook) Route {
	p := PRIORITY_DEFAULT
	if h, ok := hook.(Prioritizer); ok == true {
		p = h.Priority()
	}

	if r.hooks[p] == nil {
		r.hooks[p] = new(Slice)
	}

	r.hooks[p].Append(hook)
	return r
}

func (r *FactoryRoute) Match(request Request) (Route, bool) {
	method := request.Method()
	host := request.Host()
	uri := request.Uri()
	if !r.matchMethod(method) || !r.matchHost(host) || !r.matchUri(uri) {
		return nil, false
	}

	r.modifyRequestOnMatch(request, r.pvHost, host)
	r.modifyRequestOnMatch(request, r.pvUri, uri)
	return r, true
}

func (r *FactoryRoute) RequestInput() interface{} {
	return Clone(r.requestInput)
}

func (r *FactoryRoute) WithRequestInput(input interface{}) Route {
	r.requestInput = StructOf(input)
	return r
}

func (r *FactoryRoute) ResponseOutput() interface{} {
	return Clone(r.responseOutput)
}

func (r *FactoryRoute) WithResponseOutput(output interface{}) Route {
	r.responseOutput = StructOf(output)
	return r
}

func (r *FactoryRoute) Tags() []string {
	return r.tags
}

func (r *FactoryRoute) WithTag(tag string) Route {
	r.tags = append(r.tags, tag)
	return r
}

func (r *FactoryRoute) WithTags(tags ...string) Route {
	r.tags = tags
	return r
}

func (r *FactoryRoute) genRouteName() string {
	return fmt.Sprintf("%s_%s", r.Method(), strings.Replace(r.Uri(), "/", "_", -1))
}

func (r *FactoryRoute) extractKeyPattern(pattern string) (string, []string) {
	p := `(\<(\w+):([^\>]+)\>)`
	re, err := regexp.Compile(p)
	PanicOnError(err)
	if !re.MatchString(pattern) {
		return pattern, make([]string, 0)
	}
	v := re.FindAllStringSubmatch(pattern, -1)
	keys := make([]string, len(v))
	for i, m := range v {
		keys[i] = m[2]
		pattern = strings.Replace(pattern, m[1], fmt.Sprintf("(%s)", m[3]), 1)
	}
	return pattern, keys
}

func (r *FactoryRoute) matchMethod(method string) bool {
	if r.method == "" {
		return true
	}

	return r.method == method
}

func (r *FactoryRoute) matchHost(host string) bool {
	if r.host == "" {
		return true
	}

	return r.pvHost.reg.MatchString(host)
}

func (r *FactoryRoute) matchUri(uri string) bool {
	if r.uri == "" && uri != "" {
		return false
	}

	return r.pvUri.reg.MatchString(uri)
}

func (r *FactoryRoute) modifyRequestOnMatch(request Request, pv *patternVerifier, s string) {
	numKeys := len(pv.keys)
	if numKeys == 0 {
		return
	}

	m := pv.reg.FindStringSubmatch(s)
	numMatches := len(m)
	if numKeys+1 != numMatches {
		return
	}
	for i, key := range pv.keys {
		request.WithParam(key, m[i+1])
	}
}
