package ccd_test

import (
	"testing"

	"github.com/kdar/health/ccd"
)

func TestAllScriptsVitalSigns(t *testing.T) {
	c := ccd.NewCCD()
	c.AddParsers(ccd.VitalSignsParser)
	err := c.ParseFile("testdata/given/allscriptsccd.xml")
	if err != nil {
		t.Fatal(err)
	}

	if len(c.VitalSigns) != 9 {
		t.Errorf("Expected 9 vital signs, got: %d", len(c.VitalSigns))
	}
}
