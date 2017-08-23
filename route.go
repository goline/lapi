package lapi

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
	WithHooks(hooks ...Hook)
}
