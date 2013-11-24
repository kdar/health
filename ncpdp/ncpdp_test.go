package ncpdp

import (
	//"fmt"
	"github.com/davecgh/go-spew/spew"
	//"github.com/kr/pretty"
	"bytes"
	//"reflect"
	"rxmg/lib/edi/edifact"
	"testing"
	"time"
)

func timeParse(value string) *time.Time {
	v, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", value)
	return &v
}

func durationParse(value string) *time.Duration {
	v, _ := time.ParseDuration(value)
	return &v
}

var (
	UM_OUT1 = &RXHRES{
		RXHREX: &RXHREX{
			RXH: &RXH{
				UIB: edifact.Values{
					"UIB",
					edifact.Values{
						"UNOA",
						"0",
					},
					"",
					"hJIAmKH0FGDSt",
					"",
					"",
					edifact.Values{
						"Sender1",
						"ZZZ",
					},
					edifact.Values{
						"Recepient1",
						"ZZZ",
						"Recepient2",
					},
					edifact.Values{
						"20130113",
						"121625,0",
					},
				},
				UIH: edifact.Values{
					"UIH",
					edifact.Values{
						"SCRIPT",
						"008",
						"001",
						"RXHRES",
					},
					"hJIAmKH0FGDSt",
				},
				UIT: edifact.Values{
					"UIT",
					"hJIAmKH0FGDSt",
					"8",
				},
				UIZ: edifact.Values{
					"UIZ",
					"",
					"",
				},
			},
			PVD: edifact.Values{},
			PTT: edifact.Values{},
			COO: edifact.Values{
				"COO",
				edifact.Values{
					"Per-Se",
					"2U",
				},
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				edifact.Values{
					edifact.Values{
						"07",
						"20120115",
						"102",
					},
					edifact.Values{
						"36",
						"20130113",
						"102",
					},
				},
				"",
				"",
				"",
				"Y",
				"",
			},
			Patient: &Patient{
				Relationship: "1",
				Dob:          *timeParse("1985-06-30 00:00:00 +0000 UTC"),
				Name: &Name{
					Last:   "Smith",
					First:  "John",
					Middle: "",
					Suffix: "",
					Prefix: "",
				},
				Gender:             "M",
				ReferenceNumber:    "",
				ReferenceQualifier: "",
				Address: &Address{
					Line1:             "",
					City:              "",
					State:             "",
					Postal:            "33165",
					LocationQualifier: "",
					Location:          "",
				},
				Phones: []*Phone{},
			},
		},
		Response: &Response{
			ResponseType:      "A",
			CodeListQualifier: "",
			ReferenceNumber:   "",
			Text:              "",
		},
		RequestingPhysician: nil,
		Drugs: []*Drug{{
			ItemDescriptionIdentification: "D",
			ItemDescription:               "XOLEGEL 2% GEL",
			ItemNumber:                    "16110008045",
			CodeListResponsibilityAgency:  "ND",
			QuantityQualifier:             "ZZ",
			Quantity:                      "45",
			CodeListQualifier:             "38",
			DosageId:                      "",
			Dosage1:                       "APPLY TO AFFECTED AREA TWICE A DAY AS NEEDED",
			Dosage2:                       "",
			DaysSupply:                    durationParse("720h0m0s"),
			DateIssued:                    nil,
			LastDemand:                    timeParse("2012-04-22 00:00:00 +0000 UTC"),
			Substitution:                  "0",
			Prescriber: &Provider{
				ProviderCode:       "",
				ReferenceNumber:    "BW7412396",
				ReferenceQualifier: "DH",
				Name: &Name{
					Last:   "Betterton",
					First:  "Jill",
					Middle: "",
					Suffix: "",
					Prefix: "",
				},
				PartyName: "",
			},
			Pharmacy: &Provider{
				ProviderCode:       "",
				ReferenceNumber:    "1031232",
				ReferenceQualifier: "D3",
				Name:               nil,
				PartyName:          "CVS PHARMACY",
			},
		},
		},
	}
)

const (
	UM_IN1 = `UNA~|.?^'UIB|UNOA~0||hJIAmKH0FGDSt|||Sender1~ZZZ|Recepient1~ZZZ~Recepient2|20130113~121625,0'UIH|SCRIPT~008~001~RXHRES|hJIAmKH0FGDSt'RES|A'PTT|1|19850630|Smith~John|M||~~~33165~~|'COO|Per-Se~2U||||||||07~20120115~102^36~20130113~102||||Y|'DRU|D~XOLEGEL?2%?GEL~16110008045~ND~~~~~~~~|ZZ~45~38|~APPLY?TO?AFFECTED?AREA?TWICE?A?DAY?AS?NEEDED~|LD~20120422~102^ZDS~30~804|0|R~0|||||'PVD|P2|1031232~D3|||||CVS?PHARMACY|55596?SW?16TH?ST~ATLANTA~GA~12175|3053872415~TE^~FX|'PVD|PC|BW7412396~DH|||Betterton~Jill~~~|||9900?SW?87th?Ave~Atlanta~GA~12173|~TE'UIT|hJIAmKH0FGDSt|8'UIZ||'`
	//UM_IN2 = `UNA~|.?^'UIB|UNOA~0||SsKyHHXRef6lo|||Sender1~ZZZ|Recepient1~ZZZ~Recepient2|20130113~122506,0'UIH|SCRIPT~008~001~RXHRES|SsKyHHXRef6lo'RES|D|||Cannot find any available prescriptions for this individual.'PTT|1|19820514|doe~jane|F||~~~15672~~|'COO|Per-Se~2U||||||||07~20120115~102^36~20130113~102||||Y|'UIT|SsKyHHXRef6lo|5'UIZ||'`
)

var unmarshalTests = []struct {
	in  []byte
	out interface{}
}{
	{[]byte(UM_IN1), UM_OUT1},
}

var unmarshalValuesTests = []struct {
	in  edifact.Value
	out interface{}
}{
	{
		edifact.Values{edifact.Header{"UNA", ":+./*'"}, edifact.Values{"UIB", edifact.Values{"UNOA", "0"}, "", edifact.Values{"UYdDxBInCbq8", "0"}, "", "", edifact.Values{"Per-Se", "ZZZ"}, edifact.Values{"Sender1", "ZZZ"}, edifact.Values{"20130222", "015815"}}, edifact.Values{"UIH", edifact.Values{"SCRIPT", "008", "001", "ERROR"}, "UYdDxBInCbq8"}, edifact.Values{"STS", "900", "007", edifact.Values{"Practice", "Facility Not Found."}}, edifact.Values{"UIT", "", "3"}, edifact.Values{"UIZ", "", "1"}},
		nil, // TODO: fill this in
	},
}

func TestUnmarshal(t *testing.T) {
	for i, tt := range unmarshalTests {
		out, err := Unmarshal(tt.in)
		if err != nil {
			t.Fatalf("%d. received error: %s", i, err)
		}

		b1 := &bytes.Buffer{}
		b2 := &bytes.Buffer{}

		spew.Fprintf(b1, "%#v", out)
		spew.Fprintf(b2, "%#v", tt.out)

		if !bytes.Equal(b1.Bytes(), b2.Bytes()) {
			t.Fatalf("#%d: mismatch\nhave: %s\nwant: %s", i, b1.String(), b2.String())
		}

		// Cannot deep equal because the data contains pointers which are
		// compared by their address.
		// if !reflect.DeepEqual(out, tt.out) {
		//   t.Fatalf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, out, tt.out)
		// }
	}
}

func TestUnmarshalValues(t *testing.T) {

}
