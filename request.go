package lapi

// Request represents for an application's request
type Request interface {
	RequestResolver
}

// RequestResolver returns routing information
type RequestResolver interface {
	// Route returns matched route for request
	Route() Route
}