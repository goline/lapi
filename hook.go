package lapi

// Hook acts as a middleware of processing request
type Hook interface {
	// SetUp executes Hook before handler runs
	// if one of hooks return false during process, hooking should be stopped
	SetUp(req Request, res Response) bool

	// TearDown executes Hook after handler runs
	// if one of hooks return false during process, hooking should be stopped
	TearDown(req Request, res Response, result interface{}, err error) bool
}

type ProcessHandlerResultHook struct{}

func (h *ProcessHandlerResultHook) SetUp(req Request, res Response) bool {
	return true
}

func (h *ProcessHandlerResultHook) TearDown(req Request, res Response, result interface{}, err error) bool {
	if err != nil {
		if e, ok := err.(ErrorStatus); ok == true {
			res.WithStatus(e.Status())
		}
	} else if result != nil {
		res.WithContent(result)
	}

	return true
}
