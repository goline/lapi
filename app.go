package lapi

// App is a central application
type App interface {
	// Container returns an instance of Container
	Container() Container

	// Register allows to register application's loader
	Register(loader Loader) App

	// Request returns application's request
	Request() Request

	// Response returns application's response
	Response() Response

	// Run brings application up
	Run(config Config, container Container)
}