package lapi

import (
	"fmt"
	"reflect"

	"github.com/goline/errors"
)

type Async struct {
	runners []asyncRunner
}
type asyncRunner struct {
	runner    interface{}
	arguments []interface{}
}

func (a *Async) Add(runner interface{}, arguments ...interface{}) errors.Error {
	if t := reflect.TypeOf(runner); t.Kind() != reflect.Func {
		return errors.New(ERR_ASYNC_INVALID_TYPE, fmt.Sprintf("runner must be a function. Got %s", t.Kind()))
	}

	a.runners = append(a.runners, asyncRunner{runner: runner, arguments: arguments})
	return nil
}

func (a *Async) Run() errors.Error {
	n := len(a.runners)
	fin := 1
	ec := make(chan errors.Error, 1)
	for _, runner := range a.runners {
		go func(runner asyncRunner) {
			in := make([]reflect.Value, len(runner.arguments))
			for i, argument := range runner.arguments {
				in[i] = reflect.ValueOf(argument)
			}

			v := reflect.ValueOf(runner.runner).Call(in)[0]
			if v.IsNil() {
				fin++
				if fin == n {
					ec <- nil
				}
			} else {
				ec <- v.Interface().(errors.Error)
			}
		}(runner)
	}

	return <-ec
}

func (a *Async) Reset() {
	a.runners = make([]asyncRunner, 0)
}
