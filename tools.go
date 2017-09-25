package lapi

import (
	"reflect"
	"sort"
	"sync"

	"github.com/goline/errors"
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

// StructOf returns type of struct
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

type iteratorFunc func(item interface{})

func Parallel(list map[int][]interface{}, f iteratorFunc) {
	indexes := make(sort.IntSlice, 0)
	for i := range list {
		indexes = append(indexes, i)
	}
	sort.Sort(indexes)

	var wg sync.WaitGroup
	for _, i := range indexes {
		items := list[i]
		for _, item := range items {
			wg.Add(1)
			go func(f iteratorFunc, item interface{}) {
				defer wg.Done()
				f(item)
			}(f, item)
		}
		wg.Wait()
	}
}
