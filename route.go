package lapi

import (
	"fmt"
	"strings"
)

// Route acts a route describer
type Route interface {
	RouteHooker
	RouteHandler
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

	// Version returns API version
	Version() string

	// WithVersion sets route's version
	WithVersion(version string) Route
}

// RouteHandler manages route's handler
type RouteHandler interface {
	// Handler shows Handler of this route
	Handler() Handler

	// WithHandler sets route's handler
	WithHandler(handler Handler) Route
}

// RouteHooker manages route's hooks
type RouteHooker interface {
	// Hooks returns all hooks for route
	Hooks() []Hook

	// WithHooks allows to add hooks
	WithHooks(hooks ...Hook) Route
}

func NewRoute(method string, version string, uri string, handler Handler) Route {
	r := &FactoryRoute{}
	return r.
		WithMethod(method).
		WithVersion(version).
		WithUri(uri).
		WithHandler(handler).
		WithName(genRouteName(r))
}

type FactoryRoute struct {
	name    string
	host    string
	method  string
	uri     string
	version string
	handler Handler
	hooks   []Hook
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
	r.host = host
	return r
}

func (r *FactoryRoute) Uri() string {
	return r.uri
}

func (r *FactoryRoute) WithUri(uri string) Route {
	r.uri = uri
	return r
}

func (r *FactoryRoute) Version() string {
	return r.version
}

func (r *FactoryRoute) WithVersion(version string) Route {
	r.version = version
	return r
}

func (r *FactoryRoute) Handler() Handler {
	return r.handler
}

func (r *FactoryRoute) WithHandler(handler Handler) Route {
	r.handler = handler
	return r
}

func (r *FactoryRoute) Hooks() []Hook {
	return r.hooks
}

func (r *FactoryRoute) WithHooks(hooks ...Hook) Route {
	r.hooks = hooks
	return r
}

func genRouteName(r Route) string {
	return fmt.Sprintf("%s_%s_%s", r.Method(), r.Version(), strings.Replace(r.Uri(), "/", "_", -1))
}
