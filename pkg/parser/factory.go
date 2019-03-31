package parser

import (
	"fmt"
)

const (
	PlainTplType = "plain"
)

type Factory struct {
	parsers map[string]func() Parser
}

func NewFactory() *Factory {
	return &Factory{parsers: make(map[string]func() Parser)}
}

func (f *Factory) AddParser(tplType string, parser func() Parser) {
	f.parsers[tplType] = parser
}

func (f *Factory) Create(tplType string) (Parser, error) {
	if p, ok := f.parsers[tplType]; ok {
		return p(), nil
	}

	return nil, fmt.Errorf("could not find parser by type: %s", tplType)
}
