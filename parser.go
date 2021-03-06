package lapi

import (
	"encoding/json"
)

type Parser interface {
	// ContentType returns supported content-type
	ContentType() string

	// Decode decodes data into a specific value
	Decode(data []byte, v interface{}) error

	// Encode encodes value to data
	Encode(v interface{}) ([]byte, error)
}

type ParserManager interface {
	// Parser returns an appropriate parser
	Parser(contentType string) (Parser, bool)

	// WithParser registers a parser
	WithParser(parser Parser) ParserManager
}

func NewParserManager() ParserManager {
	return &FactoryParserManager{make(map[string]Parser)}
}

type FactoryParserManager struct {
	parsers map[string]Parser
}

func (pm *FactoryParserManager) Parser(contentType string) (Parser, bool) {
	parser, ok := pm.parsers[contentType]
	return parser, ok
}

func (pm *FactoryParserManager) WithParser(parser Parser) ParserManager {
	pm.parsers[parser.ContentType()] = parser
	return pm
}

type JsonParser struct{}

func (p *JsonParser) ContentType() string {
	return CONTENT_TYPE_JSON
}

func (p *JsonParser) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (p *JsonParser) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
