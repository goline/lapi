package lapi

import "io/ioutil"

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

func (h *SystemHook) SetUp(connection Connection) error {
	request := connection.Request()
	if request.Ancestor().Body == nil {
		return nil
	}

	rb, ok := request.Route().(RouteBodyIO)
	if ok == false {
		return nil
	}

	body, err := ioutil.ReadAll(request.Ancestor().Body)
	if err != nil {
		return err
	}
	request.WithContentBytes(body, rb.RequestInput())
	return nil
}

func (h *SystemHook) TearDown(connection Connection, result interface{}, err error) error {
	if err != nil {
		// let rescuer handle error
		return err
	}

	connection.Response().WithContent(result)
	return nil
}
