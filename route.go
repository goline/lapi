package lapi

// Route acts a route describer
type Route interface {
	// Method returns HTTP Method string
	Method() string

	// Uri gives HTTP Uri
	Uri() string

	// Handler shows Handler of this route
	Handler() Handler
}