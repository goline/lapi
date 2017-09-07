package checker

type RequireChecker struct{}

func (c *RequireChecker) Check(value interface{}) bool {
	if value == nil {
		return false
	}
	return true
}

func (c *RequireChecker) Message() string {
	return "%s is required"
}
