package lapi

// Hook acts as a middleware of processing request
type Hook interface {
	// SetUp executes Hook before handler runs
	// if one of hooks returns an error during process, hooking should be stopped
	SetUp(connection Connection) error

	// TearDown executes Hook after handler runs
	// if one of hooks returns an error during process, hooking should be stopped
	TearDown(connection Connection, result interface{}, err error) error
}
