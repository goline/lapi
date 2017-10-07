package lapi

import "github.com/goline/errors"

// Hook acts as a middleware of processing request
type Hook interface {
	// SetUp executes Hook before handler runs
	// if one of hooks returns an error during process, hooking should be stopped
	SetUp(connection Connection) errors.Error

	// TearDown executes Hook after handler runs
	// if one of hooks returns an error during process, hooking should be stopped
	TearDown(connection Connection, result interface{}, err errors.Error) errors.Error
}

// SystemHook acts as mandatory hook
type SystemHook struct{}

func (h *SystemHook) SetUp(_ Connection) errors.Error {
	return nil
}

func (h *SystemHook) TearDown(c Connection, result interface{}, err errors.Error) errors.Error {
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

func (h *ParserHook) SetUp(c Connection) errors.Error {
	parser := new(JsonParser)
	c.Request().Body().WithParser(parser)
	c.Response().Body().WithParser(parser)
	return nil
}

func (h *ParserHook) TearDown(_ Connection, _ interface{}, _ errors.Error) errors.Error {
	return nil
}
