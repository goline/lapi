package lapi

import (
	"io/ioutil"
	"net/http"
)

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
	if (request.Method() != http.MethodPost &&
		request.Method() != http.MethodPut &&
		request.Method() != http.MethodPatch) ||
		request.Route().RequestInput() == nil {
		return nil
	}

	body := request.Ancestor().Body
	defer body.Close()

	content, err := ioutil.ReadAll(request.Ancestor().Body)
	if err != nil {
		return err
	}
	request.WithContentBytes(content, request.Route().RequestInput())

	return nil
}

func (h *SystemHook) TearDown(connection Connection, result interface{}, err error) error {
	if err != nil {
		// let rescuer handle error
		return err
	}

	if result == nil {
		return nil
	}

	connection.Response().WithContent(result)
	return nil
}

type ParserHook struct{}

func (h *ParserHook) SetUp(connection Connection) error {
	parser := new(JsonParser)
	connection.Request().WithParser(parser)
	connection.Response().WithParser(parser)
	return nil
}

func (h *ParserHook) TearDown(connection Connection, result interface{}, err error) error {
	return nil
}
