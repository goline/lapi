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
		switch e := err.(type) {
		case Error:
			es = make([]Error, 1)
			es[0] = e
		case StackError:
			es = e.Errors()
		default:
			es[0] = NewError("", "INVALID_ERROR", errors.New("Error's type is not supported."))
		}

		ei := make([]errorItemResponse, len(es))
		for i, e := range es {
			ei[i] = errorItemResponse{e.Code(), e.Message()}
		}
		er := &errorStackResponse{ei}
		res.WithContent(er)
	} else if result != nil {
		res.WithContent(result)
	}

	return true
}
