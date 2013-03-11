// Package token is a token in the EDIFACT data.

package token

import "strconv"

// Pos represents a byte position in the original input text
type Pos int

func (p Pos) Position() Pos {
  return p
}

type TokenType int

// Represents a token in the EDIFACT data.
type Token struct {
  Typ TokenType // The type of this token
  Pos Pos       // The starting position, in bytes, of this item in the input string.
  Val string    // The value of this item.
}

const (
  ERROR TokenType = iota // error occurred; value is text of error
  EOF
  COMPONENT_DELIMITER
  DATA_DELIMITER
  DECIMAL
  RELEASE_INDICATOR
  REPETITION_DELIMITER
  SEGMENT_TERMINATOR

  SEGMENT // Identifies beginning of segment and is the segment name
  TEXT
  UNA_SEGMENT
  UNA_TEXT
)

var tokens = [...]string{
  ERROR:                "Error",
  EOF:                  "EOF",
  COMPONENT_DELIMITER:  "ComponentDelimiter",
  DATA_DELIMITER:       "DataDelimiter",
  DECIMAL:              "Decimal",
  RELEASE_INDICATOR:    "ReleaseIndicator",
  REPETITION_DELIMITER: "RepetitionDelimiter",
  SEGMENT_TERMINATOR:   "SegmentTerminator",

  SEGMENT:     "Segment",
  TEXT:        "Text",
  UNA_SEGMENT: "UNASegment",
  UNA_TEXT:    "UNAText",
}

// String returns the string corresponding to the token tok.
func (tok Token) String() string {
  switch {
  case tok.Typ == EOF:
    return "EOF"
  case tok.Typ == ERROR:
    return tok.Val
    //case len(i.val) > 10:
    //  return fmt.Sprintf("%.10q...", i.val)
  }

  s := ""
  if 0 <= tok.Typ && tok.Typ < TokenType(len(tokens)) {
    s = tokens[tok.Typ]
  }
  if s == "" {
    s = "token(" + strconv.Itoa(int(tok.Typ)) + ")"
  }

  s += ": " + tok.Val

  return s
}
