package lapi

import (
	"fmt"
	"net/http"
)

// App is a central application
type App interface {
	ContainerAware
	AppLoader
	AppRunner
	AppConfigger
	AppRouter
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

func (a *FactoryApp) WithErrorHandler(handler ErrorHandler) App {
	a.errorHandler = handler
	return a
}

func (a *FactoryApp) Run() {
	for _, loader := range a.loaders {
		loader.SetUp(a)
	}

	http.Handle("/", a)
	if c, ok := a.config.(ServerConfig); ok == true {
		http.ListenAndServe(c.Address(), nil)
	} else {
		panic(NewSystemError(ERROR_SERVER_CONFIG_MISSING, fmt.Sprint("Server configuration is missing")))
	}
}

func (a *FactoryApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := NewRequest(r)
	if err != nil {
		a.sendBadRequestResponse(w)
		return
	}

	response, err := NewResponse(w)
	if err != nil {
		a.sendBadRequestResponse(w)
		return
	}

	connection := NewConnection(request, response)
	defer a.forceSendResponse(connection)

	if a.router == nil {
		a.handleError(connection, NewSystemError(ERROR_ROUTER_NOT_DEFINED, fmt.Sprint("Router is not defined yet.")))
	}

	err = a.router.Route(request)
	if err != nil {
		a.handleError(connection, err)
		return
	}

	for _, hook := range request.Route().Hooks() {
		isContinue := hook.SetUp(connection)
		if isContinue == false {
			break
		}
	}
	if response.IsSent() == true {
		return
	}

	handler := request.Route().Handler()
	a.container.Inject(handler)
	if h, ok := handler.(ContainerAware); ok == true {
		h.WithContainer(a.container)
	}
	result, err := handler.Handle(connection)
	if err != nil {
		a.handleError(connection, err)
		return
	}
	for _, hook := range request.Route().Hooks() {
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

func (a *FactoryApp) sendBadRequestResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func (a *FactoryApp) handleError(connection Connection, err error) {
	for _, loader := range a.loaders {
		loader.TearDown(a, err)
	}

	if connection.Response().IsSent() == true {
		return
	}

	h := a.errorHandler
	if h == nil {
		h = NewErrorHandler()
	}
	PanicOnError(h.HandleError(connection, err))
}
