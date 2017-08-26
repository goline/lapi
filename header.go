package lapi

import (
	"strings"
)

type Header interface {
	// Get returns value of specific case-insensitive key
	Get(key string) ([]string, bool)

	// Set allows to set value for a proposed case-insensitive key
	Set(key string, values ...string)

	// Remove deletes a specific case-insensitive key from Bag
	Remove(key string)

	// Has helps to check if a case-insensitive key exists
	Has(key string) bool

	// All returns all key-value of bag
	All() map[string][]string

	// Line shows a specific case-insensitive key as a pair of name-line
	// it returns empty string for line if key is not found
	Line(key string) (string, string)

	// Lines returns header as lines
	Lines() map[string]string
}

// NewHeader returns an instance of Header
func NewHeader() Header {
	return &FactoryHeader{make(map[string][]string)}
}

type FactoryHeader struct {
	items map[string][]string
}

func (h *FactoryHeader) Get(key string) ([]string, bool) {
	value, ok := h.items[h.formatKey(key)]
	return value, ok
}

func (h *FactoryHeader) Set(key string, values ...string) {
	h.items[h.formatKey(key)] = values
}

func (h *FactoryHeader) Remove(key string) {
	delete(h.items, h.formatKey(key))
}

func (h *FactoryHeader) Has(key string) bool {
	_, ok := h.items[h.formatKey(key)]
	return ok
}

func (h *FactoryHeader) All() map[string][]string {
	return h.items
}

func (h *FactoryHeader) Line(key string) (string, string) {
	line := ""
	key = h.formatKey(key)
	values, ok := h.items[key]
	if ok == true {
		line = strings.Join(values, ", ")
	}

	return key, line
}

func (h *FactoryHeader) Lines() map[string]string {
	lines := make(map[string]string, len(h.items))
	for key := range h.items {
		name, line := h.Line(key)
		lines[name] = line
	}
	return lines
}

func (h *FactoryHeader) formatKey(key string) string {
	return strings.Title(strings.ToLower(key))
}
