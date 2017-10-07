package lapi

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/goline/errors"
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
	Bind(abstract interface{}, concrete interface{}) errors.Error
}

// Resolver helps to resolve dependencies
type Resolver interface {
	// Resolve processes and returns a concrete of proposed abstract
	Resolve(abstract interface{}, args ...interface{}) (concrete interface{}, err errors.Error)
}

// Injector works as a tool to inject dependencies
type Injector interface {
	// Inject resolves target's dependencies
	Inject(target interface{}) errors.Error
}

// ContainerAware handles a container
type ContainerAware interface {
	// Container returns an instance of Container
	Container() Container

	// WithContainer allows to set container
	WithContainer(container Container) ContainerAware
}

func NewContainer() Container {
	return &FactoryContainer{items: new(sync.Map)}
}

type FactoryContainer struct {
	items *sync.Map
}

func (c *FactoryContainer) Bind(abstract interface{}, concrete interface{}) errors.Error {
	at, isInterface := c.interfaceOf(abstract)
	if isInterface == nil {
		return c.bindInterface(at, concrete)
	}

	at, isStruct := c.structOf(abstract)
	if isStruct == nil {
		return c.bindStruct(at, concrete)
	}

	return errors.New(ERR_BIND_INVALID_ARGUMENTS, "Binding error! Invalid arguments.")
}

func (c *FactoryContainer) Resolve(abstract interface{}, args ...interface{}) (concrete interface{}, err errors.Error) {
	at, isInterface := c.interfaceOf(abstract)
	if isInterface == nil {
		return c.resolveInterface(at, args...)
	}

	at, isStruct := c.structOf(abstract)
	if isStruct == nil {
		return c.resolveStruct(at, args...)
	}

	return nil, errors.New(ERR_RESOLVE_INVALID_ARGUMENTS, "Resolving error! Invalid arguments.")
}

func (c *FactoryContainer) Inject(target interface{}) errors.Error {
	t := reflect.TypeOf(target)
	switch t.Kind() {
	case reflect.Ptr:
	default:
		return errors.New(ERR_INJECT_INVALID_TARGET_TYPE, fmt.Sprintf("Injecting to %v is not supported", t.Kind()))
	}

	s := t.Elem()
	n := s.NumField()
	if n == 0 {
		return nil
	}
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < n; i++ {
		sf := s.Field(i)
		if _, ok := sf.Tag.Lookup("inject"); ok == false {
			continue
		}

		if sf.Type.Kind() != reflect.Interface &&
			sf.Type.Kind() != reflect.Struct && sf.Type.Kind() != reflect.Ptr {
			continue
		}

		f := v.Field(i)
		if f.CanSet() == false {
			continue
		}

		o, err := c.Resolve(sf.Type)
		if err != nil {
			return err
		}
		c.Inject(o)
		f.Set(reflect.ValueOf(o))
	}
	return nil
}

func (c *FactoryContainer) bindInterface(at reflect.Type, concrete interface{}) errors.Error {
	ct := reflect.TypeOf(concrete)
	switch ct.Kind() {
	case reflect.Func:
	case reflect.Ptr:
		if c.instanceOf(at, ct) == false {
			return errors.New(ERR_BIND_NOT_IMPLEMENT_INTERFACE, fmt.Sprintf("%v is not an instance of %v", ct, at))
		}
	default:
		return errors.New(ERR_BIND_INVALID_CONCRETE, fmt.Sprintf("Non-supported kind of concrete. Got %v", ct.Kind()))
	}

	c.items.Store(at.String(), reflect.ValueOf(concrete))
	return nil
}

func (c *FactoryContainer) bindStruct(at reflect.Type, concrete interface{}) errors.Error {
	ct, err := c.structOf(concrete)
	if err != nil {
		return err
	}

	if at.String() != ct.String() {
		return errors.New(ERR_BIND_INVALID_STRUCT_CONCRETE, fmt.Sprintf("Expects %s. Got %s", at.String(), ct.String()))
	}

	c.items.Store(at.String(), reflect.ValueOf(concrete))
	return nil
}

func (c *FactoryContainer) resolveInterface(at reflect.Type, args ...interface{}) (concrete interface{}, err errors.Error) {
	v, ok := c.items.Load(at.String())
	if ok == false {
		return nil, errors.New(ERR_RESOLVE_NOT_EXIST_ABSTRACT, fmt.Sprintf("%v is not bound yet", at))
	}
	value := v.(reflect.Value)

	switch value.Kind() {
	case reflect.Func:
		return c.resolveFunc(value, args...)
	case reflect.Ptr:
		return value.Interface(), nil
	default:
		return nil, errors.New(ERR_RESOLVE_INVALID_CONCRETE, fmt.Sprintf("Type %v is not supported", value.Kind()))
	}
}

func (c *FactoryContainer) resolveStruct(at reflect.Type, args ...interface{}) (concrete interface{}, err errors.Error) {
	v, ok := c.items.Load(at.String())
	if ok == false {
		return nil, errors.New(ERR_RESOLVE_NOT_EXIST_ABSTRACT, fmt.Sprintf("%v is not bound yet", at))
	}
	value := v.(reflect.Value)

	switch value.Kind() {
	case reflect.Struct, reflect.Ptr:
		return value.Interface(), nil
	default:
		return nil, errors.New(ERR_RESOLVE_INVALID_CONCRETE, fmt.Sprintf("Type %v is not supported", value.Kind()))
	}
}

func (c *FactoryContainer) structOf(value interface{}) (reflect.Type, errors.Error) {
	if t, ok := value.(reflect.Type); ok == true {
		return c.structOfType(t)
	}

	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New(ERR_BIND_INVALID_STRUCT, "Called structOf with a value that is not a pointer to a struct. (*MyStruct)(nil)")
	}

	return t, nil
}

func (c *FactoryContainer) structOfType(t reflect.Type) (reflect.Type, errors.Error) {
	switch t.Kind() {
	case reflect.Struct:
		return t, nil
	case reflect.Ptr:
		return t.Elem(), nil
	default:
		return nil, errors.New(ERR_BIND_INVALID_STRUCT, "Called structOfType with a value that is not a pointer to a struct. (*MyStruct)(nil)")
	}
}

func (c *FactoryContainer) interfaceOf(value interface{}) (reflect.Type, errors.Error) {
	if t, ok := value.(reflect.Type); ok == true && t.Kind() == reflect.Interface {
		return t, nil
	}
	t := reflect.TypeOf(value)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		return nil, errors.New(ERR_BIND_INVALID_INTERFACE, "Called interfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)")
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

func (c *FactoryContainer) resolveFunc(value reflect.Value, args ...interface{}) (concrete interface{}, err errors.Error) {
	t := value.Type()
	if len(args) != t.NumIn() {
		return nil, errors.New(ERR_RESOLVE_INSUFFICIENT_ARGUMENTS, fmt.Sprintf("Expects to have %v input arguments. Got %v", t.NumIn(), len(args)))
	}

	in := make([]reflect.Value, t.NumIn())
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	out := value.Call(in)
	if len(out) == 0 {
		return nil, errors.New(ERR_RESOLVE_NON_VALUES_RETURNED, "Expects to have at least 1 value returned. Got 0")
	}
	return out[0].Interface(), nil
}
