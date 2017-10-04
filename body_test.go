package lapi

import (
	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Body", func() {
	It("NewBody should return an instance of Body", func() {
		Expect(NewBag()).NotTo(BeNil())
	})
})

type sampleBodyItem struct {
	Price float64 `json:"price"`
}

type sampleParserForBody struct{}

func (p *sampleParserForBody) Encode(v interface{}) ([]byte, error) {
	return make([]byte, 0), errors.New("UNABLE_TO_ENCODE", "")
}
func (p *sampleParserForBody) Decode(data []byte, v interface{}) error {
	return errors.New("UNABLE_TO_DECODE", "")
}
func (p *sampleParserForBody) ContentType() string {
	return CONTENT_TYPE_JSON
}

var _ = Describe("FactoryBody", func() {
	It("Charset should return UTF-8", func() {
		b := &FactoryBody{}
		b.charset = "UTF-8"
		Expect(b.Charset()).To(Equal("UTF-8"))
	})

	It("WithCharset should allow to set charset", func() {
		b := &FactoryBody{}
		b.WithCharset("UTF-8")
		Expect(b.charset).To(Equal("UTF-8"))
	})

	It("ContentType should return a string represents for content-type", func() {
		b := &FactoryBody{}
		b.contentType = CONTENT_TYPE_XML
		Expect(b.ContentType()).To(Equal(CONTENT_TYPE_XML))
	})

	It("WithContentType should allow to set content-type", func() {
		b := &FactoryBody{}
		Expect(b.WithContentType(CONTENT_TYPE_XML).ContentType()).To(Equal(CONTENT_TYPE_XML))
	})
})
