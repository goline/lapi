package lapi

import (
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
