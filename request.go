package lapi

// Requester represents for an application's request
type Requester interface {
	RequestIdentifier
	RequestResolver
	RequestParameter
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

// RequestParameter helps to parameterize request
type RequestParameter interface {
	// Get returns value of specific key
	Get(key string, def interface{}) interface{}

	// Set allows to set value for a proposed key
	Set(key string, value interface{})

	// Has helps to check if a key exists
	Has(key string) bool
}