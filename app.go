package lapi

// App is a central application
type App interface {
	// Container returns an instance of Container
	Container() Container

	// Register allows to register application's loader
	Register(loader Loader) App

	// Request returns application's request
	Request() Requester

	// Response returns application's response
	Response() Responser
}

// Loader is an application loader which could be useful for set things up
type Loader interface {
	// SetUp runs when application is booting
	SetUp(app App)

	// TearDown runs when application is encountered an error
	TearDown(app App, err Error)
}