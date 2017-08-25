package lapi

import (
	"errors"
	"testing"
)

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	if _, ok := c.(Container); ok == false {
		t.Errorf("Expects an instance of Container. Got %+v", c)
	}
}

func TestFactoryContainer_Bind_ErrorAbstractIsNotAnInterface(t *testing.T) {
	c := NewContainer()
	err := c.Bind(&FactoryBag{}, &FactoryError{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
	if err.Code() != BindErrorInvalidInterface {
		t.Errorf("Expects BindErrorInvalidInterface. Got %v", err.Code())
	}
}

func TestFactoryContainer_Bind_ErrorConcreteNotImplementAbstract(t *testing.T) {
	c := NewContainer()
	err := c.Bind((*Bag)(nil), &FactoryError{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
	if err.Code() != BindErrorNotImplementedInterface {
		t.Errorf("Expects BindErrorNotImplementedInterface code. Got %v", err.Code())
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
		if err.Code() != BindErrorInvalidConcrete {
			t.Errorf("Expects BindErrorInvalidConcrete code. Got %v", err.Code())
		}
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
	if err.Code() != ResolveErrorInsufficientArguments {
		t.Errorf("Expects ResolveErrorInsufficientArguments code. Got %v", err.Code())
	}
}

func TestFactoryContainer_Resolve_Func_ErrorNonValuesReturned(t *testing.T) {
	c := NewContainer()
	c.Bind((*Error)(nil), func(code string, message string, err error) {})
	_, err := c.Resolve((*Error)(nil), "my_code", "my_message", errors.New("my_err"))
	if err == nil {
		t.Errorf("Expects error is not nil")
	}
	if err.Code() != ResolveErrorNonValuesReturned {
		t.Errorf("Expects ResolveErrorNonValuesReturned code. Got %v", err.Code())
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

func TestFactoryContainer_Inject_ErrorInjectInvalidTargetType(t *testing.T) {
	c := NewContainer()
	e := FactoryError{"my_code", "my_message", errors.New("my_err")}
	err := c.Inject(e)
	if err == nil {
		t.Errorf("Expects error is not nil")
	}
	if err.Code() != InjectErrorInvalidTargetType {
		t.Errorf("Expects InjectErrorInvalidTargetType code. Got %v", err.Code())
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
	if err.Code() != ResolveErrorNotExistAbstract {
		t.Errorf("Expects ResolveErrorNotExistAbstract code. Got %v", err.Code())
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
	Err Error       `inject:"*"`
	Foo InjectFooer `inject:"*"`
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
