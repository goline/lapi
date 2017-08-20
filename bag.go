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
	inputs map[string]interface{}
}

func (i *FactoryBag) Get(key string) (interface{}, bool) {
	value, ok := i.inputs[key]
	return value, ok
}

func (i *FactoryBag) Set(key string, value interface{}) {
	i.inputs[key] = value
}

func (i *FactoryBag) Remove(key string) {
	delete(i.inputs, key)
}

func (i *FactoryBag) Has(key string) bool {
	_, ok := i.inputs[key]
	return ok
}

func (i *FactoryBag) All() map[string]interface{} {
	return i.inputs
}