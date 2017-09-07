package lapi

// Loader is an application loader which could be useful for set things up
type Loader interface {
	// Load runs when application is starting up
	Load(app App)
}
