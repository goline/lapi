package lapi

import (
	"fmt"
	"net/http"

	"github.com/goline/errors"
)

// Loader is an application loader which could be useful for set things up
type Loader interface {
	// Load runs when application is starting up
	Load(app App)
}

type ServerLoader struct {
	PriorityAware
}

func (l *ServerLoader) Load(app App) {
	PanicOnError(app.Container().Inject(app.Rescuer()))

	http.Handle("/", app)
	if address, ok := app.Config().GetString("server.address"); ok {
		PanicOnError(http.ListenAndServe(address, nil))
	} else {
		panic(errors.New(ERR_SERVER_CONFIG_MISSING, fmt.Sprint("Server configuration is missing")))
	}
}

func NewLoader(runner func(app App), priority int) *ServiceLoader {
	l := new(ServiceLoader)
	l.runner = runner
	l.WithPriority(priority)
	return l
}

type ServiceLoader struct {
	PriorityAware
	runner func(app App)
}

func (l *ServiceLoader) Load(app App) {
	l.runner(app)
}
