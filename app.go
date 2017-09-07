package lapi

import (
	"fmt"
	"net/http"
)

// App is a central application
type App interface {
	AppLoader
	AppRunner
	AppRouter
	AppConfigger
	ContainerAware
	AppErrorHandler
}

// AppLoader handles application's loader
type AppLoader interface {
	// WithLoader allows to register application's loader
	WithLoader(loader Loader) App
}

// AppConfigger handles config
type AppConfigger interface {
	// Config returns instance of config
	Config() Config

	// WithConfig allows to set config
	WithConfig(config Config) App
}

// AppRouter handles router
type AppRouter interface {
	// Router returns an instance of Router
	Router() Router

	// WithRouter sets router
	WithRouter(router Router) App
}

// AppErrorHandler manages error handler
type AppErrorHandler interface {
	// ErrorHandler returns an instance of ErrorHandler
	ErrorHandler() ErrorHandler

	// WithErrorHandler sets error handler
	WithErrorHandler(handler ErrorHandler) App
}

// AppRunner runs application
type AppRunner interface {
	// Run brings application up
	// Any errors should manage inside this method
	Run()
}

func NewApp() App {
	return &FactoryApp{
		config:    NewConfig(),
		container: NewContainer(),
		loaders:   make([]Loader, 0),
	}
}

type FactoryApp struct {
	config       Config
	container    Container
	loaders      []Loader
	request      Request
	response     Response
	router       Router
	errorHandler ErrorHandler
}

func (a *FactoryApp) Container() Container {
	return a.container
}

func (a *FactoryApp) WithContainer(container Container) ContainerAware {
	a.container = container
	return a
}

func (a *FactoryApp) WithLoader(loader Loader) App {
	a.loaders = append(a.loaders, loader)
	return a
}

func (a *FactoryApp) Config() Config {
	return a.config
}

func (a *FactoryApp) WithConfig(config Config) App {
	a.config = config
	return a
}

func (a *FactoryApp) Request() Request {
	return a.request
}

func (a *FactoryApp) Response() Response {
	return a.response
}

func (a *FactoryApp) Router() Router {
	return a.router
}

func (a *FactoryApp) WithRouter(router Router) App {
	a.router = router
	return a
}

func (a *FactoryApp) ErrorHandler() ErrorHandler {
	if a.errorHandler == nil {
		a.errorHandler = NewErrorHandler()
	}

	return a.errorHandler
}

func (a *FactoryApp) WithErrorHandler(handler ErrorHandler) App {
	a.errorHandler = handler
	return a
}

func (a *FactoryApp) Run() {
	a.setUp().handle()
}

func (a *FactoryApp) setUp() *FactoryApp {
	if a.container == nil {
		panic("App requires a container to run")
	}
	for _, loader := range a.loaders {
		loader.Load(a)
	}
	a.container.Inject(a.errorHandler)
	return a
}

func (a *FactoryApp) handle() *FactoryApp {
	if a.router == nil {
		panic(NewSystemError(ERROR_ROUTER_NOT_DEFINED, fmt.Sprint("Router is not defined yet.")))
	}

	http.Handle("/", a)
	if c, ok := a.config.(ServerConfig); ok == true {
		http.ListenAndServe(c.Address(), nil)
	} else {
		panic(NewSystemError(ERROR_SERVER_CONFIG_MISSING, fmt.Sprint("Server configuration is missing")))
	}
	return a
}

func (a *FactoryApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection := a.setUpConnection(w, r)
	defer a.forceSendResponse(connection)

	err := a.router.Route(connection.Request())
	if err != nil {
		a.handleError(connection, err)
		return
	}

	for _, hook := range connection.Request().Route().Hooks() {
		if hook.SetUp(connection) == false {
			break
		}
	}
	if connection.Response().IsSent() == true {
		return
	}

	handler := connection.Request().Route().Handler()
	if handler == nil {
		a.handleError(connection, NewSystemError(ERROR_NO_HANDLER_FOUND, "No handler found"))
		return
	}

	a.container.Inject(handler)
	if h, ok := handler.(ContainerAware); ok == true {
		h.WithContainer(a.container)
	}
	result, err := handler.Handle(connection)
	if err != nil {
		a.handleError(connection, err)
		return
	}
	for _, hook := range connection.Request().Route().Hooks() {
		if hook.TearDown(connection, result, err) == false {
			break
		}
	}
}

func (a *FactoryApp) forceSendResponse(connection Connection) {
	if connection.Response().IsSent() == false {
		connection.Response().Send()
	}
}

func (a *FactoryApp) setUpConnection(w http.ResponseWriter, r *http.Request) Connection {
	request, err := NewRequest(r)
	if err != nil {
		a.handleError(nil, err)
	}

	response, err := NewResponse(w)
	if err != nil {
		a.handleError(nil, err)
	}

	return NewConnection(request, response)
}

func (a *FactoryApp) handleError(connection Connection, err error) {
	PanicOnError(a.errorHandler.HandleError(connection, err))
}
