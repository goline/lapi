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
	Config() Bag

	// WithConfig sets application's config
	WithConfig(config Bag) App
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
		loaders:   make(map[int]*Slice),
		router:    NewRouter(),
		container: NewContainer(),
		rescuer:   NewRescuer(),
		config:    NewBag(),
	}
}

type FactoryApp struct {
	config    Bag
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

func (a *FactoryApp) Config() Bag {
	return a.config
}

func (a *FactoryApp) WithConfig(config Bag) App {
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
	if a.container == nil {
		panic(errors.New(ERR_CONTAINER_NOT_DEFINED, "App requires a container to run"))
	}

	if a.rescuer == nil {
		panic(errors.New(ERR_RESCUER_NOT_DEFINED, "App requires a rescuer to be defined"))
	}

	if a.router == nil {
		panic(errors.New(ERR_ROUTER_NOT_DEFINED, fmt.Sprint("Router is not defined yet.")))
	}

	Parallel(a.loaders, func(l interface{}) {
		l.(Loader).Load(a)
	})
}

func (a *FactoryApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection := a.setUpConnection(w, r)
	defer a.forceSendResponse(connection)
	defer a.forceRecover(connection)

	PanicOnError(a.router.Route(connection.Request()))
	Parallel(connection.Request().Route().Hooks(), func(item interface{}) {
		if hook, ok := item.(BootableHook); ok == true {
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
		if hook, ok := item.(HaltableHook); ok == true {
			defer a.forceRecover(connection)
			PanicOnError(hook.TearDown(connection, result, err))
		}
	})
}

func (a *FactoryApp) forceSendResponse(connection Connection) {
	if connection.Response().IsSent() == false {
		connection.Response().Send()
	}
}

func (a *FactoryApp) forceRecover(connection Connection) {
	if r := recover(); r != nil {
		PanicOnError(a.rescuer.Rescue(connection, r))

		// After handling error, we must send response out
		// However, we don't use defer here, as we don't want
		// to send any response if error is not handled well
		a.forceSendResponse(connection)
	}
}

func (a *FactoryApp) setUpConnection(w http.ResponseWriter, r *http.Request) Connection {
	request := NewRequest(r)
	response := NewJsonResponse(w)

	return NewConnection(request, response)
}
