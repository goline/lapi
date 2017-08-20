package lapi

import (
	"testing"
)

func TestNewBag(t *testing.T) {
	b := NewBag()
	if _, ok := b.(Bag); ok == false {
		t.Errorf("Expect an instance of Bag. Got %+v", b)
	}
}

func TestFactoryBag_Get(t *testing.T) {
	b := &FactoryBag{make(map[string]interface{})}
	b.inputs["my_key"] = "my_value"
	if v, ok := b.Get("my_key"); v != "my_value" || ok == false {
		t.Errorf("Expects my_key equals my_value. Got %+v", v)
	}
}

func TestFactoryBag_Has(t *testing.T) {
	b := &FactoryBag{make(map[string]interface{})}
	b.inputs["my_key"] = "my_value"
	if b.Has("my_key") == false {
		t.Errorf("Expects my_key exists")
	}
	if b.Has("my_another_key") == true {
		t.Errorf("Expects my_another_key not exist")
	}
}

func TestFactoryBag_Set(t *testing.T) {
	b := &FactoryBag{make(map[string]interface{})}
	if b.Has("my_key") == false {
		b.Set("my_key", "my_value")
		if b.Has("my_key") == false {
			t.Errorf("Expects my_key exists")
		}
	}
}

func TestFactoryBag_Remove(t *testing.T) {
	b := &FactoryBag{make(map[string]interface{})}
	b.inputs["my_key"] = "my_value"
	b.Remove("my_key")
	if _, ok := b.Get("my_key"); ok == true {
		t.Errorf("Expects my_key is removed")
	}
}

func TestFactoryBag_All(t *testing.T) {
	b := &FactoryBag{make(map[string]interface{})}
	b.inputs["my_key"] = "my_value"
	b.inputs["my_another_key"] = 1
	a := b.All()
	if len(a) != 2 {
		t.Errorf("Expects a only has 2 items")
	}
}