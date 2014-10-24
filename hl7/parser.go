// Uses an operator-precedence parser to parse the lexer tokens
// into a data structure

package hl7

import (
	"encoding/hex"

	"github.com/iNamik/go_lexer"
	"github.com/iNamik/go_parser"
)

// opPrec is the operator precedence of the different
// separators and terminators
var opPrec = map[lexer.TokenType]int{
	tokSubComponentSeparator: 4,
	tokComponentSeparator:    3,
	tokFieldRepeatSeparator:  2,
	tokFieldSeparator:        1,
}

// parserState represents the state for the parser.
type parserState struct {
	lexState *lexerState
}

func newParserState(lexState *lexerState) *parserState {
	return &parserState{lexState: lexState}
}

func (s *parserState) parse(p parser.Parser) parser.StateFn {
	segments, err := s.parseSegments(p)
	if err != nil {
		p.Emit(err)
		return nil
	}

	p.Emit(segments)
	return nil
}

// parseSegments takes tokens from the parser and parses the segments from it.
// uses Shunting-yard algorithm in postfix notation (reverse polish)
// http://rosettacode.org/wiki/Parsing/Shunting-yard_algorithm#Go
func (s *parserState) parseSegments(p parser.Parser) ([]Segment, error) {
	var stack []*lexer.Token
	var result []*lexer.Token
	for {
		tok := p.NextToken()
		//fmt.Println(tokenTypeAsString(tok.Type()), string(tok.Bytes()))
		if tok.Type() == tokEOF {
			break
		}

		// This token is a special case because it occurs always at
		// the end of a segment and always will have last priority.
		// It doesn't fit the algorithm because it's a terminator and
		// not a binary operator.
		// When we encounter this, we just drain the rest of the stack
		// to the result, add our segment token to the result, and continue
		// with the next tokens.
		if tok.Type() == tokSegmentTerminator {
			// drain stack to result
			for len(stack) > 0 {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			result = append(result, tok)
			continue
		}

		if prec1, isOp := opPrec[tok.Type()]; isOp {
			// token is an operator
			for len(stack) > 0 {
				// consider top item on stack
				op := stack[len(stack)-1]
				if prec2, isOp := opPrec[op.Type()]; !isOp || prec1 > prec2 {
					break
				}
				// top item is an operator that needs to come off
				stack = stack[:len(stack)-1] // pop it
				result = append(result, op)  // add it to result
			}
			// push operator (the new one) to stack
			stack = append(stack, tok)
		} else { // token is an operand
			result = append(result, tok) // add operand to result
		}
	}

	// for _, t := range result {
	// 	fmt.Println(tokenTypeAsString(t.Type()), string(t.Bytes()))
	// }
	return s.process(result)
}

// process takes the tokens that are in reverse polish notiation, and
// converts them into Segments.
// http://en.wikipedia.org/wiki/Reverse_Polish_notation
func (s *parserState) process(input []*lexer.Token) ([]Segment, error) {
	var segments []Segment
	var stack []Data

	for i := range input {
		tok := input[i]

		switch tok.Type() {
		case tokSubComponentSeparator:
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if v, ok := operand1.(SubComponent); ok {
				v.Append(operand2)
				stack = append(stack, v)
			} else {
				stack = append(stack, SubComponent{operand1.(Field), operand2.(Field)})
			}
		case tokComponentSeparator:
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if v, ok := operand1.(Component); ok {
				v.Append(operand2)
				stack = append(stack, v)
			} else {
				stack = append(stack, Component{operand1, operand2})
			}
		case tokFieldRepeatSeparator:
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if v, ok := operand1.(Repeated); ok {
				v.Append(operand2)
				stack = append(stack, v)
			} else {
				stack = append(stack, Repeated{operand1, operand2})
			}
		case tokFieldSeparator:
		case tokSegmentTerminator:
			var segment Segment
			segment.Append(stack...)
			segments = append(segments, segment)
			stack = []Data{}
		default:
			if tok.Type() == tokSegmentName || tok.Type() == tokSeparators {
				stack = append(stack, Field(tok.Bytes()))
			} else {
				stack = append(stack, Field(s.convertEscaped(tok.Bytes())))
			}
		}
	}

	return segments, nil
}

// convertEscaped replaces escape sequences in the bytes passed
func (s *parserState) convertEscaped(e []byte) []byte {
	var final []byte

	for i := 0; i < len(e); i++ {
		if rune(e[i]) == s.lexState.escapeCharacter {
			// need at least more characters in the escape sequence
			if i+1 >= len(e) {
				continue
			}

			x := i + 1
			// find the end of the escape sequence
			for ; x < len(e) && rune(e[x]) != s.lexState.escapeCharacter; x++ {
			}

			// if we get this case, that means they put two escape characters
			// back-to-back. just skip it.
			if i+1 == x {
				i++
				continue
			}

			switch sequence := e[i+1 : x]; sequence[0] {
			case 'E':
				final = append(final, byte(s.lexState.escapeCharacter))
			case 'F':
				final = append(final, byte(s.lexState.fieldSeparator))
			case 'R':
				final = append(final, byte(s.lexState.fieldRepeatSeparator))
			case 'S':
				final = append(final, byte(s.lexState.componentSeparator))
			case 'T':
				final = append(final, byte(s.lexState.subComponentSeparator))
			case 'X':
				sequence = sequence[1:]
				out := make([]byte, len(sequence)/2)
				hex.Decode(out, sequence)
				final = append(final, out...)
			default:
				final = append(final, byte(s.lexState.escapeCharacter))
				final = append(final, sequence...)
				final = append(final, byte(s.lexState.escapeCharacter))
			}
			i = x
			continue
		}

		final = append(final, e[i])
	}

	return final
}
