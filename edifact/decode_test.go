package edifact

import (
  //"fmt"
  "reflect"
  "testing"
)

var (
  CRAZY_OUT1 = Values{
    Header{"UNA", ":+./*'"},
    Values{
      "UIB",
      Values{"UNOA", "0"},
      "",
      Values{"", "0"},
      "", "",
      Values{"Per-Se", "ZZZ"},
      Values{"Per-Se", "ZZZ"},
      Values{"20130222", "065442"},
    },
    Values{
      "UIH",
      Values{"SCRIPT", "008", "001", "ERROR"},
    },
    Values{
      "STS",
      "900", "007",
      "Missing 'Request='",
    },
    Values{
      "UIT", "", "3",
    },
    Values{
      "UIZ", "", "1",
    },
  }
)

const (
  // crazy freaking case from RelayHealth where they
  // use ' as a segment terminator, but fail to escape
  // it inside the text. You can see this where it says
  // "Missing 'Request='"
  CRAZY_IN1 = `UNA:+./*'UIB+UNOA:0++:0+++Per-Se:ZZZ+Per-Se:ZZZ+20130222:065442'UIH+SCRIPT:008:001:ERROR'STS+900+007+Missing 'Request=''UIT++3'UIZ++1'`
)

var unmarshalTests = []struct {
  in  []byte
  out Values
}{
  // just reversed marshal test cases
  {[]byte(M_OUT1), M_IN1},
  {[]byte(M_OUT2), M_IN2},
  {[]byte(M_OUT3), M_IN3},

  {[]byte(CRAZY_IN1), CRAZY_OUT1},
}

func TestUnmarshal(t *testing.T) {
  // just use the marshal tests
  for i, tt := range marshalTests {
    out, err := Unmarshal(tt.out)
    if err != nil {
      t.Fatalf("%d. received error: %s", i, err)
    }

    if !reflect.DeepEqual(out, tt.in) {
      t.Fatalf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, out, tt.out)
    }
  }
}
