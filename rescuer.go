package lapi

import (
	"fmt"
	"net/http"

	"github.com/goline/errors"
)

// Rescuer handles error
type Rescuer interface {
	// Rescue handles error, it returns nil if error is handled,
	// and error itself if could not be handled properly
	// Server should panic if an error is returned
	Rescue(connection Connection, v interface{}) error
}

func NewRescuer() Rescuer {
	return &FactoryRescuer{}
}

type ErrorResponse struct {
	// The error code
	// Required: true
	Code string `json:"code"`

	// The error message
	// Required: true
	Message string `json:"message"`
}

type FactoryRescuer struct {
	parser Parser
}

func (r *FactoryRescuer) Rescue(c Connection, v interface{}) error {
	if c == nil {
		return errors.New(ERR_INVALID_ARGUMENT, "Connection must be not nil")
	}
	if r.parser == nil {
		r.parser = new(JsonParser)
	}
	c.Response().Body().
		WithContentType(CONTENT_TYPE_JSON).
		WithParser(r.parser)

	var code, message string
	code = ERR_HTTP_UNKNOWN_ERROR
	if e, ok := v.(errors.Error); ok == true {
		code = e.Code()
		switch code {
		case ERR_HTTP_NOT_FOUND:
			c.Response().WithStatus(http.StatusNotFound)
		case ERR_HTTP_BAD_REQUEST:
			c.Response().WithStatus(http.StatusBadRequest)
		case ERR_HTTP_INTERNAL_SERVER_ERROR:
			c.Response().WithStatus(http.StatusInternalServerError)
		default:
			if e.Status() == http.StatusOK {
				c.Response().WithStatus(http.StatusInternalServerError)
			} else if c.Response().Status() == http.StatusOK {
				c.Response().WithStatus(e.Status())
			}
		}
		message = e.Message()
	} else if e, ok := v.(error); ok == true {
		message = e.Error()
		c.Response().WithStatus(http.StatusInternalServerError)
	} else {
		message = fmt.Sprintf("%s", v)
		c.Response().WithStatus(http.StatusInternalServerError)
	}
	if err := c.Response().Body().Write(&ErrorResponse{code, message}); err != nil {
		return err
	}

	return nil
}
