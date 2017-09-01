package lapi

import (
	"errors"
	"fmt"
	"net/http"
)

// Router is an application's router
type Router interface {
	RouteRestfuller
	RouteInformer
	RouteRegister
	RouteGrouper
	RouteMatcher
	RouteManager
}

// RouteRestfuller uses common HTTP verbs to register routes
type RouteRestfuller interface {
	// Get registers a GET route handler
	Get(version string, uri string, handler Handler) Route

	// Post registers a POST route handler
	Post(version string, uri string, handler Handler) Route

	// Put registers a PUT route handler
	Put(version string, uri string, handler Handler) Route

	// Patch registers a PATCH route handler
	Patch(version string, uri string, handler Handler) Route

	// Delete registers a DELETE route handler
	Delete(version string, uri string, handler Handler) Route
}

// RouteInformer allows to register special actions, such as Head, Options
type RouteInformer interface {
	// Head registers a HEAD route handler
	Head(version string, uri string, handler Handler) Route

	// Options registers an OPTION route handler
	Options(version string, uri string, handler Handler) Route
}

// RouteRegister lets manually register a route
type RouteRegister interface {
	// Register enrolls a http route handler
	Register(method string, version string, uri string, handler Handler) Route
}

// RouteGrouper groups sub routes
type RouteGrouper interface {
	// Group collects a number of routes
	Group(prefix string) Router
}

// RouteMatcher matches request to route
type RouteMatcher interface {
	// Route performs routing
	Route(request Request) error

	// Match tests and returns matched route for proposed request
	Match(request Request) Route
}

// RouteManager manages inner routes
type RouteManager interface {
	// ByName returns a route by name
	ByName(name string) (Route, bool)

	// Routes returns all registered routes
	Routes() []Route

	// Set allows to set a route to router
	Set(name string, route Route) Router

	// Remove deletes a route by name
	Remove(name string) Router
}

// NewRouter returns an instance of Router
func NewRouter() Router {
	return &FactoryRouter{
		routes: make([]Route, 0),
	}
}

// NewGroupRouter returns a sub (group) router
func NewGroupRouter(parent Router, prefix string) Router {
	return &FactoryRouter{
		parent: parent,
		prefix: prefix,
	}
}

type FactoryRouter struct {
	routes []Route
	parent Router
	prefix string
}

func (r *FactoryRouter) Get(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodGet, version, uri, handler)
}

func (r *FactoryRouter) Post(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodPost, version, uri, handler)
}

func (r *FactoryRouter) Put(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodPut, version, uri, handler)
}

func (r *FactoryRouter) Patch(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodPatch, version, uri, handler)
}

func (r *FactoryRouter) Delete(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodDelete, version, uri, handler)
}

func (r *FactoryRouter) Head(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodHead, version, uri, handler)
}

func (r *FactoryRouter) Options(version string, uri string, handler Handler) Route {
	return r.Register(http.MethodOptions, version, uri, handler)
}

func (r *FactoryRouter) Register(method string, version string, uri string, handler Handler) Route {
	if r.parent != nil && r.prefix != "" {
		uri = fmt.Sprintf("%s%s", r.prefix, uri)
		return r.parent.Register(method, version, uri, handler)
	} else {
		route := NewRoute(method, version, uri, handler)
		_, ok := r.ByName(route.Name())
		if ok == true {
			panic(errors.New(fmt.Sprintf("Route with name %s has already been defined", route.Name())))
		}

		r.routes = append(r.routes, route)
		return route
	}
}

func (r *FactoryRouter) Group(prefix string) Router {
	return NewGroupRouter(r, prefix)
}

func (r *FactoryRouter) ByName(name string) (Route, bool) {
	for _, route := range r.routes {
		if route.Name() == name {
			return route, true
		}
	}
	return nil, false
}

func (r *FactoryRouter) Routes() []Route {
	return r.routes
}

func (r *FactoryRouter) Set(name string, route Route) Router {
	i, ok := r.routeIndex(name)
	if ok == true {
		route.WithName(name)
		r.routes[i] = route
	}
	return r
}

func (r *FactoryRouter) Remove(name string) Router {
	i, ok := r.routeIndex(name)
	if ok == true {
		r.routes = append(r.routes[:i], r.routes[i+1:]...)
	}
	return r
}

func (r *FactoryRouter) Route(request Request) error {
	return nil
}

func (r *FactoryRouter) Match(request Request) Route {
	return nil
}

func (r *FactoryRouter) routeIndex(name string) (int, bool) {
	for i, route := range r.routes {
		if route.Name() == name {
			return i, true
		}
	}
	return -1, false
}
