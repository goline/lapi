package lapi

import (
	"fmt"
	"net/http"

	"github.com/goline/errors"
)

// Router is an application's router
type Router interface {
	RouteRestfuller
	RouteDispatcher
	RouteInformer
	RouteRegister
	RouteGrouper
	RouteManager
}

// RouteRestfuller uses common HTTP verbs to register routes
type RouteRestfuller interface {
	// Get registers a GET route handler
	Get(uri string, handler Handler) Route

	// Post registers a POST route handler
	Post(uri string, handler Handler) Route

	// Put registers a PUT route handler
	Put(uri string, handler Handler) Route

	// Patch registers a PATCH route handler
	Patch(uri string, handler Handler) Route

	// Delete registers a DELETE route handler
	Delete(uri string, handler Handler) Route
}

// RouteInformer allows to register special actions, such as Head, Options
type RouteInformer interface {
	// Head registers a HEAD route handler
	Head(uri string, handler Handler) Route

	// Options registers an OPTION route handler
	Options(uri string, handler Handler) Route
}

// RouteRegister lets manually register a route
type RouteRegister interface {
	// Any registers route for any methods
	Any(uri string, handler Handler) Route

	// Register enrolls a http route handler
	Register(method string, uri string, handler Handler) Route

	// WithRoute registers a route
	WithRoute(route Route) Router
}

// RouteGrouper groups sub routes
type RouteGrouper interface {
	// Group collects a number of routes
	Group(prefix string) Router

	// WithHook allows to add hook to all router's routes
	WithHook(hook Hook) Router

	// WithTag adds a tag to all routes
	WithTag(tag string) Router
}

// RouteMatcher matches request to route
type RouteDispatcher interface {
	// Route performs routing
	Route(request Request) errors.Error
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

type RouteCopier interface {
	// Copy add all routes of input router to current router
	Copy(router Router) Router
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

func (r *FactoryRouter) Any(uri string, handler Handler) Route {
	return r.Register("", uri, handler)
}

func (r *FactoryRouter) Get(uri string, handler Handler) Route {
	return r.Register(http.MethodGet, uri, handler)
}

func (r *FactoryRouter) Post(uri string, handler Handler) Route {
	return r.Register(http.MethodPost, uri, handler)
}

func (r *FactoryRouter) Put(uri string, handler Handler) Route {
	return r.Register(http.MethodPut, uri, handler)
}

func (r *FactoryRouter) Patch(uri string, handler Handler) Route {
	return r.Register(http.MethodPatch, uri, handler)
}

func (r *FactoryRouter) Delete(uri string, handler Handler) Route {
	return r.Register(http.MethodDelete, uri, handler)
}

func (r *FactoryRouter) Head(uri string, handler Handler) Route {
	return r.Register(http.MethodHead, uri, handler)
}

func (r *FactoryRouter) Options(uri string, handler Handler) Route {
	return r.Register(http.MethodOptions, uri, handler)
}

func (r *FactoryRouter) Register(method string, uri string, handler Handler) Route {
	if r.parent != nil && r.prefix != "" {
		uri = fmt.Sprintf("%s%s", r.prefix, uri)
		route := r.parent.Register(method, uri, handler)
		r.routes = append(r.routes, route)
		return route
	} else {
		route := NewRoute(method, uri, handler)
		_, ok := r.ByName(route.Name())
		if ok == true {
			panic(errors.New(ERR_ROUTER_DUPLICATE_ROUTE_NAME, fmt.Sprintf("Route with name %s has already been defined", route.Name())))
		}

		r.routes = append(r.routes, route)
		return route
	}
}

func (r *FactoryRouter) WithRoute(route Route) Router {
	r.routes = append(r.routes, route)
	return r
}

func (r *FactoryRouter) Group(prefix string) Router {
	return NewGroupRouter(r, prefix)
}

func (r *FactoryRouter) WithHook(hook Hook) Router {
	for _, route := range r.routes {
		route.WithHook(hook)
	}
	return r
}

func (r *FactoryRouter) WithTag(tag string) Router {
	for _, route := range r.routes {
		if tagger, ok := route.(RouteTagger); ok == true {
			tagger.WithTag(tag)
		}
	}
	return r
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

func (r *FactoryRouter) Route(request Request) errors.Error {
	for _, route := range r.routes {
		if matchedRoute, ok := route.Match(request); ok == true {
			request.WithRoute(matchedRoute)
			return nil
		}
	}
	return errors.New(ERR_HTTP_NOT_FOUND, fmt.Sprintf("Url (%s %s) could not be found", request.Method(), request.Uri())).
		WithLevel(errors.LEVEL_WARN)
}

func (r *FactoryRouter) Copy(router Router) Router {
	for _, route := range router.Routes() {
		r.routes = append(r.routes, route)
	}
	return r
}

func (r *FactoryRouter) routeIndex(name string) (int, bool) {
	for i, route := range r.routes {
		if route.Name() == name {
			return i, true
		}
	}
	return -1, false
}
