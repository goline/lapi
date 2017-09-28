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
	if c, ok := app.Config().(ServerConfig); ok == true {
		PanicOnError(http.ListenAndServe(c.Address(), nil))
	} else {
		panic(errors.New(ERR_SERVER_CONFIG_MISSING, fmt.Sprint("Server configuration is missing")))
	}
}
