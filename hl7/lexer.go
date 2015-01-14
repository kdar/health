package hl7

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/iNamik/go_lexer"
	"github.com/iNamik/go_lexer/rangeutil"
)

const (
	tokEOF lexer.TokenType = lexer.TokenTypeEOF
	tokNIL                 = lexer.TokenTypeEOF + iota
	tokError
	tokSegmentName
	tokSeparators
	tokSegmentTerminator
	tokFieldSeparator
	tokComponentSeparator
	tokFieldRepeatSeparator
	tokSubComponentSeparator
	tokField
	//tokEscaped
)

const (
	fieldSeparatorPos = iota
	componentSeparatorPos
	fieldRepeatSeparatorPos
	escapeCharacterPos
	subComponentSeparatorPos
)

var (
	bytesUpperChars     = rangeutil.RangeToBytes("A-Z0-9")
	bytesSegmentGarbage = []byte{'\r', '\n'}
)

// tokenTypeAsString converts a token type to a string.
func tokenTypeAsString(t lexer.TokenType) string {
	var typeString string

	switch t {
	case tokEOF:
		typeString = "tokEOF"
	case tokNIL:
		typeString = "tokNIL"
	case tokError:
		typeString = "tokError"
	case tokSegmentName:
		typeString = "tokSegmentName"
	case tokSeparators:
		typeString = "tokSeparators"
	case tokSegmentTerminator:
		typeString = "tokSegmentTerminator"
	case tokFieldSeparator:
		typeString = "tokFieldSeparator"
	case tokComponentSeparator:
		typeString = "tokComponentSeparator"
	case tokFieldRepeatSeparator:
		typeString = "tokFieldRepeatSeparator"
	case tokSubComponentSeparator:
		typeString = "tokSubComponentSeparator"
	case tokField:
		typeString = "tokField"
	//case tokEscaped:
	//	typeString = "tokEscaped"

	default:
		typeString = strconv.Itoa(int(t))
	}

	return typeString
}

// lexerState represents the state for the lexer.
type lexerState struct {
	segmentTerminator     rune
	fieldSeparator        rune
	componentSeparator    rune
	fieldRepeatSeparator  rune
	escapeCharacter       rune
	subComponentSeparator rune
	err                   error

	lexFieldSeparator        func(lexer.Lexer) lexer.StateFn
	lexComponentSeparator    func(lexer.Lexer) lexer.StateFn
	lexFieldRepeatSeparator  func(lexer.Lexer) lexer.StateFn
	lexSubComponentSeparator func(lexer.Lexer) lexer.StateFn
}

// newLexerState returns the state that helps in lexing
func newLexerState() *lexerState {
	s := &lexerState{
		segmentTerminator: '\r',
	}

	// make a function to error if these functions are called
	// before we find the message header
	f := func(l lexer.Lexer) lexer.StateFn {
		s.err = errors.New("did not find message header")
		l.EmitToken(tokError)
		l.EmitEOF()
		return nil
	}
	s.lexFieldSeparator = f
	s.lexComponentSeparator = f
	s.lexFieldRepeatSeparator = f
	s.lexSubComponentSeparator = f

	return s
}

// lexHeader scans for the HL7 message header
func (s *lexerState) lexHeader(l lexer.Lexer) lexer.StateFn {
	l.MatchZeroOrMoreRunes([]rune{s.segmentTerminator, '\n'})

	// if l.MatchEOF() {
	// 	l.EmitEOF()
	// 	return nil
	// }

	if matchRunes(l, []rune("MSH")) {
		l.EmitTokenWithBytes(tokSegmentName)
	} else {
		s.err = errors.New("could not find message header")
		l.EmitToken(tokError)
		return nil
	}

	return s.lexHeaderSeparators
}

// lexHeaderSeparators scans for the separators used to parse the
// rest of the HL7 message.
func (s *lexerState) lexHeaderSeparators(l lexer.Lexer) lexer.StateFn {
	for i := 0; i <= subComponentSeparatorPos-fieldSeparatorPos; i++ {
		r := l.NextRune()
		if r == lexer.RuneEOF {
			s.err = errors.New("found eof while reading message header")
			l.EmitToken(tokError)
			return nil
		} else if i != fieldSeparatorPos && r == s.fieldSeparator {
			s.err = fmt.Errorf("missing %d separators", 5-i)
			l.EmitToken(tokError)
			return nil
		}

		switch i {
		case fieldSeparatorPos:
			s.fieldSeparator = r
			l.EmitTokenWithBytes(tokField)
		case componentSeparatorPos:
			s.componentSeparator = r
		case fieldRepeatSeparatorPos:
			s.fieldRepeatSeparator = r
		case escapeCharacterPos:
			s.escapeCharacter = r
		case subComponentSeparatorPos:
			s.subComponentSeparator = r
		}
	}

	s.lexFieldSeparator = s.getLexSeparator("field", s.fieldSeparator, tokFieldSeparator)
	s.lexComponentSeparator = s.getLexSeparator("component", s.componentSeparator, tokComponentSeparator)
	s.lexFieldRepeatSeparator = s.getLexSeparator("field repeat", s.fieldRepeatSeparator, tokFieldRepeatSeparator)
	s.lexSubComponentSeparator = s.getLexSeparator("sub component", s.subComponentSeparator, tokSubComponentSeparator)

	l.EmitTokenWithBytes(tokSeparators)

	return s.lexFieldSeparator
}

// lexSegment scans for a HL7 segment.
func (s *lexerState) lexSegment(l lexer.Lexer) lexer.StateFn {
	if l.MatchEOF() {
		l.EmitEOF()
		return nil // We're done here
	}

	if !l.MatchMinMaxBytes(bytesUpperChars, 3, 3) {
		s.err = fmt.Errorf("incorrect segment name, found: %c%c%c", l.PeekRune(0), l.PeekRune(1), l.PeekRune(2))
		l.EmitToken(tokError)
		return nil
	}

	l.EmitTokenWithBytes(tokSegmentName)
	return s.lexFieldSeparator
}

// getLexSeparator is a function generator that returns a function
// that scans for a particular separator. This is used to prevent
// a lot of code duplication for all the separators/terminators.
func (s *lexerState) getLexSeparator(typ string, sep rune, tok lexer.TokenType) lexer.StateFn {
	return func(l lexer.Lexer) lexer.StateFn {
		r := l.NextRune()
		if r != sep {
			s.err = fmt.Errorf("expected a %s separator '%c'. got: '%c'", typ, sep, r)
			l.EmitToken(tokError)
			return nil
		}

		l.EmitTokenWithBytes(tok)
		return s.lexField
	}
}

// lexField scans for a HL7 field.
func (s *lexerState) lexField(l lexer.Lexer) lexer.StateFn {
	for {
		switch r := l.NextRune(); {
		case r == s.fieldSeparator:
			l.BackupRune()
			l.EmitTokenWithBytes(tokField)
			return s.lexFieldSeparator
		case r == s.componentSeparator:
			l.BackupRune()
			l.EmitTokenWithBytes(tokField)
			return s.lexComponentSeparator
		case r == s.fieldRepeatSeparator:
			l.BackupRune()
			l.EmitTokenWithBytes(tokField)
			return s.lexFieldRepeatSeparator
		// case r == s.escapeCharacter:
		// 	l.NonMatchOneOrMoreRunes([]rune{s.escapeCharacter})
		// 	l.NextRune()
		// 	l.EmitTokenWithBytes(tokEscaped)
		// 	return s.lexField
		case r == s.subComponentSeparator:
			l.BackupRune()
			l.EmitTokenWithBytes(tokField)
			return s.lexSubComponentSeparator
		case r == s.segmentTerminator || r == '\n':
			l.BackupRune()
			l.EmitTokenWithBytes(tokField)
			l.NewLine()
			l.NextRune()
			l.EmitTokenWithBytes(tokSegmentTerminator)
			// ignore any multiple \r and \n
			l.MatchOneOrMoreBytes(bytesSegmentGarbage)
			l.IgnoreToken()
			return s.lexSegment
		case r == lexer.RuneEOF:
			//l.BackupRune()
			l.EmitTokenWithBytes(tokField)
			l.EmitToken(tokSegmentTerminator)
			l.EmitEOF()
			return nil
		default:
			// absorb
		}
	}

	l.EmitEOF()
	return nil
}

// matchRunes consumes a specified run of matching runes. Every rune
// must match sequentially.
func matchRunes(l lexer.Lexer, match []rune) bool {
	i := 0
	return l.MatchMinMaxFunc(func(r rune) bool {
		truth := match[i] == r
		i++
		return truth
	}, len(match), len(match))
}
