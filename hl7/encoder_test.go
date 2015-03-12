package hl7

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	simple_hl7_input = []Segment{
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

	sample_hl7_input = []Segment{
		Segment{
			Field("MSH"),
			Field("|"),
			Field(`^~\&`),
			Field(`EP^IC`),
			Field("EPICADT"),
			Field("SMS"),
			Field("SMSADT"),
			Field("199912271408"),
			Field("CHARRIS"),
			Component{
				Field("ADT"),
				Field("A04"),
			},
			Field("1817457"),
			Field("D"),
			Field("2.5"),
			Field(""),
		},
		Segment{
			Field("PID"),
			Field(""),
			Component{
				Field("0493575"),
				Field(""),
				Field(""),
				Field("2"),
				Field("ID 1"),
			},
			Field("454721"),
			Field(""),
			Component{
				Field("DOE"),
				Field("JOHN"),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
			},
			Component{
				Field("DOE"),
				Field("JOHN"),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
			},
			Field("19480203"),
			Field("M"),
			Field(""),
			Field("B"),
			Component{
				Field("254 MYSTREET AVE"),
				Field(""),
				Field("MYTOWN"),
				Field("OH"),
				Field("44123"),
				Field("USA"),
			},
			Field(""),
			Field("(216)123-4567"),
			Field(""),
			Field(""),
			Field("M"),
			Field("NON"),
			Repeated{
				Field("400003403"),
				Field("1129086"),
			},
			Field(""),
		},
		Segment{
			Field("NK1"),
			Field(""),
			Component{
				Field("ROE"),
				Field("MARIE"),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
			},
			Field("SPO"),
			Field(""),
			Field("(216)123-4567"),
			Field(""),
			Field("EC"),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
		},
		Segment{
			Field("PV1"),
			Field(""),
			Field("O"),
			Component{
				Repeated{
					Field("168 "),
					Field("219"),
					Field("C"),
					Field("PMA"),
				},
				Field(""),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
			},
			Field(""),
			Field(""),
			Field(""),
			Component{
				Field("277"),
				Field("ALLEN MYLASTNAME"),
				Field("BONNIE"),
				Field(""),
				Field(""),
				Field(""),
				Field(""),
			},
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(" "),
			Field(""),
			Field("2688684"),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field("199912271408"),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(""),
			Field(`002376853"`),
		},
	}
)

var marshalTests = []struct {
	file string
	in   []Segment
}{
	{"simple_nohex.hl7", simple_hl7_input},
	{"sample.hl7", sample_hl7_input},
}

func TestMarshalFiles(t *testing.T) {
	for i, tt := range marshalTests {
		fp, err := os.Open("testdata/" + tt.file)
		if err != nil {
			t.Fatalf("#%d. received error: %s", i, err)
		}
		defer fp.Close()

		expected, err := ioutil.ReadAll(fp)
		if err != nil {
			t.Fatalf("#%d. received error: %s", i, err)
		}

		actual, err := Marshal(tt.in)
		if err != nil {
			t.Fatalf("#%d. received error: %s", i, err)
		}

		if !bytes.Equal(actual, expected) {
			if len(actual) != len(expected) {
				ioutil.WriteFile(fmt.Sprintf("/tmp/bad%d.hl7", i), actual, 0777)
				t.Fatalf("#%d: length mismatch\nhave: %d\nwant: %d", i, len(actual), len(expected))
			} else {

				loc := -1
				l := len(actual)
				if l > len(expected) {
					l = len(expected)
				}

				for j := 0; j < l; j++ {
					if actual[j] != expected[j] {
						loc = j
						break
					}
				}

				sloc := loc - 5
				if sloc < 0 {
					sloc = 0
				}

				eloc := loc + 5
				if eloc >= l {
					eloc = l
				}

				t.Fatalf("#%d: mismatch at byte %d\nhave: %s (%v)\nwant: %s (%v)", i, loc, string(actual[sloc:eloc]), actual[sloc:eloc], string(expected[sloc:eloc]), expected[sloc:eloc])
			}
		}
	}
}
