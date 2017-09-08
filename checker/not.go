package checker

import (
	"fmt"
	. "github.com/goline/lapi"
)

func Not(checker Checker) Checker {
	return &NotChecker{checker}
}

type NotChecker struct {
	checker Checker
}

func (c *NotChecker) Check(value interface{}) bool {
	return !c.checker.Check(value)
}

func (c *NotChecker) Message() string {
	return fmt.Sprintf("MUST NOT (%s)", c.checker.Message())
}
