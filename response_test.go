package lapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Response", func() {
	It("NewResponse should return an instance of Response", func() {
		Expect(NewResponse(nil)).NotTo(BeNil())
	})
})

var _ = Describe("FactoryResponse", func() {
	It("Status should return http.StatusOk as default", func() {
		r := NewResponse(nil)
		Expect(r.Status()).To(Equal(http.StatusOK))
	})

	It("Status should return http status", func() {
		r := &FactoryResponse{}
		r.status = http.StatusBadRequest
		Expect(r.Status()).To(Equal(http.StatusBadRequest))
	})

	It("WithStatus should allow to set http status", func() {
		r := &FactoryResponse{}
		r.WithStatus(http.StatusBadRequest)
		Expect(r.status).To(Equal(http.StatusBadRequest))
	})

	It("Message should return http Message", func() {
		r := &FactoryResponse{}
		r.message = "my_own_message"
		Expect(r.Message()).To(Equal("my_own_message"))
	})

	It("WithMessage should set http Message", func() {
		r := &FactoryResponse{}
		r.WithMessage("my_own_message")
		Expect(r.message).To(Equal("my_own_message"))
	})

	It("Header should return Header instance", func() {
		r := &FactoryResponse{}
		h := NewHeader()
		h.Set("content-type", "application/json")
		r.header = h
		rh, ok := r.Header().Get("content-type")
		Expect(ok).To(BeTrue())
		Expect(rh).To(Equal("application/json"))
	})

	It("WithHeader should set header instance", func() {
		r := &FactoryResponse{}
		h := NewHeader()
		h.Set("content-type", "application/json")
		rh, ok := r.WithHeader(h).Header().Get("content-type")
		Expect(ok).To(BeTrue())
		Expect(rh).To(Equal("application/json"))
	})

	It("Cookies should return cookies", func() {
		r := &FactoryResponse{}
		c1 := &http.Cookie{Name: "my_c1", Value: "val_c1"}
		c2 := &http.Cookie{Name: "my_c2", Value: "val_c2"}
		r.cookies = make([]*http.Cookie, 2)
		r.cookies[0] = c1
		r.cookies[1] = c2
		Expect(len(r.Cookies())).To(Equal(2))
		Expect(r.Cookies()[0].Name).To(Equal("my_c1"))
		Expect(r.Cookies()[1].Value).To(Equal("val_c2"))
	})

	It("WithCookies should set cookies", func() {
		r := &FactoryResponse{}
		c1 := &http.Cookie{Name: "my_c1", Value: "val_c1"}
		c2 := &http.Cookie{Name: "my_c2", Value: "val_c2"}
		c := make([]*http.Cookie, 2)
		c[0] = c1
		c[1] = c2
		Expect(len(r.WithCookies(c).Cookies())).To(Equal(2))
		Expect(r.Cookies()[0].Name).To(Equal("my_c1"))
		Expect(r.Cookies()[1].Value).To(Equal("val_c2"))
	})

	It("Send should flush response", func() {
		h := func(w http.ResponseWriter, r *http.Request) {
			rs := &FactoryResponse{ancestor: w, body: NewBody(nil, w), header: NewHeader()}
			rs.Body().WithContentType(CONTENT_TYPE_JSON)
			rs.Body().WithParser(new(JsonParser))
			rs.WithStatus(http.StatusInternalServerError)
			rs.Header().Set("x-custom-flag", "1234")
			c := make(map[string]interface{})
			c["status"] = true
			c["errors"] = false
			rs.Body().Write(c)
			rs.WithCookies([]*http.Cookie{
				&http.Cookie{Name: "my_c1", Value: "my_val1"},
				&http.Cookie{Name: "my_c2", Value: "my_val2"},
			})
			rs.Send()
		}
		rq := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		h(w, rq)

		r := w.Result()
		Expect(r.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(r.Header.Get("x-custom-flag")).To(Equal("1234"))
		Expect(r.Header.Get("Content-Type")).To(Equal("application/json"))

		b, err := ioutil.ReadAll(r.Body)
		Expect(err).To(BeNil())
		var v map[string]interface{}
		err = json.Unmarshal(b, &v)
		Expect(err).To(BeNil())
		Expect(v["errors"]).To(BeFalse())
		Expect(v["status"]).To(BeTrue())

		Expect(r.Cookies()[0].Name).To(Equal("my_c1"))
		Expect(r.Cookies()[0].Value).To(Equal("my_val1"))
		Expect(r.Cookies()[1].Name).To(Equal("my_c2"))
		Expect(r.Cookies()[1].Value).To(Equal("my_val2"))
	})

	It("Send should return error code ERR_NO_WRITER_FOUND", func() {
		r := &FactoryResponse{}
		err := r.Send()
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_NO_WRITER_FOUND))
	})

	It("Send should return error code ERR_RESPONSE_ALREADY_SENT", func() {
		r := &FactoryResponse{ancestor: httptest.NewRecorder()}
		r.isSent = true
		err := r.Send()
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESPONSE_ALREADY_SENT))
	})

	It("Send should return error code ERR_CONTENT_TYPE_EMPTY", func() {
		w := httptest.NewRecorder()
		r := &FactoryResponse{ancestor: w, body: NewBody(nil, w)}
		err := r.Send()
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_CONTENT_TYPE_EMPTY))
	})

	It("Send should delivery header content-type", func() {
		w := httptest.NewRecorder()
		r := &FactoryResponse{ancestor: w, header: NewHeader(), body: NewBody(nil, w)}
		r.Body().WithContentType(CONTENT_TYPE_JSON)
		r.Body().WithCharset(CONTENT_CHARSET_DEFAULT)
		r.Send()
		v, ok := r.header.Get(HEADER_CONTENT_TYPE)
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(fmt.Sprintf("%s; charset=%s", CONTENT_TYPE_JSON, CONTENT_CHARSET_DEFAULT)))
	})

	It("Send should return error code ERR_NO_PARSER_FOUND", func() {
		w := httptest.NewRecorder()
		r := &FactoryResponse{ancestor: w, header: NewHeader(), body: NewBody(nil, w)}
		r.Body().WithContentType(CONTENT_TYPE_JSON)
		r.Body().WithCharset(CONTENT_CHARSET_DEFAULT)
		r.Body().Write("a_string")
		err := r.Send()
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_NO_PARSER_FOUND))
	})

	It("Send should send http message", func() {
		w := httptest.NewRecorder()
		r := &FactoryResponse{ancestor: w, header: NewHeader(), body: NewBody(nil, w)}
		r.Body().WithContentType(CONTENT_TYPE_JSON)
		r.Body().WithCharset(CONTENT_CHARSET_DEFAULT)
		r.WithMessage("HTTP NOT FOUND")
		err := r.Send()
		Expect(err).To(BeNil())
	})

	It("Send should return error code ERR_RESPONSE_IS_SENDING", func() {
		w := httptest.NewRecorder()
		r := &FactoryResponse{ancestor: w, header: NewHeader(), body: NewBody(nil, w)}
		r.Body().WithContentType(CONTENT_TYPE_JSON)
		r.Body().WithCharset(CONTENT_CHARSET_DEFAULT)
		r.lock()
		err := r.Send()
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_RESPONSE_IS_SENDING))

		r.unlock()
		err = r.Send()
		Expect(err).To(BeNil())
	})

	It("IsSent should return boolean", func() {
		r := &FactoryResponse{}
		Expect(r.IsSent()).To(BeFalse())
		r.isSent = true
		Expect(r.IsSent()).To(BeTrue())
	})
})
