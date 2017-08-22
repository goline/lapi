package lapi

// Request represents for an application's request
type Request interface {
	RequestIdentifier
	RequestVersioner
	RequestResolver
	RequestInput
}

// RequestVersioner handles multiple versions
type RequestVersioner interface {
	// Version returns version string, e.g, v1, v1.1
	Version() string
}

// RequestIdentifier identifies a request
type RequestIdentifier interface {
	// Method returns request's method
	Method() string

	// Uri returns request's path
	Uri() string
}

// RequestResolver returns routing information
type RequestResolver interface {
	// Route returns matched route for request
	Route() Route
}

// RequestInput contains request's input
type RequestInput interface {
	// Input returns an instance of Bag
	Input() Bag
}