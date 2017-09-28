package lapi

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("App", func() {
	It("NewApp should return an instance of App", func() {
		Expect(NewApp()).NotTo(BeNil())
	})
})

var _ = Describe("FactoryApp", func() {
	It("should return response with error code ERR_HTTP_NOT_FOUND", func() {
		app := NewApp()
		app.Router().WithHook(new(SystemHook)).WithHook(new(ParserHook))
		app.Run()

		req := httptest.NewRequest("GET", "/foo", nil)
		rw := httptest.NewRecorder()
		app.ServeHTTP(rw, req)
		res := rw.Result()
		Expect(res.StatusCode).To(Equal(http.StatusNotFound))

		resErr := new(ErrorResponse)
		body, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())
		err = json.Unmarshal(body, resErr)
		Expect(err).To(BeNil())
		Expect(resErr.Code).To(Equal(ERR_HTTP_NOT_FOUND))
	})
})
