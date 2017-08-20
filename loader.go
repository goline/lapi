package lapi

// Loader is an application loader which could be useful for set things up
type Loader interface {
	// SetUp runs when application is starting up
	SetUp(app App)

	// TearDown runs when application is encountered an error
	TearDown(app App, err error)
}