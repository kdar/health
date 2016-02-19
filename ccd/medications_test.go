package ccd_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kdar/health/ccd"
	"github.com/kylelemons/godebug/pretty"
)

func TestParse_Medications(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/medications.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	meds := []ccd.Medication{
		ccd.Medication{
			Name:           "Albuterol 0.09 MG/ACTUAT inhalant solution",
			Administration: "",
			Dose: ccd.MedicationDose{
				LowValue:  "0.09",
				LowUnit:   "mg/actuat",
				HighValue: "",
				HighUnit:  "",
			},
			Status:     "Active",
			StatusCode: "completed",
			StartDate:  time.Time{},
			StopDate:   time.Date(2012, 8, 6, 0, 0, 0, 0, time.UTC),
			Period:     time.Duration(43200000000000),
			Code: ccd.Code{
				CodeSystemName: "",
				Type:           "",
				CodeSystem:     "2.16.840.1.113883.6.88",
				Code:           "573621",
				DisplayName:    "Albuterol 0.09 MG/ACTUAT inhalant solution",
				Translations: []ccd.Code{ccd.Code{
					CodeSystemName: "RxNorm",
					Type:           "",
					CodeSystem:     "2.16.840.1.113883.6.88",
					Code:           "573621",
					DisplayName:    "Proventil 0.09 MG/ACTUAT inhalant solution",
					OriginalText:   "",
				}},
			},
			Reason: &ccd.MedicationReason{
				Value: ccd.Code{
					CodeSystemName: "",
					Type:           "CD",
					CodeSystem:     "2.16.840.1.113883.6.96",
					Code:           "233604007",
					DisplayName:    "Pneumonia",
				},
				Date: time.Time{},
			},
		},
	}

	if len(meds) != len(c.Medications) {
		t.Fatalf("Expected %d medications. Got %d", len(meds), len(c.Medications))
	}

	for i, _ := range meds {
		if !reflect.DeepEqual(meds[i], c.Medications[i]) {
			fmt.Println(pretty.Compare(meds[i], c.Medications[i]))
			t.Fatal()
		}

	}

}
