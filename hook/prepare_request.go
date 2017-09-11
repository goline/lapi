package hook

import (
	. "github.com/goline/lapi"
	"io/ioutil"
)

type PrepareRequestHook struct{}

func (h *PrepareRequestHook) SetUp(request Request, _ Response) error {
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

func (h *PrepareRequestHook) TearDown(_ Request, res Response, result interface{}, err error) error {
	return nil
}
