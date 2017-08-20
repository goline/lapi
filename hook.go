package lapi

// Hook acts as a middleware of processing request
type Hook interface {
	// Run executes Hook, hooking process should be stopped
	// if one of hooks return false during process
	Run(req Request, res Response) bool
}