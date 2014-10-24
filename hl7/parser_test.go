package hl7

import (
	"bytes"
	"testing"
)

var convertEscapedTests = []struct {
	in  []byte
	out []byte
}{
	{
		[]byte(`\E\\F\\R\\S\\T\\X484559\`),
		[]byte(`\|~^&HEY`),
	},
	{
		[]byte(`\X00407F`),
		[]byte{0, 0x40, 0x7F},
	},
	{
		[]byte(`\\\F\`),
		[]byte(`|`),
	},
	{
		[]byte(`\\F\\`),
		[]byte(`F`),
	},
	{
		[]byte(`\\\\\\\\\\\\`),
		[]byte(``),
	},
	{
		// all defined in the spec, but we don't convert them
		[]byte(`\C4845\ \M484950\ \H\ \N\ \Z45\`),
		[]byte(`\C4845\ \M484950\ \H\ \N\ \Z45\`),
	},
	{
		[]byte(`\unknown`),
		[]byte(`\unknown\`),
	},
	{
		[]byte(`\unknown\`),
		[]byte(`\unknown\`),
	},
}

func TestConvertEscaped(t *testing.T) {
	ls := newLexerState()
	ls.componentSeparator = '^'
	ls.fieldRepeatSeparator = '~'
	ls.escapeCharacter = '\\'
	ls.subComponentSeparator = '&'
	ls.fieldSeparator = '|'
	ps := newParserState(ls)

	for i, tt := range convertEscapedTests {
		out := ps.convertEscaped(tt.in)
		if !bytes.Equal(out, tt.out) {
			t.Fatalf("#%d: incorrect escape conversion.\nexpected:%v\ngot:     %v", i, tt.out, out)
		}
	}
}
