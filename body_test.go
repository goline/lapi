package lapi

import (
	"github.com/goline/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
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

type sampleIOReader struct{}

func (r *sampleIOReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("", "")
}

type sampleIOWriter struct{}

func (r *sampleIOWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("", "")
}

type sampleBodyParser struct{}

func (p *sampleBodyParser) ContentType() string {
	return CONTENT_TYPE_JSON
}

func (p *sampleBodyParser) Decode(data []byte, v interface{}) error {
	return nil
}

func (p *sampleBodyParser) Encode(v interface{}) ([]byte, error) {
	return []byte(""), errors.New("", "")
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

	It("Read should return error code ERR_BODY_READER_MISSING", func() {
		b := &FactoryBody{}
		err := b.Read(nil)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BODY_READER_MISSING))
	})

	It("Read should return error code ERR_BODY_READER_FAILURE", func() {
		b := &FactoryBody{reader: new(sampleIOReader)}
		err := b.Read(nil)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BODY_READER_FAILURE))
	})

	It("Read should return error code ERR_NO_PARSER_FOUND", func() {
		b := &FactoryBody{reader: strings.NewReader(`{"status": true}`), ParserManager: NewParserManager()}
		err := b.Read(nil)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_NO_PARSER_FOUND))
	})

	It("Read should return error code ERR_PARSE_DECODE_FAILURE", func() {
		b := &FactoryBody{reader: strings.NewReader(`{"status": true}`), ParserManager: NewParserManager()}
		b.WithParser(new(JsonParser))
		b.WithContentType(CONTENT_TYPE_JSON)
		err := b.Read(nil)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_PARSE_DECODE_FAILURE))
	})

	It("Read should return nil", func() {
		b := &FactoryBody{reader: strings.NewReader(`{"price": 10.2}`), ParserManager: NewParserManager()}
		b.WithParser(new(JsonParser))
		b.WithContentType(CONTENT_TYPE_JSON)
		i := new(sampleBodyItem)
		err := b.Read(i)
		Expect(err).To(BeNil())
		Expect(i.Price).To(Equal(float64(10.2)))
	})

	It("Write should return nil when writing nil", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		err := b.Write(nil)
		Expect(err).To(BeNil())
	})

	It("Write should return nil when writing string", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		err := b.Write("a string")
		Expect(err).To(BeNil())
		Expect(string(b.contentBytes)).To(Equal("a string"))
	})

	It("Write should return error code ERR_NO_PARSER_FOUND", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		err := b.Write(new(sampleBodyItem))
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_NO_PARSER_FOUND))
	})

	It("Write should return error code ERR_PARSE_ENCODE_FAILURE", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		b.WithParser(new(sampleBodyParser))
		b.WithContentType(CONTENT_TYPE_JSON)
		err := b.Write(new(sampleBodyItem))
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_PARSE_ENCODE_FAILURE))
	})

	It("Write should return nil when writing bytes", func() {
		b := &FactoryBody{ParserManager: NewParserManager()}
		err := b.Write([]byte("sample string"))
		Expect(err).To(BeNil())
		Expect(string(b.contentBytes)).To(Equal("sample string"))
	})

	It("Flush should return error code ERR_BODY_WRITER_MISSING", func() {
		b := &FactoryBody{}
		err := b.Flush()
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_BODY_WRITER_MISSING))
	})
})
