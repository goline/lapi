package lapi

import (
	"fmt"
	"reflect"
)

const (
	BindErrorInvalidInterface         = 1
	BindErrorInvalidConcrete          = 2
	BindErrorNotImplementedInterface  = 3
	ResolveErrorNotExistAbstract      = 4
	ResolveErrorInvalidConcrete       = 5
	ResolveErrorInsufficientArguments = 6
	ResolveErrorNonValuesReturned     = 7
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
	Resolve(abstract interface{}, args ...interface{}) (concrete interface{}, err SystemError)
}

// Injector works as a tool to inject dependencies
type Injector interface {
	// Inject resolves target's dependencies
	Inject(target interface{}) SystemError
}

func NewContainer() Container {
	return &FactoryContainer{make(map[reflect.Type]reflect.Value)}
}

type FactoryContainer struct {
	items map[reflect.Type]reflect.Value
}

func (c *FactoryContainer) Bind(abstract interface{}, concrete interface{}) SystemError {
	at, err := c.interfaceOf(abstract)
	if err != nil {
		return err
	}

	ct := reflect.TypeOf(concrete)
	switch ct.Kind() {
	case reflect.Func:
	case reflect.Ptr:
		if c.instanceOf(at, ct) == false {
			return NewSystemError(BindErrorNotImplementedInterface, fmt.Sprintf("%v is not an instance of %v", ct, at))
		}
	default:
		return NewSystemError(BindErrorInvalidConcrete, fmt.Sprintf("Non-supported kind of concrete. Got %v", ct.Kind()))
	}

	cv := reflect.ValueOf(concrete)
	c.items[at] = cv
	return nil
}

func (c *FactoryContainer) Resolve(abstract interface{}, args ...interface{}) (concrete interface{}, err SystemError) {
	at, err := c.interfaceOf(abstract)
	if err != nil {
		return nil, err
	}

	value, ok := c.items[at]
	if ok == false {
		return nil, NewSystemError(ResolveErrorNotExistAbstract, fmt.Sprintf("%v is not bound yet", at))
	}

	switch value.Kind() {
	case reflect.Func:
		return c.resolveFunc(value, args...)
	case reflect.Ptr:
		return value.Interface(), nil
	default:
		return nil, NewSystemError(ResolveErrorInvalidConcrete, fmt.Sprintf("Type %v is not supported", value.Kind()))
	}
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

func (c *FactoryContainer) instanceOf(abstract reflect.Type, concrete reflect.Type) bool {
	if abstract.Kind() != reflect.Interface {
		return false
	}

	switch concrete.Kind() {
	case reflect.Struct, reflect.Ptr:
		return concrete.Implements(abstract)
	default:
		return false
	}
}

func (c *FactoryContainer) resolveFunc(value reflect.Value, args ...interface{}) (concrete interface{}, err SystemError) {
	t := value.Type()
	if len(args) != t.NumIn() {
		return nil, NewSystemError(ResolveErrorInsufficientArguments, fmt.Sprintf("Expects to have %v input arguments. Got %v", t.NumIn(), len(args)))
	}

	in := make([]reflect.Value, t.NumIn())
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	out := value.Call(in)
	if len(out) == 0 {
		return nil, NewSystemError(ResolveErrorNonValuesReturned, fmt.Sprintf("Expects to have at least 1 value returned. Got 0"))
	}
	return out[0].Interface(), nil
}
