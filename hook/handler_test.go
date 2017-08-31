package hook

import (
	"encoding/json"
	"github.com/goline/lapi"
	"net/http"
	"testing"
)

func TestProcessHandlerResultHook_SetUp(t *testing.T) {
	h := &ProcessHandlerResultHook{}
	if h.SetUp(lapi.NewRequest(nil), lapi.NewResponse(nil)) != true {
		t.Errorf("Expects SetUp to return true. Got false")
	}
}

func TestProcessHandlerResultHook_TearDown_Error(t *testing.T) {
	h := &ProcessHandlerResultHook{}
	r := lapi.NewResponse(nil)
	h.TearDown(nil, r, nil, lapi.NewError("001.002.003", "MyMsg", nil))
	b, ok := json.Marshal(r.Content())
	if ok != nil || string(b) != `{"errors":[{"code":"001.002.003","message":"MyMsg"}]}` {
		t.Errorf("Expects correct format of error response. Got %v", string(b))
	}
}

func TestProcessHandlerResultHook_TearDown_HttpError(t *testing.T) {
	h := &ProcessHandlerResultHook{}
	r := lapi.NewResponse(nil)
	h.TearDown(nil, r, nil, lapi.NewHttpError(http.StatusInternalServerError, "001.002.003", "MyMsg", nil))
	b, ok := json.Marshal(r.Content())
	if ok != nil || string(b) != `{"errors":[{"code":"001.002.003","message":"MyMsg"}]}` {
		t.Errorf("Expects correct format of error response. Got %v", string(b))
	}
	if r.Status() != http.StatusInternalServerError {
		t.Errorf("Expects response has StatusInternalServerError status. Got %v", r.Status())
	}
}

func TestProcessHandlerResultHook_TearDown_StackError(t *testing.T) {
	h := &ProcessHandlerResultHook{}
	r := lapi.NewResponse(nil)
	se := lapi.NewStackError(400, []lapi.Error{
		lapi.NewError("001.002.003", "MSG1", nil),
		lapi.NewError("101.102.103", "MSG2", nil),
	})
	h.TearDown(nil, r, nil, se)
	b, ok := json.Marshal(r.Content())
	if ok != nil || string(b) != `{"errors":[{"code":"001.002.003","message":"MSG1"},{"code":"101.102.103","message":"MSG2"}]}` {
		t.Errorf("Expects correct format of stack error response. Got %v", string(b))
	}
	if r.Status() != http.StatusBadRequest {
		t.Errorf("Expects response has HttpBadRequest status. Got %v", r.Status())
	}
}

type itemResponse struct {
	Sku   string  `json:"sku"`
	Value float64 `json:"value"`
}

func TestProcessHandlerResultHook_TearDown(t *testing.T) {
	h := &ProcessHandlerResultHook{}
	r := lapi.NewResponse(nil)
	v := &itemResponse{"ABC-111", 6.9}
	h.TearDown(nil, r, v, nil)
	b, ok := json.Marshal(r.Content())
	if ok != nil || string(b) != `{"sku":"ABC-111","value":6.9}` {
		t.Errorf("Expects correct format of item response. Got %v", string(b))
	}
}
