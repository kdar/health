package ccd_test

import (
	"github.com/shurcooL/go-goon"
	"reflect"
	"testing"
	"time"

	"github.com/kdar/health/ccd"
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
			Status:    "completed",
			StartDate: time.Time{},
			StopDate:  time.Date(2012, 8, 6, 0, 0, 0, 0, time.UTC),
			Period:    time.Duration(43200000000000),
			Code: ccd.Code{
				CodeSystemName: "",
				Type:           "",
				CodeSystem:     "2.16.840.1.113883.6.88",
				Code:           "573621",
				DisplayName:    "Albuterol 0.09 MG/ACTUAT inhalant solution",
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

	if !reflect.DeepEqual(meds, c.Medications) {
		goon.Dump(c.Medications)
		t.Fatalf("Expected:\n%#v, got:\n%#v", meds, c.Medications)
	}
}
