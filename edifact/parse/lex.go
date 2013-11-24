// A lexer to lex EDIFACT data.

package parse

import (
	"fmt"
	"github.com/kdar/health/edifact/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

// our own defined eof so we can track it
const eof = -1

// default values for the UNA segment. if no
// una segment is specified, then these are used
const (
	UNA_SEGMENT_NAME         = "UNA"
	UNA_COMPONENT_DELIMITER  = ':'
	UNA_DATA_DELIMITER       = '+'
	UNA_DECIMAL              = '.'
	UNA_RELEASE_INDICATOR    = '?'
	UNA_REPETITION_DELIMITER = ' ' // This is really "reserved" in the spec.
	UNA_SEGMENT_TERMINATOR   = '\''
)

const (
	COMPONENT_DELIMITER_POS = iota
	DATA_DELIMITER_POS
	DECIMAL_POS
	RELEASE_INDICATOR_POS
	REPETITION_DELIMITER_POS
	SEGMENT_TERMINATOR_POS
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name                string           // the name of the input; used only for error reports
	input               string           // the string being scanned
	componentDelimiter  rune             // separates components within data
	dataDelimiter       rune             // separates data within a segment
	decimal             rune             // the character used to signify decimal numbers
	releaseIndicator    rune             // escapes the next character
	repetitionDelimiter rune             // used as repetition in some other specs. We will use it as repitition.
	segmentTerminator   rune             // terminates each segment
	state               stateFn          // the next lexing function to enter
	pos                 token.Pos        // current position in the input
	start               token.Pos        // start position of this item
	width               token.Pos        // width of last rune read from input
	lastPos             token.Pos        // position of most recent item returned by nextItem
	tokens              chan token.Token // channel of scanned tokens

	// rare special case
	foundQuote rune
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = token.Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// returns the rune at the particular position.
func (l *lexer) at(pos token.Pos) (rune, int) {
	if int(pos) >= len(l.input) {
		return eof, 0
	}
	r, w := utf8.DecodeRuneInString(l.input[pos:])
	return r, w
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t token.TokenType) {
	l.tokens <- token.Token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token.Token{token.ERROR, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextToken() token.Token {
	tok := <-l.tokens
	l.lastPos = tok.Pos
	return tok
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:                name,
		input:               input,
		componentDelimiter:  UNA_COMPONENT_DELIMITER,
		dataDelimiter:       UNA_DATA_DELIMITER,
		releaseIndicator:    UNA_RELEASE_INDICATOR,
		repetitionDelimiter: UNA_REPETITION_DELIMITER,
		segmentTerminator:   UNA_SEGMENT_TERMINATOR,
		tokens:              make(chan token.Token, 2),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexSegment; l.state != nil; {
		l.state = l.state(l)
	}
}

// lex a segment. if this is a UNA segment then we
// need to parse that to determine how to parse the
// rest of the message. we know it's a segment name if
// it starts with capital letters.
func lexSegment(l *lexer) stateFn {
	if /*l.pos == 0 && */ strings.HasPrefix(l.input[l.pos:], UNA_SEGMENT_NAME) {
		return lexUNASegment
	}

	r := l.peek()

	if r == '\n' {
		return lexBeginningNewlines
	}

	if isUpper(r) {
		return lexSegmentName
	}

	l.emit(token.EOF)
	return nil
}

// lex any amount of newlines before a segment.
// some companies do this for some reason
func lexBeginningNewlines(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '\n':
			// ignore
			l.start = l.pos
		case isAlphaNumeric(r) && isUpper(r):
			l.backup()
			return lexSegment
		}
	}
	return nil
}

// parse the UNA segment and retrieve all
// our delimiters and settings.
func lexUNASegment(l *lexer) stateFn {
	l.pos += token.Pos(len(UNA_SEGMENT_NAME))
	l.emit(token.UNA_SEGMENT)
	//l.emit(token.SEGMENT)

	// read the next 6 runes, because they are the
	// data for the UNA segment that we need
	// to lex the rest
	for x := 0; x < 6; x++ {
		r := l.next()
		if r == eof {
			return l.errorf("found eof while reading UNA header")
		}

		switch x {
		case COMPONENT_DELIMITER_POS:
			l.componentDelimiter = r
		case DATA_DELIMITER_POS:
			l.dataDelimiter = r
		case DECIMAL_POS:
			l.decimal = r
		case RELEASE_INDICATOR_POS:
			l.releaseIndicator = r
		case REPETITION_DELIMITER_POS:
			l.repetitionDelimiter = r
		case SEGMENT_TERMINATOR_POS:
			l.segmentTerminator = r
		}
	}

	l.emit(token.UNA_TEXT)
	//l.emit(token.TEXT)

	return lexSegment
}

// les the segment name. usually this is just
// three uppercase letters.
func lexSegmentName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == l.dataDelimiter:
			l.backup()
			l.emit(token.SEGMENT)
			return lexDataDelimiter
		case isAlphaNumeric(r) && isUpper(r):
			// absorb.
		case r == eof:
			return l.errorf("found eof while reading segment name")
		default:
			return l.errorf("unknown character found while reading segment name: %q", r)
		}
	}

	return lexSegment
}

// lex the data delimiter
func lexDataDelimiter(l *lexer) stateFn {
	l.next() // we already know this is the data delimiter
	l.emit(token.DATA_DELIMITER)
	return lexData
}

// lex a data section. a data section can have
// components, repetitions, and texts in it.
func lexData(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == l.dataDelimiter:
			l.backup()
			l.emit(token.TEXT)
			return lexDataDelimiter
		case r == l.segmentTerminator:
			// now this might sound retarded (because it is) but some
			// companies (::cough:: relayhealth) do not escape the
			// quotations inside when it is used as a delimiter also.
			// what this does is if it detects a quote character, and
			// this isn't the end of the input or the data following isn't
			// the start of another segment, then just absorb it as if
			// it were token.TEXT.
			// note: this does not cover the case where if they don't
			// escape other delimiters. but i have not seen this yet.
			if l.foundQuote == 0 && isQuote(r) && int(l.pos) < len(l.input) {
				isTerm := true
				p := l.pos

				// test if the next 3 runes are upper case
				for x := 0; x < 3 && int(p) < len(l.input); x++ {
					subr, w := l.at(p)
					isTerm = isTerm && isUpper(subr)
					p += token.Pos(w)
				}

				// check to see if the 4th rune is a data delimiter
				if int(p) < len(l.input) {
					subr, _ := l.at(p)
					isTerm = isTerm && subr == l.dataDelimiter

					if !isTerm {
						l.foundQuote = r
					}
				}
			}

			if l.foundQuote == 0 {
				l.backup()
				l.emit(token.TEXT)
				return lexSegmentTerminator
			} else {
				// we absorb the quote
				l.foundQuote = 0
			}
		case r == l.releaseIndicator:
			// skip to next character since it is escaped
			l.next()
		case r == l.componentDelimiter:
			l.backup()
			l.emit(token.TEXT)
			return lexComponentDelimiter
		case r == l.repetitionDelimiter:
			l.backup()
			l.emit(token.TEXT)
			return lexRepetitionDelimiter
		case r == eof:
			return l.errorf("found eof while reading data")
		default:
			// absorb
		}
	}

	return lexSegment
}

// lex the component delimiter
func lexComponentDelimiter(l *lexer) stateFn {
	l.next() // we already know this is the component delimiter
	l.emit(token.COMPONENT_DELIMITER)
	return lexData
}

// lex the repetition delimiter
func lexRepetitionDelimiter(l *lexer) stateFn {
	l.next() // we already know this is the repetition delimiter
	l.emit(token.REPETITION_DELIMITER)
	return lexData
}

// les the segment terminator
func lexSegmentTerminator(l *lexer) stateFn {
	l.next() // we already know this is the segment terminator
	l.emit(token.SEGMENT_TERMINATOR)
	return lexSegment
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isUpper reports whether r is all upper or not.
func isUpper(r rune) bool {
	return unicode.IsUpper(r)
}

// isQuote reports whether r is a quote or not.
func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}
