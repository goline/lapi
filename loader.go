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

type HttpServerLoader struct{}

func (l *HttpServerLoader) Load(app App) {
	if app.Router() == nil {
		panic(errors.New(ERR_ROUTER_NOT_DEFINED, fmt.Sprint("Router is not defined yet.")))
	}

	http.Handle("/", app)
	if c, ok := app.Config().(ServerConfig); ok == true {
		PanicOnError(http.ListenAndServe(c.Address(), nil))
	} else {
		panic(errors.New(ERR_SERVER_CONFIG_MISSING, fmt.Sprint("Server configuration is missing")))
	}
}
