package lapi

import (
	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Rescuer", func() {
	It("NewRescuer should return an instance of Rescuer", func() {
		Expect(NewRescuer()).NotTo(BeNil())
	})
})

type myUnknownError struct{}

func (e *myUnknownError) Error() string { return "" }
func (e *myUnknownError) Status() int   { return http.StatusInternalServerError }

var _ = Describe("FactoryRescuer", func() {
	It("Rescue set http status to StatusNotFound", func() {
		c := NewConnection(nil, getEmptyResponse())
		e := errors.New(ERR_HTTP_NOT_FOUND, "")
		h := &FactoryRescuer{}
		h.Rescue(c, e)
		Expect(c.Response().Status()).To(Equal(http.StatusNotFound))
	})

	It("Rescue set http status to StatusBadRequest", func() {
		c := NewConnection(nil, getEmptyResponse())
		e := errors.New(ERR_HTTP_BAD_REQUEST, "")
		h := &FactoryRescuer{}
		h.Rescue(c, e)
		Expect(c.Response().Status()).To(Equal(http.StatusBadRequest))
	})

	It("Rescue set http status to StatusInternalServerError for unknown error", func() {
		c := NewConnection(nil, getEmptyResponse())
		e := errors.New("11", "err1")
		h := &FactoryRescuer{}
		h.Rescue(c, e)
		Expect(c.Response().Status()).To(Equal(http.StatusInternalServerError))
	})

	It("Rescue will not handle this case", func() {
		e := &myUnknownError{}
		h := &FactoryRescuer{}
		err := h.Rescue(nil, e)
		Expect(err).NotTo(BeNil())
	})
})

func getEmptyResponse() Response {
	return &FactoryResponse{body: NewBody(nil, nil)}
}
