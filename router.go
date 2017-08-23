package lapi

import "net/http"

// Router is an application's router
type Router interface {
	RouteRestfuller
	RouteInformer
	RouteRegister
	RouteMatcher
	RouteManager
}

// RouteRestfuller uses common HTTP verbs to register routes
type RouteRestfuller interface {
	// Get registers a GET route handler
	Get(path string, handler Handler) Route

	// Post registers a POST route handler
	Post(path string, handler Handler) Route

	// Put registers a PUT route handler
	Put(path string, handler Handler) Route

	// Patch registers a PATCH route handler
	Patch(path string, handler Handler) Route

	// Delete registers a DELETE route handler
	Delete(path string, handler Handler) Route
}

// RouteInformer allows to register special actions, such as Head, Options
type RouteInformer interface {
	// Head registers a HEAD route handler
	Head(path string, handler Handler) Route

	// Options registers an OPTION route handler
	Options(path string, handler Handler) Route
}

// RouteRegister lets manually register a route
type RouteRegister interface {
	// Register enrolls a http route handler
	Register(method string, path string, handler Handler) Route
}

// RouteGrouper groups sub routes
type RouteGrouper interface {
	// Group collects a number of routes
	Group(prefix string) Router
}

// RouteMatcher matches request to route
type RouteMatcher interface {
	// Route performs routing
	Route(r http.Request) Request

	// Match tests and returns matched route for proposed request
	Match(r http.Request) Route
}

// RouteManager manages inner routes
type RouteManager interface {
	// ByName returns a route by name
	ByName(name string) Route

	// Routes returns all registered routes
	Routes() map[string]Route

	// Set allows to set a route to router
	Set(name string, route Route)

	// Remove deletes a route by name
	Remove(name string)
}
