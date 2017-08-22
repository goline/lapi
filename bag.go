package lapi

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