package ccd_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/kdar/health/ccd"
	"github.com/shurcooL/go-goon"
)

// A lot of this data is retrieved from many CCDs in practice.
// Some are made up.
var rangeTests = []struct {
	In  string
	Out ccd.ResultRanges
}{
	{"0.00-0.00", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(0),
			High:    float64p(0),
			Text:    nil,
		},
	}},
	{"13.5 - 18", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(13.5),
			High:    float64p(18),
			Text:    nil,
		},
	}},
	{"(.0-100.0)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(0),
			High:    float64p(100),
			Text:    nil,
		},
	}},
	{"(1.000-200.6)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(1),
			High:    float64p(200.6),
			Text:    nil,
		},
	}},
	{"0.00-0.20 ng/mL", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(0),
			High:    float64p(0.2),
			Text:    nil,
		},
	}},
	{"27-32 uug", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(27),
			High:    float64p(32),
			Text:    nil,
		},
	}},
	{"80-95 um3", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(80),
			High:    float64p(95),
			Text:    nil,
		},
	}},
	{"11.5-14.5 %", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(11.5),
			High:    float64p(14.5),
			Text:    nil,
		},
	}},
	{"10-39 years: 55-110 mg/dL | 40-59 years: 70-150 mg/dL | >60 years: 80-150 mg/dL | Therapeutic Target: <100 mg/dL", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  float64p(10),
			AgeHigh: float64p(39),
			Low:     float64p(55),
			High:    float64p(110),
			Text:    nil,
		}, {
			Gender:  nil,
			AgeLow:  float64p(40),
			AgeHigh: float64p(59),
			Low:     float64p(70),
			High:    float64p(150),
			Text:    nil,
		}, {
			Gender:  nil,
			AgeLow:  float64p(60),
			AgeHigh: nil,
			Low:     float64p(80),
			High:    float64p(150),
			Text:    nil,
		}, {
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    float64p(100),
			Text:    stringp("Therapeutic Target"),
		},
	}},
	{"<130 mg/dL (calc)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    float64p(130),
			Text:    nil,
		},
	}},
	{"<=534233", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    float64p(534233),
			Text:    nil,
		},
	}},
	{">2.5 ng/mL", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(2.5),
			High:    nil,
			Text:    nil,
		},
	}},
	{">=27.0", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(27),
			High:    nil,
			Text:    nil,
		},
	}},
	{"M 13-18 g/dl; F 12-16 g/dl", ccd.ResultRanges{
		{
			Gender:  stringp("M"),
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(13),
			High:    float64p(18),
			Text:    nil,
		}, {
			Gender:  stringp("F"),
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(12),
			High:    float64p(16),
			Text:    nil,
		},
	}},
	{"mg/dL", ccd.ResultRanges{}},
	{"NA", ccd.ResultRanges{}},
	{"No data", ccd.ResultRanges{}},
	{"(Clear-Mod Cloud)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    nil,
			Text:    stringp("Clear-Mod Cloud"),
		},
	}},
	{"(Negative)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    nil,
			Text:    stringp("Negative"),
		},
	}},
	{"(Negative-250)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    nil,
			Text:    stringp("Negative-250"),
		},
	}},
	{"normal: 0.29–5.11 IU/ml", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(0.29),
			High:    float64p(5.11),
			Text:    stringp("normal"),
		},
	}},
	{"Normal (3.0-4.0 cm2), mild (1.5–2.0 cm2), moderate (1.0–1.5 cm2), severe (less than 1.0 cm2)", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(3),
			High:    float64p(4),
			Text:    stringp("Normal"),
		}, {
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(1.5),
			High:    float64p(2),
			Text:    stringp("mild"),
		}, {
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(1),
			High:    float64p(1.5),
			Text:    stringp("moderate"),
		}, {
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    float64p(1),
			Text:    stringp("severe"),
		},
	}},
	{"normal: below 1.5 mg/dL", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     nil,
			High:    float64p(1.5),
			Text:    stringp("normal"),
		},
	}},
	{"normal: above 1.5 mg/dL", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(1.5),
			High:    nil,
			Text:    stringp("normal"),
		},
	}},
	{"Normal reference range 1.0-1.5; Targeted INR 2.0-3.0", ccd.ResultRanges{
		{
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(1),
			High:    float64p(1.5),
			Text:    stringp("Normal reference range"),
		}, {
			Gender:  nil,
			AgeLow:  nil,
			AgeHigh: nil,
			Low:     float64p(2),
			High:    float64p(3),
			Text:    stringp("Targeted INR"),
		},
	}},
}

func float64p(f float64) *float64 {
	return &f
}

func stringp(s string) *string {
	return &s
}

func TestResultRange(t *testing.T) {
	for _, rt := range rangeTests {
		resultRanges := ccd.ResultRanges{}
		resultRanges.Parse(rt.In)

		// Just set the OriginalText since it will always be the input.
		for i, _ := range rt.Out {
			rt.Out[i].OriginalText = rt.In
		}

		if !reflect.DeepEqual(rt.Out, resultRanges) {
			t.Fatal(spew.Printf("\ninput:    %v\nexpected: %v\ngot:      %v\n", rt.In, rt.Out, resultRanges))
		}

		_ = spew.Dump
		_ = goon.Dump
		_ = fmt.Println
	}

	//goon.Dump(rangeTests)
}

// Helps me in generating test data for the range tests.
// Make sure you verify the data afterward. This saves me some typing.
func TestGenerateResultRangeTests(t *testing.T) {
	t.Skip("not needed.")

	tmpl, err := template.New("test").Parse(`
{{ range $_, $item := . }}
{"{{ $item.In }}", ResultRanges{
{{ range $_, $x := $item.Out }}{
  Gender:  {{ if $x.Gender }}stringp("{{ $x.Gender }}"){{ else }}nil{{ end }},
  AgeLow:  {{ if $x.AgeLow }}float64p({{ $x.AgeLow }}){{ else }}nil{{ end }},
  AgeHigh: {{ if $x.AgeHigh }}float64p({{ $x.AgeHigh }}){{ else }}nil{{ end }},
  Low:     {{ if $x.Low }}float64p({{ $x.Low }}){{ else }}nil{{ end }},
  High:    {{ if $x.High }}float64p({{ $x.High }}){{ else }}nil{{ end }},
  Text:    {{ if $x.Text }}stringp("{{ $x.Text }}"){{ else }}nil{{ end }},
},{{ end }}
}},{{ end }}`)
	if err != nil {
		panic(err)
	}

	for i, rt := range rangeTests {
		resultRanges := ccd.ResultRanges{}
		resultRanges.Parse(rt.In)
		rangeTests[i].Out = resultRanges
	}

	err = tmpl.Execute(os.Stdout, rangeTests)
	if err != nil {
		panic(err)
	}
}
