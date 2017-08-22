package lapi

import (
	"strings"
)

type Header interface {
	// Get returns value of specific case-insensitive key
	Get(key string) ([]string, bool)

	// Set allows to set value for a proposed case-insensitive key
	Set(key string, value []string)

	// Remove deletes a specific case-insensitive key from Bag
	Remove(key string)

	// Has helps to check if a case-insensitive key exists
	Has(key string) bool

	// All returns all key-value of bag
	All() map[string][]string

	// Line shows a specific case-insensitive key as a line
	// it returns empty if key is not found
	Line(key string) string

	// Lines returns header as lines
	Lines() []string
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

func (h *FactoryHeader) Set(key string, value []string) {
	h.items[h.formatKey(key)] = value
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

func (h *FactoryHeader) Line(key string) string {
	line := ""
	key = h.formatKey(key)
	values, ok := h.items[key]
	if ok == true {
		line = key + ": " + strings.Join(values, ", ")
	}

	return line
}

func (h *FactoryHeader) Lines() []string {
	lines := make([]string, len(h.items))
	i := 0
	for key := range h.items {
		lines[i] = h.Line(key)
		i++
	}
	return lines
}

func (h *FactoryHeader) formatKey(key string) string {
	return strings.Title(strings.ToLower(key))
}
