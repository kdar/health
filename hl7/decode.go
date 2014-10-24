package hl7

import (
	"bytes"
	"fmt"

	"github.com/iNamik/go_lexer"
	"github.com/iNamik/go_parser"
)

// Unmarshal takes the bytes passed and returns
// the segments of the hl7 message.
func Unmarshal(b []byte) ([]Segment, error) {
	reader := bytes.NewReader(b)
	lexState := newLexerState()
	l := lexer.New(lexState.lexHeader, reader, 3)
	parseState := newParserState(lexState)
	p := parser.New(parseState.parse, l, 3)

	val := p.Next()
	switch t := val.(type) {
	case error:
		return nil, t
	case []Segment:
		return t, nil
	}

	return nil, fmt.Errorf("received unknown type")
}
