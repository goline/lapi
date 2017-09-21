package lapi

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	if _, ok := c.(Container); ok == false {
		t.Errorf("Expects an instance of Container. Got %+v", c)
	}
}

func TestFactoryContainer_Bind_ErrorAbstractIsNotAnInterfaceOrStruct(t *testing.T) {
	c := NewContainer()
	err := c.Bind("string", &FactoryError{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	} else if e, ok := err.(Error); ok == false || e.Code() != ERR_BIND_INVALID_ARGUMENTS {
		t.Errorf("Expects ERR_BIND_INVALID_ARGUMENTS. Got %s", e.Code())
	}
}

func TestFactoryContainer_Bind_ErrorConcreteNotImplementAbstract(t *testing.T) {
	c := NewContainer()
	err := c.Bind((*Bag)(nil), &FactoryError{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
	if e, ok := err.(Error); ok == false || e.Code() != ERR_BIND_NOT_IMPLEMENT_INTERFACE {
		t.Errorf("Expects ERR_BIND_NOT_IMPLEMENT_INTERFACE code. Got %s", e.Code())
	}
}

func TestFactoryContainer_Bind_ErrorNotSupportAbstractType(t *testing.T) {
	c := NewContainer()
	a := make([]interface{}, 2)
	a[0] = "a_string"
	a[1] = FactoryError{"my_code", "my_message", errors.New("my_err")}
	for _, i := range a {
		err := c.Bind((*Bag)(nil), i)
		if err == nil {
			t.Errorf("Expects err is not nil")
		}
		if e, ok := err.(Error); ok == false || e.Code() != ERR_BIND_INVALID_CONCRETE {
			t.Errorf("Expects ERR_BIND_INVALID_CONCRETE code. Got %v", e.Code())
		}
	}
}

func TestFactoryContainer_Bind_StructOk(t *testing.T) {
	c := NewContainer()
	err := c.Bind(FactoryBag{}, &FactoryBag{})
	if err != nil {
		t.Errorf("Expects err is nil. Got %v", err)
	}
}

func TestFactoryContainer_Bind_Struct_BindErrorInvalidStruct(t *testing.T) {
	c := NewContainer()
	err := c.Bind(FactoryBag{}, "not_a_struct")
	if err == nil {
		t.Errorf("Expects err is not nil")
	} else if e, ok := err.(Error); ok == false || e.Code() != ERR_BIND_INVALID_STRUCT {
		t.Errorf("Expects ERR_BIND_INVALID_STRUCT. Got %s", e.Code())
	}
}

func TestFactoryContainer_Bind_Struct_BindErrorInvalidStructConcrete(t *testing.T) {
	c := NewContainer()
	err := c.Bind(FactoryBag{}, &FactoryError{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	} else if e, ok := err.(Error); ok == false || e.Code() != ERR_BIND_INVALID_STRUCT_CONCRETE {
		t.Errorf("Expects ERR_BIND_INVALID_STRUCT_CONCRETE. Got %s", e.Code())
	}
}

func TestFactoryContainer_Bind_Struct_StructOf(t *testing.T) {
	c := &FactoryContainer{}
	_, err := c.structOf(reflect.TypeOf(FactoryBag{}))
	if err != nil {
		t.Errorf("Expects err is nil")
	}
}

func TestFactoryContainer_Resolve_ErrorAbstractNotAnInterface(t *testing.T) {
	c := &FactoryContainer{make(map[string]reflect.Value)}
	_, err := c.Resolve(&FactoryBag{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
}

func TestFactoryContainer_Resolve_ErrorInvalidConcreteType(t *testing.T) {
	c := &FactoryContainer{make(map[string]reflect.Value)}
	c.Bind((*Bag)(nil), NewBag())
	for k := range c.items {
		c.items[k] = reflect.ValueOf("a_string")
		break
	}
	_, err := c.Resolve((*Bag)(nil))
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
}

func TestFactoryContainer_Bind_Function_Ok(t *testing.T) {
	c := NewContainer()
	err := c.Bind((*Error)(nil), func(code string, message string, err error) Error {
		return NewError(code, message, err)
	})
	if err != nil {
		t.Errorf("Expects error is not nil")
	}
}

func TestFactoryContainer_Bind_Ptr_Ok(t *testing.T) {
	c := NewContainer()
	err := c.Bind((*Error)(nil), NewError("my_code", "my_message", errors.New("my_error")))
	if err != nil {
		t.Errorf("Expects error is not nil")
	}
}

func TestFactoryContainer_Resolve_Func_ErrorInsufficientArguments(t *testing.T) {
	c := NewContainer()
	c.Bind((*Error)(nil), func(code string, message string, err error) Error {
		return NewError(code, message, err)
	})
	_, err := c.Resolve((*Error)(nil), "my_code", "my_message")
	if err == nil {
		t.Errorf("Expects error is not nil")
	}
	if e, ok := err.(Error); ok == false || e.Code() != ERR_RESOLVE_INSUFFICIENT_ARGUMENTS {
		t.Errorf("Expects ERR_RESOLVE_INSUFFICIENT_ARGUMENTS code. Got %s", e.Code())
	}
}

func TestFactoryContainer_Resolve_Func_ErrorNonValuesReturned(t *testing.T) {
	c := NewContainer()
	c.Bind((*Error)(nil), func(code string, message string, err error) {})
	_, err := c.Resolve((*Error)(nil), "my_code", "my_message", errors.New("my_err"))
	if err == nil {
		t.Errorf("Expects error is not nil")
	}
	if e, ok := err.(Error); ok == false || e.Code() != ERR_RESOLVE_NON_VALUES_RETURNED {
		t.Errorf("Expects ERR_RESOLVE_NON_VALUES_RETURNED code. Got %s", e.Code())
	}
}

func TestFactoryContainer_Resolve_Func_Ok(t *testing.T) {
	c := NewContainer()
	c.Bind((*Error)(nil), func(code string, message string, err error) Error {
		return NewError(code, message, err)
	})
	e, err := c.Resolve((*Error)(nil), "my_code", "my_message", errors.New("my_err"))
	if err != nil {
		t.Errorf("Expects error is nil")
	}
	er, ok := e.(Error)
	if ok == false {
		t.Errorf("Expects error is an instance of Error")
	}
	if er.Code() != "my_code" {
		t.Errorf("Expects error's code is my_code")
	}
}

func TestFactoryContainer_Resolve_Ptr_Ok(t *testing.T) {
	c := NewContainer()
	c.Bind((*Error)(nil), &FactoryError{"my_code", "my_message", errors.New("my_err")})
	e, err := c.Resolve((*Error)(nil))
	if err != nil {
		t.Errorf("Expects error is nil")
	}
	er, ok := e.(Error)
	if ok == false {
		t.Errorf("Expects error is an instance of Error")
	}
	if er.Code() != "my_code" {
		t.Errorf("Expects error's code is my_code")
	}
}

func TestFactoryContainer_Resolve_Struct(t *testing.T) {
	c := NewContainer()
	c.Bind(FactoryBag{}, NewBag())
	v, err := c.Resolve(FactoryBag{})
	if err != nil {
		t.Errorf("Expects err is nil. Got %v", err)
	}
	if v == nil {
		t.Errorf("Expects v is not nil")
	}
}

func TestFactoryContainer_Resolve_ResolveErrorInvalidArguments(t *testing.T) {
	c := NewContainer()
	_, err := c.Resolve("string")
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
}

func TestFactoryContainer_Inject_ErrorInjectInvalidTargetType(t *testing.T) {
	c := NewContainer()
	e := FactoryError{"my_code", "my_message", errors.New("my_err")}
	err := c.Inject(e)
	if err == nil {
		t.Errorf("Expects error is not nil")
	}
	if e, ok := err.(Error); ok == false || e.Code() != ERR_INJECT_INVALID_TARGET_TYPE {
		t.Errorf("Expects ERR_INJECT_INVALID_TARGET_TYPE code. Got %s", e.Code())
	}
}

type InjectErrorResolveErrorNotExistAbstract struct {
	Err Error `inject:"*"`
}

func TestFactoryContainer_Inject_ErrorResolveErrorNotExistAbstract(t *testing.T) {
	c := NewContainer()
	in := &InjectErrorResolveErrorNotExistAbstract{}
	err := c.Inject(in)
	if err == nil {
		t.Errorf("Expects error is not nil")
	}
	if e, ok := err.(Error); ok == false || e.Code() != ERR_RESOLVE_NOT_EXIST_ABSTRACT {
		t.Errorf("Expects ERR_RESOLVE_NOT_EXIST_ABSTRACT code. Got %s", e.Code())
	}
}

type InjectOk struct {
	Err Error `inject:"*"`
}

func TestFactoryContainer_Inject_Ok(t *testing.T) {
	c := NewContainer()
	c.Bind((*Error)(nil), NewError("my_code", "my_message", errors.New("my_err")))
	in := &InjectOk{}
	err := c.Inject(in)
	if err != nil {
		t.Errorf("Expects error is nil")
	}
	if in.Err.Code() != "my_code" {
		t.Errorf("Expects Err.Code() is my_code. Got %v", in.Err.Code())
	}
	if in.Err.Message() != "my_message" {
		t.Errorf("Expects Err.Message() is my_message. Got %v", in.Err.Message())
	}
	if in.Err.Trace().Error() != "my_err" {
		t.Errorf("Expects Err.Error() is my_err. Got %v", in.Err.Trace().Error())
	}
}

type InjectRecursiveOk struct {
	Err                          Error       `inject:"*"`
	Foo                          InjectFooer `inject:"*"`
	NotInjectableNonInterface    string      `inject:"*"`
	notInjectablePrivateProperty Error       `inject:"*"`
}
type InjectFooer interface {
	Foo() string
}
type InjectFoo struct {
	Baz InjectBazer `inject:"*"`
}

func (f *InjectFoo) Foo() string { return "Foo.." }

type InjectBazer interface {
	Baz() string
}
type InjectBaz struct{}

func (b *InjectBaz) Baz() string { return "Baz.." }
func TestFactoryContainer_Inject_Recursive_Ok(t *testing.T) {
	// InjectRecursiveOk -> Err
	// InjectRecursiveOk -> Foo -> Baz
	c := NewContainer()
	c.Bind((*Error)(nil), NewError("my_code", "my_message", errors.New("my_err")))
	c.Bind((*InjectFooer)(nil), &InjectFoo{})
	c.Bind((*InjectBazer)(nil), &InjectBaz{})
	in := &InjectRecursiveOk{}
	err := c.Inject(in)
	if err != nil {
		t.Errorf("Expects error is nil")
	}

	// Asserting for in.Err
	if in.Err.Code() != "my_code" {
		t.Errorf("Expects Err.Code() is my_code. Got %v", in.Err.Code())
	}
	if in.Err.Message() != "my_message" {
		t.Errorf("Expects Err.Message() is my_message. Got %v", in.Err.Message())
	}
	if in.Err.Trace().Error() != "my_err" {
		t.Errorf("Expects Err.Error() is my_err. Got %v", in.Err.Trace().Error())
	}

	// Asserting for in.Foo
	if in.Foo.Foo() != "Foo.." {
		t.Errorf("Expects Err.Foo() is Foo... Got %v", in.Foo.Foo())
	}
	if b, ok := in.Foo.(*InjectFoo); ok != true || b.Baz.Baz() != "Baz.." {
		t.Errorf("Expects Foo is an instance of InjectFoo and Foo.Baz.Baz() is Baz... Got %v", b.Baz.Baz())
	}
}

func TestFactoryContainer_Private_InstanceOf_AbstractNotAnInterface(t *testing.T) {
	c := &FactoryContainer{make(map[string]reflect.Value)}
	v := c.instanceOf(reflect.TypeOf("a_string"), reflect.TypeOf((*Bag)(nil)))
	if v != false {
		t.Errorf("Expects v is false. Got %v", v)
	}
}

func TestFactoryContainer_Private_InstanceOf_ConcreteTypeNotSupport(t *testing.T) {
	c := &FactoryContainer{make(map[string]reflect.Value)}
	i, _ := c.interfaceOf((*Bag)(nil))
	v := c.instanceOf(i, reflect.TypeOf("a_string"))
	if v != false {
		t.Errorf("Expects v is false. Got %v", v)
	}
}
