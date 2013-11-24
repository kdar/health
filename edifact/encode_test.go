package edifact

import (
	"bytes"
	//"fmt"
	"testing"
)

var (
	M_IN1 = Values{
		Header{
			"UNA",
			"~|.?^'",
		},
		Values{
			"UIB",
			Values{"UNOA", "0"},

			"",
			"hJIAmKH0FGDSt",
			"",
			"",

			Values{
				"Senderid1",
				"ZZZ",
				"Senderid2",
			},
			Values{
				"Recepientid1",
				"ZZZ",
				"Recepientid2",
			},
			Values{"20130113", "051612"},
		}, Values{
			"UIH",
			Values{
				"SCRIPT",
				"008",
				"001",
				"RXHREQ",
			},
			"hJIAmKH0FGDSt",
		}, Values{
			"PTT",
			"1",
			"19900807",
			Values{"Smith", "John"},
			"M",
			"",
			Values{"", "", "", "32385", "", ""},
			"",
		}, Values{
			"COO",
			Values{"Per-Se", "2U"},
			"",
			"",
			"", "", "", "", "",
			Values{Values{"07", "20120115", "102"}, Values{"36", "20130113", "102"}},
			"", "", "",
			"Y",
			"",
		}, Values{
			"UIT",
			"hJIAmKH0FGDSt",
			"4",
		}, Values{
			"UIZ",
			"",
			"",
		},
	}

	M_IN2 = Values{
		Header{
			"UNA",
			"~|.?^'",
		},
		Values{
			"TEST",
			"caca~|.?^'|",
			Values{
				Values{"LD", "20121231", "102"},
				Values{"ZDS", "90", "804"},
				Values{"85", "20121001", "102"},
			},
		},
	}

	M_IN3 = Values{
		Header{"UNA", ":+./*'"},
		Values{
			"TES",
			"hey",
			"there",
			"guy",
		},
	}
)

const (
	M_OUT1 = "UNA~|.?^'UIB|UNOA~0||hJIAmKH0FGDSt|||Senderid1~ZZZ~Senderid2|Recepientid1~ZZZ~Recepientid2|20130113~051612'UIH|SCRIPT~008~001~RXHREQ|hJIAmKH0FGDSt'PTT|1|19900807|Smith~John|M||~~~32385~~|'COO|Per-Se~2U||||||||07~20120115~102^36~20130113~102||||Y|'UIT|hJIAmKH0FGDSt|4'UIZ||'"
	M_OUT2 = "UNA~|.?^'TEST|caca?~?|.???^?'?||LD~20121231~102^ZDS~90~804^85~20121001~102'"
	M_OUT3 = `UNA:+./*'TES+hey+there+guy'`
)

var marshalTests = []struct {
	in  Values
	out []byte
}{
	{M_IN1, []byte(M_OUT1)},
	{M_IN2, []byte(M_OUT2)},
	{M_IN3, []byte(M_OUT3)},
}

func TestMarshal(t *testing.T) {
	for i, tt := range marshalTests {
		out, err := Marshal(tt.in)
		if err != nil {
			t.Fatalf("%d. received error: %s", err)
		}

		if !bytes.Equal(out, tt.out) {
			t.Fatalf("%d. unexpected output: %s. want %s", i, string(out), string(tt.out))
		}
	}
}
