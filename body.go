package lapi

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/goline/errors"
)

type Body interface {
	BodyRW
	BodyContent
	BodyFlusher
	ParserManager
}

// BodyContent handles body's content
type BodyContent interface {
	// ContentType returns type of body's content
	ContentType() string

	// WithContentType sets body's content's type
	WithContentType(contentType string) Body

	// Charset returns character set for response
	Charset() string

	// WithCharset sets charset of response
	WithCharset(charset string) Body
}

type BodyRW interface {
	BodyReader
	BodyWriter
}

type BodyReader interface {
	// Read gets value into input
	Read(input interface{}) error
}

type BodyWriter interface {
	// Writes puts output into writer
	Write(output interface{}) error
}

type BodyFlusher interface {
	// Flush sends output to writer
	Flush() error
}

func NewBody(reader io.ReadCloser, writer io.Writer) Body {
	return &FactoryBody{
		reader: reader,
		writer: writer,

		ParserManager: NewParserManager(),
	}
}

type FactoryBody struct {
	ParserManager
	reader       io.ReadCloser
	writer       io.Writer
	contentBytes []byte
	contentType  string
	charset      string
}

func (b *FactoryBody) Read(input interface{}) error {
	if b.reader == nil {
		return errors.New(ERR_BODY_READER_MISSING, "There is no readers specified")
	}
	defer b.reader.Close()

	bytes, err := ioutil.ReadAll(b.reader)
	if err != nil {
		return errors.New(ERR_BODY_READER_FAILURE, fmt.Sprintf("Unable to read resource. Got %s", err.Error()))
	}

	p, ok := b.Parser(b.contentType)
	if ok == false {
		return errors.New(ERR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", b.contentType))
	}

	err = p.Decode(bytes, input)
	if err != nil {
		return errors.New(ERR_PARSE_DECODE_FAILURE, fmt.Sprintf("Unable to decode content. Got %s", err))
	}

	return nil
}

func (b *FactoryBody) Write(output interface{}) error {
	if b.writer == nil {
		return errors.New(ERR_BODY_WRITER_MISSING, "Writer must not be nil")
	}

	p, ok := b.Parser(b.contentType)
	if ok == false {
		return errors.New(ERR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", b.contentType))
	}

	bytes, err := p.Encode(output)
	if err != nil {
		return errors.New(ERR_PARSE_ENCODE_FAILURE, fmt.Sprintf("Unable to encode content. Got %s", err))
	}
	b.contentBytes = bytes

	return nil
}

func (b *FactoryBody) Flush() error {
	_, err := b.writer.Write(b.contentBytes)
	if err != nil {
		return errors.New(ERR_BODY_WRITER_FAILURE, fmt.Sprintf("Unable to write output. Got %s", err.Error()))
	}

	return nil
}

func (b *FactoryBody) ContentType() string {
	return b.contentType
}

func (b *FactoryBody) WithContentType(contentType string) Body {
	b.contentType = contentType
	return b
}

func (b *FactoryBody) Charset() string {
	return b.charset
}

func (b *FactoryBody) WithCharset(charset string) Body {
	b.charset = charset
	return b
}
