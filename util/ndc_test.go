package util

import (
	"testing"
)

var idNDCFmtTests = []struct {
	in  string
	out string
}{
	{"000406-0522-05", "6-4-2"},
	{"054868-5338-*3", "6-4-2"},
	{"076564-7654-5", "6-4-1"},
	{"045983-295-86", "6-3-2"},
	{"059372-837-6", "6-3-1"},
	{"45803-1234-34", "5-4-2"},
	{"12345-6789-0", "5-4-1"},
	{"60951-700-85", "5-3-2"},
	{"0591-0933-01", "4-4-2"},
	{"000406052201", ""},
}

func TestIdentifyNDCFormat(t *testing.T) {
	for _, tt := range idNDCFmtTests {
		out := IdentifyNDCFormat(tt.in)
		if out != tt.out {
			t.Fatalf("\ninput:    %v\noutput:   %v\nexpected: %v", tt.in, out, tt.out)
		}
	}
}

var normalizeNDCTests = []struct {
	in  string
	out string
}{
	{"000406-0522-05", "00406052205"},
	{"054868-5338-*3", "54868533803"},
	{"076564-7654-5", "76564765405"},
	{"045983-295-86", "45983029586"},
	{"059372-837-6", "59372083706"},
	{"045803-1234-34", "45803123434"},
	{"12345-6789-0", "12345678900"},
	{"60951-700-85", "60951070085"},
	{"0591-0933-01", "00591093301"},
	{"000406052201", "00406052201"},
	{"a000406052201", ""},
}

func TestNormalizeNDC(t *testing.T) {
	for _, tt := range normalizeNDCTests {
		out, err := NormalizeNDC(tt.in)
		if out != tt.out {
			t.Fatalf("\ninput:    %v\noutput:   %v\nexpected: %v\nerror:    %v", tt.in, out, tt.out, err)
		}
	}
}
