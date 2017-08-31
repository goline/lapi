package parser

import (
	"encoding/json"
)

type JsonParser struct{}

func (p *JsonParser) ContentType() string {
	return "application/json"
}

func (p *JsonParser) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (p *JsonParser) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
