package lapi

// Validator helps validate inputs
type Validator interface {
	// Validate checks input against rules
	Validate(input Bag, rules Rules) (bool, StackError)

	// Translator returns an instance of Translator
	Translator() Translator
}

// Checker validates input
type Checker interface {
	// Check verifies value
	Check(value interface{}) bool

	// ErrorMessage returns format of error message
	ErrorMessage() string
}

// Skipper allows to skip
// Any Checker (or other) which might takes time to process
// should implement this interface
type Skipper interface {
	// Skip returns true if we should skip
	Skip() bool
}

// Rules manages validation rules
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