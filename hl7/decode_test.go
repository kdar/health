package hl7

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/iNamik/go_lexer"
	"github.com/iNamik/go_parser"
)

var unmarshalTests = []struct {
	file string
	out  []Segment
}{
	{"simple.hl7", simple_hl7_output},
}

func TestUnmarshal(t *testing.T) {
	for i, tt := range unmarshalTests {
		fp, err := os.Open("testdata/" + tt.file)
		if err != nil {
			t.Fatalf("#%d. received error: %s", i, err)
		}
		defer fp.Close()

		data, err := ioutil.ReadAll(fp)
		if err != nil {
			t.Fatalf("#%d. received error: %s", i, err)
		}

		out, err := Unmarshal(data)
		if err != nil {
			t.Fatalf("#%d. received error: %s", i, err)
		}

		if !reflect.DeepEqual(out, tt.out) {
			t.Fatalf("#%d: mismatch\nhave: %s\nwant: %s", i, getValidGo(out), getValidGo(tt.out))
		}
	}
}

func TestUnmarshalNoError(t *testing.T) {
	filenames, err := filepath.Glob("testdata/*.hl7")
	if err != nil {
		t.Fatalf("received error: %s", err)
	}
	for _, filename := range filenames {
		fp, err := os.Open(filename)
		if err != nil {
			t.Fatalf("%s: received error: %s", filename, err)
		}
		defer fp.Close()

		data, err := ioutil.ReadAll(fp)
		if err != nil {
			t.Fatalf("%s: received error: %s", filename, err)
		}

		out, err := Unmarshal(data)
		if err != nil {
			t.Fatalf("%s: received error: %s", filename, err)
		}

		if out == nil || len(out) == 0 {
			t.Fatalf("%s: expected segments, got none", filename)
		}
	}
}

var simpleTests = []struct {
	input     []byte
	shouldErr bool
	err       string
}{
	{[]byte(""), true, "did not error on empty HL7"},
	{[]byte("BAD|^~\\&|stuff|things"), true, "did not error when given a bad header"},
	{[]byte("MSH|^~\\&|\rbadseg|field"), true, "did not error on a bad segment"},
	{[]byte("MSH\r|^~\\&"), true, "did not error on a bad header"},
	{[]byte("MSH|hellothere|field"), true, "did not error on a bad header"},
	{[]byte("\r\n\r\n\r\n\r\n\r\n"), true, "did not error on a bad HL7"},
	{
		[]byte("\x0bMSH|^~\\&|||||20071203173658|||20071203173658.98 \x0d`/usr/bin/whoami`|||XXX|\x0d\x1c\x0d"),
		true,
		"did not error on a bad HL7",
	},
	{[]byte("MSH|\r~\\&|field"), true, "did not error on invalid separator"},

	{[]byte("MSH|||||||||||||||||||"), true, "did not error on missing all separators"},
	{[]byte(`MSH|@%\|||||||||`), true, "did not error on missing one separator"},
	{[]byte(`MSH|@%\&$?!|||||||||`), true, "did not error on having too many separators"},
	{
		[]byte(`MSH|^&&^|SENDING APP|SENDING FAC|REC APP|REC FAC|20080115153000||ADT^A01^ADT_A01|0123456789|P|2.6||||AL\r`),
		true,
		"did not error on having duplicate separators",
	},

	{[]byte(`MSH$%~\&$GHH LAB\rPID$$$555-44-4444$$EVERYWOMAN%EVE%E%%%L`), false, "should not error on non-standard separators"},
	{[]byte("\r\n\r\r\r\n\r\nMSH|^~\\&|hey\r\n\n\r\n"), false, "should not error on multiple CR and NL"},
}

func TestUnmarshalSimple(t *testing.T) {
	for _, tt := range simpleTests {
		_, err := Unmarshal(tt.input)
		if tt.shouldErr && err == nil {
			t.Fatal(tt.err)
		}
		if !tt.shouldErr && err != nil {
			t.Fatalf("%v: %v", tt.err, err)
		}
	}
}

func TestMultiple(t *testing.T) {
	data := []byte("MSH|^~\\&|||1^2^3^4^^^s1&s2&s3&&~r1~r2~r3~r4~~\r\n\nPV1|1^2^3\rPV2|1^2^3\r\r\r\n\r")
	expected := []Segment{
		Segment{
			Field("MSH"),
			Field("|"),
			Field("^~\\&"),
			Field(nil),
			Field(nil),
			Repeated{
				Component{
					Field("1"),
					Field("2"),
					Field("3"),
					Field("4"),
					Field(nil),
					Field(nil),
					SubComponent{
						Field("s1"),
						Field("s2"),
						Field("s3"),
						Field(nil),
						Field(nil),
					},
				},
				Field("r1"),
				Field("r2"),
				Field("r3"),
				Field("r4"),
				Field(nil),
				Field(nil),
			},
		},
		Segment{
			Field("PV1"),
			Component{
				Field("1"),
				Field("2"),
				Field("3"),
			},
		},
		Segment{
			Field("PV2"),
			Component{
				Field("1"),
				Field("2"),
				Field("3"),
			},
		},
	}

	out, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("received error: %s", err)
	}

	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("mismatch\nhave: %s\nwant: %s", getValidGo(out), getValidGo(expected))
	}
}

func BenchmarkParser(b *testing.B) {
	fp, err := os.Open("testdata/simple.hl7")
	if err != nil {
		b.Fatalf("received error: %s", err)
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		b.Fatalf("received error: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		lexState := newLexerState()
		l := lexer.New(lexState.lexHeader, reader, 3)
		parseState := newParserState(lexState)
		p := parser.New(parseState.parse, l, 3)
		p.Next()
	}
}

func getValidGo(segments []Segment) string {
	b := &bytes.Buffer{}
	printValidGo(b, segments)
	return b.String()
}

func printValidGo(w io.Writer, segments []Segment) {
	indent := "  "
	fmt.Fprintln(w, "[]Segment{")
	for _, segment := range segments {
		printValidGo_(w, indent, segment)
	}
	fmt.Fprintln(w, "}")
}

func printValidGo_(w io.Writer, indent string, data Data) {
	switch t := data.(type) {
	case Segment:
		fmt.Fprintln(w, indent+"Segment{")
		for _, v := range t {
			printValidGo_(w, indent+"  ", v)
		}
		fmt.Fprintln(w, indent+"},")
	case Component:
		fmt.Fprintln(w, indent+"Component{")
		for _, v := range t {
			printValidGo_(w, indent+"  ", v)
		}
		fmt.Fprintln(w, indent+"},")
	case SubComponent:
		fmt.Fprintln(w, indent+"SubComponent{")
		for _, v := range t {
			printValidGo_(w, indent+"  ", v)
		}
		fmt.Fprintln(w, indent+"},")
	case Repeated:
		fmt.Fprintln(w, indent+"Repeated{")
		for _, v := range t {
			printValidGo_(w, indent+"  ", v)
		}
		fmt.Fprintln(w, indent+"},")
	case Field:
		v := "nil"
		if t != nil {
			v = `"` + t.String() + `"`
		}
		fmt.Fprintf(w, indent+"Field(%s),\n", v)
	}
}

var (
	simple_hl7_output = []Segment{
		Segment{
			Field("MSH"),
			Field("|"),
			Field(`^~\&`),
			Field("field"),
			Field(`\|~^&HEY`),
			Component{
				Field("component1"),
				Field("component2"),
			},
			Component{
				SubComponent{
					Field("subcomponent1a"),
					Field("subcomponent2a"),
				},
				SubComponent{
					Field("subcomponent1b"),
					Field("subcomponent2b"),
				},
			},
			Repeated{
				Component{
					Field("component1a"),
					Field("component2a"),
				},
				Component{
					Field("component1b"),
					Field("component2b"),
				},
			},
		},
	}
)
