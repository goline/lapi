package lapi

// Validator helps validate inputs
type Validator interface {
	// Validate checks input against rules
	Validate(input Input, rules Rules) (bool, []error)
}

// Checker validates input
type Checker interface {
	// Check verifies value
	Check(value interface{}) bool

	// Name returns it's name
	Name() string

	// SetTranslator sets translator
	SetTranslator(translator Translator)
}

type Rules interface {
	// Add queues proposed key with checkers
	Add(key string, checkers ...Checker)

	// Del removes specific key from rules
	Del(key string)

	// Get returns all checkers of provided key
	Get(key string) []Checker

	// All returns all keys-checkers
	All() map[string][]Checker
}