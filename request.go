package lapi

// Requester represents for an application's request
type Requester interface {
	RequestIdentifier
	RequestResolver
	RequestInput
}

// RequestVersioner handles multiple versions
type RequestVersioner interface {
	// Version returns version string
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
	Bag
}