package lapi

import (
	"fmt"
	"reflect"
)

const (
	BindErrorInvalidInterface        = 1
	BindErrorInvalidConcrete         = 2
	BindErrorNotImplementedInterface = 3
)

// Container acts as a dependency-injection manager
type Container interface {
	Binder
	Resolver
	Injector
}

// Binder uses to bind a concrete to an abstract
type Binder interface {
	// Bind stores a concrete of an abstract, as default sharing is enable
	Bind(abstract interface{}, concrete interface{}) SystemError
}

// Resolver helps to resolve dependencies
type Resolver interface {
	// Resolve processes and returns a concrete of proposed abstract
	Resolve(abstract interface{}) (concrete interface{}, err SystemError)
}

// Injector works as a tool to inject dependencies
type Injector interface {
	// Inject resolves target's dependencies
	Inject(target interface{}) SystemError
}

func NewContainer() Container {
	return &FactoryContainer{make(map[reflect.Type]interface{})}
}

type FactoryContainer struct {
	items map[reflect.Type]interface{}
}

func (c *FactoryContainer) Bind(abstract interface{}, concrete interface{}) SystemError {
	at, err := c.interfaceOf(abstract)
	if err != nil {
		return err
	}

	ct := reflect.TypeOf(concrete)
	switch ct.Kind() {
	case reflect.Func:
	case reflect.Struct, reflect.Ptr:
		if c.instanceOf(at, ct) == false {
			return NewSystemError(BindErrorNotImplementedInterface, fmt.Sprintf("%v is not an instance of %v", ct, at))
		}
	default:
		return NewSystemError(BindErrorInvalidConcrete, fmt.Sprintf("Non-supported kind of concrete. Got %v", ct.Kind()))
	}
	return nil
}

func (c *FactoryContainer) Resolve(abstract interface{}) (concrete interface{}, err SystemError) {
	return nil, nil
}

func (c *FactoryContainer) Inject(target interface{}) SystemError {
	return nil
}

func (c *FactoryContainer) interfaceOf(value interface{}) (reflect.Type, SystemError) {
	t := reflect.TypeOf(value)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		return nil, NewSystemError(BindErrorInvalidInterface, "Called interfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)")
	}

	return t, nil
}

func (c *FactoryContainer) instanceOf(a reflect.Type, b reflect.Type) bool {
	if b.Kind() != reflect.Struct || b.Kind() != reflect.Ptr || a.Kind() != reflect.Interface {
		return false
	}

	return b.Implements(a)
}
