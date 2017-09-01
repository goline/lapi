package lapi

import (
	"errors"
	"fmt"
	"net/http"
)

// App is a central application
type App interface {
	ContainerAware
	AppLoader
	AppRunner
	AppConfigger
	AppConnector
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

type AppConnector interface {
	// Request returns an instance of Request
	Request() Request

	// Response returns an instance of Response
	Response() Response
}

// AppRouter handles router
type AppRouter interface {
	// Router returns an instance of Router
	Router() Router

	// WithRouter sets router
	WithRouter(router Router) App
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
	a.request = NewRequest(r)
	a.response = NewResponse(w)
	defer a.forceSendResponse()

	if a.router == nil {
		a.handleError(NewSystemError(ERROR_ROUTER_NOT_DEFINED, fmt.Sprint("Router is not defined yet.")))
	}

	err := a.router.Route(a.request)
	if err != nil {
		a.handleError(err)
		return
	}

	for _, hook := range a.request.Route().Hooks() {
		isContinue := hook.SetUp(a.request, a.response)
		if isContinue == false {
			break
		}
	}
	if a.response.IsSent() == true {
		return
	}

	result, err := a.request.Route().Handler().Handle(a.request, a.response)
	for _, hook := range a.request.Route().Hooks() {
		isContinue := hook.TearDown(a.request, a.response, result, err)
		if isContinue == false {
			break
		}
	}
}

func (a *FactoryApp) forceSendResponse() {
	if a.response.IsSent() == false {
		a.response.Send()
	}
}

func (a *FactoryApp) handleError(err error) {
	for _, loader := range a.loaders {
		loader.TearDown(a, err)
	}

	if a.response.IsSent() == false {
		a.modifyResponseOnError(err)
	}
}

func (a *FactoryApp) modifyResponseOnError(err error) {
	if e, ok := err.(ErrorStatus); ok == true {
		a.response.WithStatus(e.Status())
	} else {
		a.response.WithStatus(http.StatusInternalServerError)
	}

	var es []Error
	switch e := err.(type) {
	case Error:
		es = make([]Error, 1)
		es[0] = e
	case StackError:
		es = e.Errors()
	default:
		es[0] = NewError("", "ERROR_HANDLE_INVALID_ERROR", errors.New("Error's type is not supported."))
	}

	ei := make([]errorItemResponse, len(es))
	for i, e := range es {
		ei[i] = errorItemResponse{e.Code(), e.Message()}
	}
	er := &errorStackResponse{ei}
	a.response.WithContent(er)
}

type errorStackResponse struct {
	Errors []errorItemResponse `json:"errors"`
}

type errorItemResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
