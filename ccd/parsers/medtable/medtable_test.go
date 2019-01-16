package medtable

import (
	//"github.com/mattn/go-pkg-xmlx"
	"github.com/kdar/health/ccd"
	"testing"
)

func TestMedDate(t *testing.T) {
	cc := ccd.NewDefaultCCD()
	// first three meds missing date in xml data
	err := cc.ParseFile("testdata/ccd.xml")
	if err != nil {
		t.Fatal(err)
	}

	for _, med := range cc.Medications[:3] {
		if !med.StartDate.IsZero() {
			t.Fatal("Expected medication StartDate to be zero.")
		}
	}

	cc.AddParsers(Parser())

	// reparse with our new parser added
	err = cc.ParseFile("testdata/ccd.xml")
	if err != nil {
		t.Fatal(err)
	}

	for _, med := range cc.Medications[:3] {
		if med.StartDate.IsZero() {
			t.Fatal("Expected medication StartDate to be non-zero.")
		}
	}
}
