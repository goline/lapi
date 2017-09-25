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

	It("Content should return bytes value", func() {
		b := &FactoryBody{}
		Expect(b.content).To(BeNil())
		Expect(b.err).To(BeNil())

		b.content = "a_string"
		Expect(b.content).To(Equal("a_string"))
		Expect(b.err).To(BeNil())
	})

	It("WithContent should return error code ERR_NO_PARSER_FOUND", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		b.WithContent("a_string")
		Expect(b.err).NotTo(BeNil())
		Expect(b.err.(errors.Error).Code()).To(Equal(ERR_NO_PARSER_FOUND))
	})

	It("WithContent should allow to set content", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		b.WithParser(new(JsonParser))
		b.WithContentType(CONTENT_TYPE_JSON)
		b.WithContent(&sampleBodyItem{Price: 5.67})
		Expect(b.err).To(BeNil())
		Expect(string(b.contentBytes)).To(Equal("{\"price\":5.67}"))
	})

	It("WithContent should make an error when parsing invalid content", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		b.WithParser(new(sampleParserForBody))
		b.WithContentType(CONTENT_TYPE_JSON)
		b.WithContent("a_string")
		Expect(b.err).NotTo(BeNil())
		Expect(b.err.(errors.Error).Code()).To(Equal(ERR_PARSE_ENCODE_FAILURE))
	})

	It("ContentBytes should return bytes equivalent to a_string", func() {
		b := &FactoryBody{}
		b.contentBytes = []byte("a_string")

		bb, err := b.ContentBytes()
		Expect(err).To(BeNil())
		Expect(bb).To(Equal([]byte("a_string")))
	})

	It("WithContentBytes should allow to set content bytes", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		b.WithContentBytes([]byte("a_string"), nil)
		Expect(b.err).To(BeNil())
		Expect(b.contentBytes).NotTo(BeNil())
		Expect(b.content).To(BeNil())

		b.WithContentType(CONTENT_TYPE_XML)
		b.WithContentBytes([]byte("a_string"), nil)
		Expect(b.err).NotTo(BeNil())

		v := &sampleBodyItem{}
		b.WithParser(new(JsonParser))
		b.WithContentType(CONTENT_TYPE_JSON)
		b.WithContentBytes([]byte("{\"price\":5.67}"), v)
		Expect(b.err).To(BeNil())
		Expect(v.Price).To(Equal(5.67))
	})

	It("WithContentBytes should return error code ERR_PARSE_DECODE_FAILURE", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		b.WithParser(new(sampleParserForBody))
		b.WithContentType(CONTENT_TYPE_JSON)
		b.WithContentBytes([]byte("a_string"), nil)
		Expect(b.err).NotTo(BeNil())
		Expect(b.err.(errors.Error).Code()).To(Equal(ERR_PARSE_DECODE_FAILURE))
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
