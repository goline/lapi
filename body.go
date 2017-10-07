package lapi

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

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
	Read(input interface{}) errors.Error
}

type BodyWriter interface {
	// Writes puts output into writer
	Write(output interface{}) errors.Error
}

type BodyFlusher interface {
	// Flush sends output to writer
	Flush() errors.Error
}

func NewBody(reader io.Reader, writer io.Writer) Body {
	return &FactoryBody{
		reader: reader,
		writer: writer,

		ParserManager: NewParserManager(),
	}
}

type FactoryBody struct {
	ParserManager
	reader       io.Reader
	writer       io.Writer
	contentBytes []byte
	contentType  string
	charset      string
}

func (b *FactoryBody) Read(input interface{}) errors.Error {
	if b.reader == nil {
		return errors.New(ERR_BODY_READER_MISSING, "There is no readers specified")
	}
	defer func() {
		if c, ok := b.reader.(io.Closer); ok == true {
			c.Close()
		}
	}()

	bytes, err := ioutil.ReadAll(b.reader)
	if err != nil {
		return errors.New(ERR_BODY_READER_FAILURE, "Unable to read resource.").WithDebug(err.Error())
	}

	p, ok := b.Parser(b.contentType)
	if ok == false {
		return errors.New(ERR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", b.contentType))
	}

	err = p.Decode(bytes, input)
	if err != nil {
		return errors.New(ERR_PARSE_DECODE_FAILURE, "Unable to decode content").WithDebug(err.Error())
	}

	return nil
}

func (b *FactoryBody) Write(output interface{}) errors.Error {
	if output == nil {
		return nil
	}

	t := reflect.TypeOf(output)
	switch t.Kind() {
	case reflect.Slice:
		if t.String() == "[]uint8" {
			// receives []byte
			b.contentBytes = reflect.ValueOf(output).Bytes()
			return nil
		}
	case reflect.String:
		s := output.(string)
		b.contentBytes = []byte(s)
	default:
		p, ok := b.Parser(b.contentType)
		if ok == false {
			return errors.New(ERR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", b.contentType))
		}

		bytes, err := p.Encode(output)
		if err != nil {
			return errors.New(ERR_PARSE_ENCODE_FAILURE, "Unable to encode content").WithDebug(err.Error())
		}
		b.contentBytes = bytes
	}

	return nil
}

func (b *FactoryBody) Flush() errors.Error {
	if b.writer == nil {
		return errors.New(ERR_BODY_WRITER_MISSING, "Writer must not be nil")
	}

	_, err := b.writer.Write(b.contentBytes)
	if err != nil {
		return errors.New(ERR_BODY_WRITER_FAILURE, "Unable to write output").WithDebug(err.Error())
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
