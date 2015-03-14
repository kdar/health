package ccd_test

import (
	"reflect"
	"testing"

	"github.com/kdar/health/ccd"
)

func TestParse_Address(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/address.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	addr := ccd.Address{
		Line1:   "Line1",
		Line2:   "Line2",
		City:    "City",
		County:  "County",
		State:   "ST",
		Zip:     "12345",
		Country: "Country",
		Use:     "HP",
	}

	if !reflect.DeepEqual(addr, c.Patient.Addresses[0]) {
		t.Fatalf("Expected:\n%#v, got:\n%#v", addr, c.Patient.Addresses[0])
	}

	if !c.Patient.Name.IsZero() {
		t.Fatalf("Patient.Name was suppose to be empty, but it's not")
	}
}

func TestParse_Name(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/name.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	name := ccd.Name{
		First:    "First",
		Middle:   "Middle",
		Last:     "Last",
		Suffix:   "Suffix",
		Prefix:   "Prefix",
		Type:     "PN",
		NickName: "NickName",
	}

	if !reflect.DeepEqual(name, c.Patient.Name) {
		t.Fatalf("Expected:\n%#v, got:\n%#v", name, c.Patient.Name)
	}
}

func TestParse_MissingName(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/name_missing.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	name := ccd.Name{
		First:    "",
		Middle:   "",
		Last:     "",
		Suffix:   "",
		Prefix:   "",
		Type:     "",
		NickName: "",
	}

	if !reflect.DeepEqual(name, c.Patient.Name) {
		t.Fatalf("Expected:\n%#v, got:\n%#v", name, c.Patient.Name)
	}
}
