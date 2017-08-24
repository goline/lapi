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
	err := c.Bind((*Bag)(nil), FactoryError{})
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
	if err.Code() != BindErrorNotImplementedInterface {
		t.Errorf("Expects BindErrorNotImplementedInterface code. Got %v", err.Code())
	}
}

func TestFactoryContainer_Bind_ErrorNotSupportAbstractType(t *testing.T) {
	c := NewContainer()
	err := c.Bind((*Bag)(nil), "a_string")
	if err == nil {
		t.Errorf("Expects err is not nil")
	}
	if err.Code() != BindErrorInvalidConcrete {
		t.Errorf("Expects BindErrorInvalidConcrete code. Got %v", err.Code())
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

func TestFactoryContainer_Bind_Struct_Ok(t *testing.T) {
	c := NewContainer()
	err := c.Bind((*Error)(nil), &FactoryError{"my_code", "my_message", errors.New("my_error")})
	if err != nil {
		t.Errorf("Expects error is nil")
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
