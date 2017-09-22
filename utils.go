package lapi

import (
	"github.com/goline/errors"
	"reflect"
)

// PanicOnError panics if input value is an error
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Must panics if there is an error
func Must(errors ...error) {
	for _, err := range errors {
		PanicOnError(err)
	}
}

// Clone returns a pointer which is a copied of input type
func Clone(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Ptr:
		return reflect.New(t.Elem()).Interface()
	case reflect.Struct:
		return reflect.New(t).Interface()
	default:
		panic(errors.New(ERR_CLONE_INVALID_TYPE, "Clone invalid type. Support: Ptr, Struct"))
	}
}

func StructOf(v interface{}) reflect.Type {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Ptr:
		return t.Elem()
	case reflect.Struct:
		return t
	default:
		panic(errors.New(ERR_STRUCT_INVALID_TYPE, "Calls StructOf with invalid type. Support: Ptr, Struct"))
	}
}
