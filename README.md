# lapi
Light API for Golang

[![Build Status](https://api.travis-ci.org/goline/lapi.svg)](https://travis-ci.org/goline/lapi)

## Concepts

In LAPI, we have 2 most important concepts, and they are *Loader*, *Hook*
- Loader: loads when application starts
- Hook: processes request after a route is found

```text
Loader and Hook are loaded concurrently; therefore, you should set it priority in order to make it run as expected order.
```

## Quick Start

```go
package main

import (
	. "github.com/goline/lapi"
	"github.com/goline/validation"
	"github.com/goline/errors"
	"github.com/goline/log"
	
	"net/http"
)

type UsersHandler Handler
type UsersPostHandler Handler

func main() {
	app := NewApp()
    app.WithLoader(NewLoader(func(a App) {
    		c := a.Container()
    
    		Must(
    			c.Bind((*log.Logger)(nil), new(log.ConsoleLogger)),
    			c.Bind((*validation.Validator)(nil), validation.New()),
    		)
    	}, PRIORITY_DEFAULT)).
    	WithLoader(NewLoader(func(a App) {
            r := app.Router()
            r.Get("/users", new(UsersHandler))
            r.Post("/users/<id:\\d+>", new(UsersPostHandler))
        }, PRIORITY_DEFAULT)).
    	WithLoader(NewLoader(func(a App) {
    		app.Router().
    			WithHook(new(SystemHook))
    	}, 100))	
}
```