package lapi

// Hook acts as a middleware of processing request
type Hook interface {
	// SetUp executes Hook before handler runs
	// if one of hooks return false during process, hooking should be stopped
	SetUp(req Request, res Response) bool

	// TearDown executes Hook after handler runs
	// if one of hooks return false during process, hooking should be stopped
	TearDown(req Request, res Response) bool
}
