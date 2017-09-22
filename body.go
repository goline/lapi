package lapi

import (
	"fmt"

	"github.com/goline/errors"
)

type Body interface {
	BodyContent
	ParserManager
}

// BodyContent handles body's content
type BodyContent interface {
	// Content returns body's content
	Content() (interface{}, error)

	// WithContent sets body's content
	WithContent(content interface{}) Body

	// ContentBytes returns body's content as bytes
	ContentBytes() ([]byte, error)

	// WithContentBytes sets body's content's bytes
	WithContentBytes(bytes []byte, v interface{}) Body

	// ContentType returns type of body's content
	ContentType() string

	// WithContentType sets body's content's type
	WithContentType(contentType string) Body

	// Charset returns character set for response
	Charset() string

	// WithCharset sets charset of response
	WithCharset(charset string) Body
}

func NewBody() Body {
	return &FactoryBody{
		ParserManager: NewParserManager(),
	}
}

type FactoryBody struct {
	ParserManager
	content      interface{}
	contentBytes []byte
	contentType  string
	charset      string
	err          error
}

func (b *FactoryBody) Content() (interface{}, error) {
	return b.content, b.err
}

func (b *FactoryBody) WithContent(content interface{}) Body {
	b.reset()

	p, ok := b.Parser(b.contentType)
	if ok == true {
		bytes, err := p.Encode(content)
		if err != nil {
			b.err = errors.New(ERR_PARSE_ENCODE_FAILURE, fmt.Sprintf("Unable to encode content. Got %s", err))
		} else {
			b.content = content
			b.contentBytes = bytes
		}
	} else {
		b.err = errors.New(ERR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", b.contentType))
	}

	return b
}

func (b *FactoryBody) ContentBytes() ([]byte, error) {
	return b.contentBytes, b.err
}

func (b *FactoryBody) WithContentBytes(bytes []byte, v interface{}) Body {
	b.reset()
	if v == nil && b.contentType == "" {
		b.contentBytes = bytes
		return b
	}

	p, ok := b.Parser(b.contentType)
	if ok == true {
		err := p.Decode(bytes, v)
		if err != nil {
			b.err = errors.New(ERR_PARSE_DECODE_FAILURE, fmt.Sprintf("Unable to decode content. Got %s", err))
		} else {
			b.content = v
			b.contentBytes = bytes
		}
	} else {
		b.err = errors.New(ERR_NO_PARSER_FOUND, fmt.Sprintf("Unable to find an appropriate parser for %s", b.contentType))
	}

	return b
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

func (b *FactoryBody) reset() {
	b.content = nil
	b.contentBytes = make([]byte, 0)
	b.err = nil
}
