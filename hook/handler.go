package hook

import (
	"errors"
	. "github.com/goline/lapi"
)

type errorStackResponse struct {
	Errors []errorItemResponse `json:"errors"`
}

type errorItemResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ProcessHandlerResultHook struct{}

func (h *ProcessHandlerResultHook) SetUp(_ Request, _ Response) bool {
	return true
}

func (h *ProcessHandlerResultHook) TearDown(_ Request, res Response, result interface{}, err error) bool {
	if err != nil {
		if e, ok := err.(ErrorStatus); ok == true {
			res.WithStatus(e.Status())
		}

		var es []Error
		if e, ok := err.(Error); ok == true {
			es = make([]Error, 1)
			es[0] = e
		} else if e, ok := err.(StackError); ok == true {
			es = e.Errors()
		}

		if len(es) > 0 {
			ei := make([]errorItemResponse, len(es))
			for i, e := range es {
				ei[i] = errorItemResponse{e.Code(), e.Message()}
			}
			er := &errorStackResponse{ei}
			res.WithContent(er)
		} else {
			panic(errors.New("Expects to catch at least 1 error. Got 0"))
		}
	} else if result != nil {
		res.WithContent(result)
	}

	return true
}
