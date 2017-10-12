package lapi

import (
	"reflect"
	"strconv"
	"strings"
)

// Bag manages key-value pairs
type Bag interface {
	// Get returns value of specific key
	Get(key string) (interface{}, bool)

	// Set allows to set value for a proposed key
	Set(key string, value interface{})

	// Remove deletes a specific key from Bag
	Remove(key string)

	// Has helps to check if a key exists
	Has(key string) bool

	// All returns all key-value of bag
	All() map[string]interface{}

	BagGetter
}

type BagGetter interface {
	// GetInt returns int64 value
	GetInt(key string) (int64, bool)

	// GetFloat returns float64 value
	GetFloat(key string) (float64, bool)

	// GetString returns string value
	GetString(key string) (string, bool)

	// GetBool returns boolean value
	GetBool(key string) (bool, bool)
}

// NewBag returns an instance of Bag
func NewBag() Bag {
	return &FactoryBag{make(map[string]interface{})}
}

type FactoryBag struct {
	items map[string]interface{}
}

func (b *FactoryBag) Get(key string) (interface{}, bool) {
	value, ok := b.items[key]
	return value, ok
}

func (b *FactoryBag) Set(key string, value interface{}) {
	b.items[key] = value
}

func (b *FactoryBag) Remove(key string) {
	delete(b.items, key)
}

func (b *FactoryBag) Has(key string) bool {
	_, ok := b.items[key]
	return ok
}

func (b *FactoryBag) All() map[string]interface{} {
	return b.items
}

func (b *FactoryBag) GetInt(key string) (int64, bool) {
	value, ok := b.Get(key)
	if ok == false {
		return 0, false
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(value).Int(), true
	case reflect.String:
		s := reflect.ValueOf(value).String()
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return i, true
		}
	}

	return 0, false
}

func (b *FactoryBag) GetFloat(key string) (float64, bool) {
	value, ok := b.Get(key)
	if ok == false {
		return 0.0, false
	}

	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f, true
		}
	}

	return 0.0, false
}

func (b *FactoryBag) GetString(key string) (string, bool) {
	value, ok := b.Get(key)
	if ok == false {
		return "", false
	}

	v, ok := value.(string)
	if ok == false {
		return "", false
	}

	return v, true
}

func (b *FactoryBag) GetBool(key string) (bool, bool) {
	value, ok := b.Get(key)
	if ok == false {
		return false, false
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool:
		return reflect.ValueOf(value).Bool(), true
	case reflect.String:
		v := reflect.ValueOf(value).String()
		if strings.Compare(v, "true") == 0 {
			return true, true
		} else if strings.Compare(v, "false") == 0 {
			return false, true
		}

		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i == int64(1), true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(value).Int() == int64(1), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(value).Uint() == uint64(1), true
	}

	return false, false
}
