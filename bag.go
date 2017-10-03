package lapi

import (
	"strconv"
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
	// GetInt64 returns int64 value
	GetInt64(key string) (int64, bool)

	// GetFloat64 returns float64 value
	GetFloat64(key string) (float64, bool)

	// GetString returns string value
	GetString(key string) (string, bool)
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

func (b *FactoryBag) GetInt64(key string) (int64, bool) {
	value, ok := b.Get(key)
	if ok == false {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i, true
		}
	}

	return 0, false
}

func (b *FactoryBag) GetFloat64(key string) (float64, bool) {
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
