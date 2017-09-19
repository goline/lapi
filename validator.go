package lapi

import (
	"errors"
	"fmt"
)

// Validator helps validate inputs
type Validator interface {
	// Validate checks input against rules
	Validate(input Bag, rules Rules) error
}

func NewValidator() Validator {
	return &FactoryValidator{}
}

type FactoryValidator struct{}

func (v *FactoryValidator) Validate(input Bag, rules Rules) error {
	for key, value := range input.All() {
		checks := rules.Get(key)
		if len(checks) == 0 {
			continue
		}
		for _, checker := range checks {
			if err := checker.Check(value); err != nil {
				return errors.New(fmt.Sprintf("%s: %s", key, err.Error()))
			}
		}
	}

	return nil
}

// Checker validates input
type Checker interface {
	// Check verifies value
	Check(value interface{}) error
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
