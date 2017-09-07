package lapi

// Validator helps validate inputs
type Validator interface {
	// Validate checks input against rules
	Validate(input Bag, rules Rules) (bool, StackError)
}

// Checker validates input
type Checker interface {
	// Check verifies value
	Check(value interface{}) bool
	ErrorMessager
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
	Add(key string, checkers ...Checker) Rules

	// Remove deletes specific key from rules
	Remove(key string) Rules

	// Get returns all checkers of provided key, returns empty if key is not found
	Get(key string) []Checker

	// All returns all keys-checkers
	All() map[string][]Checker
}

func NewRules() Rules {
	return &FactoryRules{make(map[string][]Checker)}
}

type FactoryRules struct {
	rules map[string][]Checker
}

func (r *FactoryRules) Add(key string, checkers ...Checker) Rules {
	r.rules[key] = make([]Checker, len(checkers))
	r.rules[key] = checkers
	return r
}

func (r *FactoryRules) Remove(key string) Rules {
	delete(r.rules, key)
	return r
}

func (r *FactoryRules) Get(key string) []Checker {
	checkers, ok := r.rules[key]
	if ok == false {
		return make([]Checker, 0)
	}

	return checkers
}

func (r *FactoryRules) All() map[string][]Checker {
	return r.rules
}
