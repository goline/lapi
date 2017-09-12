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
	AppRescuer
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

// AppRescuer manages error handler
type AppRescuer interface {
	// Rescuer returns an instance of Rescuer
	Rescuer() Rescuer

	// WithRescuer sets error handler
	WithRescuer(handler Rescuer) App
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
	config    Config
	container Container
	loaders   []Loader
	request   Request
	response  Response
	router    Router
	rescuer   Rescuer
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

func (a *FactoryApp) Rescuer() Rescuer {
	if a.rescuer == nil {
		a.rescuer = NewRescuer()
	}

	return a.rescuer
}

func (a *FactoryApp) WithRescuer(handler Rescuer) App {
	a.rescuer = handler
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
	PanicOnError(a.container.Inject(a.rescuer))
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

	PanicOnError(a.router.Route(connection.Request()))
	for _, hook := range connection.Request().Route().Hooks() {
		PanicOnError(hook.SetUp(connection))
	}
	if connection.Response().IsSent() == true {
		return
	}

	handler := connection.Request().Route().Handler()
	if handler == nil {
		panic(NewSystemError(ERROR_NO_HANDLER_FOUND, "No handler found"))
	}

	PanicOnError(a.container.Inject(handler))
	if h, ok := handler.(ContainerAware); ok == true {
		h.WithContainer(a.container)
	}
	result, err := handler.Handle(connection)
	for _, hook := range connection.Request().Route().Hooks() {
		PanicOnError(hook.TearDown(connection, result, err))
	}
}

func (a *FactoryApp) forceSendResponse(connection Connection) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok == true {
			a.rescuer.Rescue(connection, err)
		}
	}
	if connection.Response().IsSent() == false {
		connection.Response().Send()
	}
}

func (a *FactoryApp) setUpConnection(w http.ResponseWriter, r *http.Request) Connection {
	request := NewRequest(r)
	response := NewResponse(w)

	// let make request and response have same content-type and charset as default
	response.WithContentType(request.ContentType()).WithCharset(request.Charset())

	return NewConnection(request, response)
}
