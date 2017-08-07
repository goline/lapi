package lapi

// Bagger manages key-value pairs
type Bagger interface {
	// Get returns value of specific key
	Get(key string, def interface{}) interface{}

	// Set allows to set value for a proposed key
	Set(key string, value interface{})

	// Has helps to check if a key exists
	Has(key string) bool
}