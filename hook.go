package lapi

import "github.com/goline/errors"

// Hook acts as a middleware of processing request
type Hook interface {
	// Since v1.0.14
	// Hook will become an empty interface
	// User should implement either BootableHook or HaltableHook
}

// BootableHook allows to register hook to be executed before handler runs
type BootableHook interface {
	// SetUp executes Hook before handler runs
	// if one of hooks returns an error during process, hooking should be stopped
	SetUp(connection Connection) errors.Error
}

// HaltableHook allows to register hook to be executed after handler runs
type HaltableHook interface {
	// TearDown executes Hook after handler runs
	// if one of hooks returns an error during process, hooking should be stopped
	TearDown(connection Connection, result interface{}, err errors.Error) errors.Error
}

// SystemHook acts as mandatory hook
type SystemHook struct {}

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

type ParserHook struct {}

func (h *ParserHook) SetUp(c Connection) errors.Error {
	parser := new(JsonParser)
	c.Request().Body().WithParser(parser)
	c.Response().Body().WithParser(parser)
	return nil
}
