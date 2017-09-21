package lapi

import "net/http"

// Rescuer handles error
type Rescuer interface {
	// Rescue handles error, it returns nil if error is handled,
	// and error itself if could not be handled properly
	// Server should panic if an error is returned
	Rescue(connection Connection, err error) error
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

type FactoryRescuer struct{}

func (r *FactoryRescuer) Rescue(c Connection, err error) error {
	if c == nil {
		return err
	}

	if e, ok := err.(ErrorStatuser); ok == true {
		c.Response().WithStatus(e.Status())
	}
	var code, message string
	if e, ok := err.(ErrorCoder); ok == true {
		code = e.Code()
		switch code {
		case ERR_HTTP_NOT_FOUND:
			c.Response().WithStatus(http.StatusNotFound)
		case ERR_HTTP_BAD_REQUEST:
			c.Response().WithStatus(http.StatusBadRequest)
		case ERR_HTTP_INTERNAL_SERVER_ERROR:
			c.Response().WithStatus(http.StatusInternalServerError)
		default:
			c.Response().WithStatus(http.StatusInternalServerError)
		}
	} else {
		code = ERR_HTTP_UNKNOWN_ERROR
	}
	if e, ok := err.(ErrorMessager); ok == true {
		message = e.Message()
	} else {
		message = err.Error()
	}
	c.Response().WithContent(&ErrorResponse{code, message})

	return nil
}
