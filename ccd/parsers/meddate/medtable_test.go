package medtable

import (
  //"github.com/jteeuwen/go-pkg-xmlx"
  "github.com/kdar/health/ccd"
  "testing"
)

func TestMedDate(t *testing.T) {
  // first three meds missing date in xml data
  cc, err := ccd.UnmarshalFile("testdata/ccd.xml")
  if err != nil {
    t.Fatal(err)
  }

  for _, med := range cc.Medications[:3] {
    if !med.StartDate.IsZero() {
      t.Fatal("Expected medication StartDate to be zero.")
    }
  }

  ccd.Register(Parser())

  // reparse with our new parser added
  cc2, err := ccd.UnmarshalFile("testdata/ccd.xml")
  if err != nil {
    t.Fatal(err)
  }

  for _, med := range cc2.Medications[:3] {
    if med.StartDate.IsZero() {
      t.Fatal("Expected medication StartDate to be non-zero.")
    }
  }
}
