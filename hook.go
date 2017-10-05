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

// SystemHook acts as mandatory hook
type SystemHook struct{}

func (h *SystemHook) SetUp(_ Connection) error {
	return nil
}

func (h *SystemHook) TearDown(c Connection, result interface{}, err error) error {
	if err != nil {
		// let rescuer handle error
		return err
	}

	if result == nil {
		return nil
	}

	if err := c.Response().Body().Write(result); err != nil {
		return err
	}

	return nil
}

// Priority implements Prioritizer interface
func (h *SystemHook) Priority() int {
	return PRIORITY_SYSTEM_HOOK
}

type ParserHook struct{}

func (h *ParserHook) SetUp(c Connection) error {
	parser := new(JsonParser)
	c.Request().Body().WithParser(parser)
	c.Response().Body().WithParser(parser)
	return nil
}

func (h *ParserHook) TearDown(_ Connection, _ interface{}, _ error) error {
	return nil
}
