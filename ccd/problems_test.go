package ccd_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kdar/health/ccd"
	"github.com/kylelemons/godebug/pretty"
)

//parseTime returns a time.Time based on a premade layout
func parseTime(in string) time.Time {
	time, _ := time.Parse("2006-01-02", in)
	return time

}

//TestParse_Problems parses a specific ccd with many different types of problems and compares it to a pre-made output.
func TestParse_Problems(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/problems.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	problems := []ccd.Problem{
		ccd.Problem{
			Name:        "Pneumonia",
			Time:        ccd.Time{Low: parseTime("2012-08-06"), Value: parseTime("2012-08-06")},
			Status:      "Active",
			ProblemType: "Complaint",
			Code: ccd.Code{
				CodeSystemName: "",
				Type:           "CD",
				CodeSystem:     "2.16.840.1.113883.6.96",
				Code:           "233604007",
				DisplayName:    "Pneumonia",
			},
		},
		ccd.Problem{
			Name:        "CEREBRAL ARTERY OCCLUSION, UNSPECIFIED, WITH CEREBRAL INFARCTION",
			Time:        ccd.Time{Low: parseTime("2009-07-09"), Value: parseTime("2009-07-09")},
			Status:      "Active",
			ProblemType: "Problem",
			Code: ccd.Code{
				CodeSystemName: "ICD-9-CM",
				Type:           "",
				CodeSystem:     "2.16.840.1.113883.6.104",
				Code:           "434.91",
				DisplayName:    "CEREBRAL ARTERY OCCLUSION, UNSPECIFIED, WITH CEREBRAL INFARCTION",
			},
		},
		ccd.Problem{
			Name:        "Body mass index 30+ - obesity",
			Time:        ccd.Time{Low: parseTime("2011-01-18"), Value: parseTime("2011-01-18")},
			Status:      "Active",
			ProblemType: "",
			Code: ccd.Code{
				CodeSystemName: "SNOMED CT",
				Type:           "",
				CodeSystem:     "2.16.840.1.113883.6.96",
				Code:           "162864005",
				DisplayName:    "",
			},
		},
		ccd.Problem{
			Name:        "GERD (gastroesophageal reflux disease)",
			Time:        ccd.Time{},
			Status:      "Active",
			ProblemType: "Problem",
			Code: ccd.Code{
				CodeSystemName: "ICD-9 CM",
				Type:           "",
				CodeSystem:     "2.16.840.1.113883.6.103",
				Code:           "530.81",
				DisplayName:    "GERD (gastroesophageal reflux disease)",
			},
		},
		ccd.Problem{
			Name:        "JOINT PAIN PELVIS",
			Time:        ccd.Time{},
			Status:      "Active",
			ProblemType: "ASSERTION",
			Code: ccd.Code{
				Type:           "CD",
				CodeSystem:     "2.16.840.1.113883.6.103",
				CodeSystemName: "",
				Code:           "719.45",
				DisplayName:    "JOINT PAIN PELVIS",
			},
		}}

	if len(problems) != len(c.Problems) {
		t.Fatalf("Wrong number of problems in specific problems test. Expected: %d got: %d", len(problems), len(c.Problems))
	}
	for i, _ := range problems {
		if !reflect.DeepEqual(problems[i], c.Problems[i]) {
			fmt.Println(pretty.Compare(problems[i], c.Problems[i]))
			t.Fatal()
			//		t.Fatalf("Expected:\n%s, got:\n%s", goon.Sdump(problems), goon.Sdump(c.Problems))
		}

	}
}
