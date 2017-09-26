package lapi

import (
	"fmt"
	"net/http"

	"github.com/goline/errors"
)

// App is a central application
type App interface {
	AppLoader
	AppRunner
	AppRouter
	AppRescuer
	AppConfigger
	ContainerAware
	http.Handler
}

// AppLoader handles application's loader
type AppLoader interface {
	// WithLoader allows to register application's loader
	WithLoader(loader Loader) App
}

type AppConfigger interface {
	// Config returns application's configuration
	Config() interface{}

	// WithConfig sets application's config
	WithConfig(config interface{}) App
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

type AppDryRunner interface {
	// DryRun brings application up
	// It should configure and load only
	DryRun()
}

func NewApp() App {
	return &FactoryApp{
		loaders:   make(map[int]*Slice),
		router:    NewRouter(),
		container: NewContainer(),
		rescuer:   NewRescuer(),
	}
}

type FactoryApp struct {
	config    interface{}
	container Container
	loaders   map[int]*Slice
	router    Router
	rescuer   Rescuer
}

func (a *FactoryApp) WithLoader(loader Loader) App {
	p := PRIORITY_DEFAULT
	if l, ok := loader.(Prioritizer); ok == true {
		p = l.Priority()
	}

	if a.loaders[p] == nil {
		a.loaders[p] = new(Slice)
	}

	a.loaders[p].Append(loader)
	return a
}

func (a *FactoryApp) Config() interface{} {
	return a.config
}

func (a *FactoryApp) WithConfig(config interface{}) App {
	a.config = config
	return a
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

func (a *FactoryApp) Container() Container {
	return a.container
}

func (a *FactoryApp) WithContainer(container Container) ContainerAware {
	a.container = container
	return a
}

func (a *FactoryApp) Run() {
	a.SetUp().Handle()
}

func (a *FactoryApp) DryRun() {
	a.SetUp()
}

func (a *FactoryApp) SetUp() *FactoryApp {
	if a.container == nil {
		panic(errors.New(ERR_CONTAINER_NOT_DEFINED, "App requires a container to run"))
	}

	Parallel(a.loaders, func(l interface{}) {
		l.(Loader).Load(a)
	})

	if a.rescuer == nil {
		panic(errors.New(ERR_RESCUER_NOT_DEFINED, "App requires a rescuer to be defined"))
	}
	PanicOnError(a.container.Inject(a.rescuer))
	return a
}

func (a *FactoryApp) Handle() *FactoryApp {
	if a.router == nil {
		panic(errors.New(ERR_ROUTER_NOT_DEFINED, fmt.Sprint("Router is not defined yet.")))
	}

	http.Handle("/", a)
	if c, ok := a.config.(ServerConfig); ok == true {
		PanicOnError(http.ListenAndServe(c.Address(), nil))
	} else {
		panic(errors.New(ERR_SERVER_CONFIG_MISSING, fmt.Sprint("Server configuration is missing")))
	}
	return a
}

func (a *FactoryApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection := a.setUpConnection(w, r)
	defer a.forceSendResponse(connection)

	PanicOnError(a.router.Route(connection.Request()))
	Parallel(connection.Request().Route().Hooks(), func(item interface{}) {
		if hook, ok := item.(Hook); ok == true {
			defer a.forceRecover(connection)
			PanicOnError(hook.SetUp(connection))
		}
	})
	if connection.Response().IsSent() == true {
		return
	}

	handler := connection.Request().Route().Handler()
	if handler == nil {
		panic(errors.New(ERR_NO_HANDLER_FOUND, "No handler found"))
	}

	PanicOnError(a.container.Inject(handler))
	if h, ok := handler.(ContainerAware); ok == true {
		h.WithContainer(a.container)
	}
	result, err := handler.Handle(connection)
	Parallel(connection.Request().Route().Hooks(), func(item interface{}) {
		if hook, ok := item.(Hook); ok == true {
			defer a.forceRecover(connection)
			PanicOnError(hook.TearDown(connection, result, err))
		}
	})
}

func (a *FactoryApp) forceSendResponse(connection Connection) {
	a.forceRecover(connection)
	if connection.Response().IsSent() == false {
		connection.Response().Send()
	}
}

func (a *FactoryApp) forceRecover(connection Connection) {
	if r := recover(); r != nil {
		PanicOnError(a.rescuer.Rescue(connection, r))
	}
}

func (a *FactoryApp) setUpConnection(w http.ResponseWriter, r *http.Request) Connection {
	request := NewRequest(r)
	response := NewResponse(w)

	// let make request and response have same content-type and charset as default
	response.WithContentType(request.ContentType()).WithCharset(request.Charset())

	return NewConnection(request, response)
}
