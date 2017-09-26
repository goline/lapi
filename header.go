package lapi

import (
	"strings"
	"sync"
)

type Header interface {
	// Get returns value of specific case-insensitive key
	Get(key string) (string, bool)

	// Set allows to set value for a proposed case-insensitive key
	Set(key string, value string)

	// Remove deletes a specific case-insensitive key from Bag
	Remove(key string)

	// Has helps to check if a case-insensitive key exists
	Has(key string) bool

	// All returns all key-value of bag
	All() map[string]string
}

// NewHeader returns an instance of Header
func NewHeader() Header {
	return &FactoryHeader{new(sync.Map)}
}

type FactoryHeader struct {
	items *sync.Map
}

func (h *FactoryHeader) Get(key string) (string, bool) {
	v, ok := h.items.Load(h.formatKey(key))
	if ok == false {
		return "", false
	}

	value := v.(string)
	return value, ok
}

func (h *FactoryHeader) Set(key string, value string) {
	h.items.Store(h.formatKey(key), value)
}

func (h *FactoryHeader) Remove(key string) {
	h.items.Delete(h.formatKey(key))
}

func (h *FactoryHeader) Has(key string) bool {
	_, ok := h.items.Load(h.formatKey(key))
	return ok
}

func (h *FactoryHeader) All() map[string]string {
	items := make(map[string]string)
	h.items.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		items[k] = v
		return true
	})
	return items
}

func (h *FactoryHeader) formatKey(key string) string {
	return strings.Title(strings.ToLower(key))
}
