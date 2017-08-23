package lapi

// Route acts a route describer
type Route interface {

	RouteHandler
}

type RouteIdentifier interface {
	// Name returns route's name
	Name() string
}

type RouteDescriber interface {
	// Method returns HTTP Method string
	Method() string

	// Uri gives HTTP Uri
	Uri() string

	// Version returns API version
	Version() string
}

type RouteHandler interface {
	// Handler shows Handler of this route
	Handler() Handler
}