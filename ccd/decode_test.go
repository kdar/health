package ccd

import (
  "strings"

  //"github.com/davecgh/go-spew/spew"
  //"github.com/jteeuwen/go-pkg-xmlx"
  "os"
  "path/filepath"
  "reflect"
  "testing"
)

func TestUnmarshal(t *testing.T) {
  filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
    if info.IsDir() {
      return nil
    }

    // if info.Name() != "SampleCCDDocument.xml" {
    //   return nil
    // }

    _, err = UnmarshalFile(path)
    shouldfail := strings.HasPrefix(info.Name(), "fail_")
    if shouldfail && err == nil {
      t.Fatalf("%s: Expected failure, instead received success.", path)
    } else if !shouldfail && err != nil {
      t.Fatalf("%s: Failed: %v", path, err)
    }

    return nil
  })
}

func TestUnmarshal_Address(t *testing.T) {
  ccd, err := UnmarshalFile("testdata/specific/address.xml")
  if err != nil {
    t.Fatal(err)
  }

  addr := Address{
    Line1:   "Line1",
    Line2:   "Line2",
    City:    "City",
    County:  "County",
    State:   "ST",
    Zip:     "12345",
    Country: "Country",
    Type:    "HP",
  }

  if !reflect.DeepEqual(addr, ccd.Patient.Address) {
    t.Fatalf("Expected:\n%#v, got:\n%#v", addr, ccd.Patient.Address)
  }

  if !ccd.Patient.Name.IsZero() {
    t.Fatalf("ccd.Patient.Name was suppose to be empty, but it's not")
  }
}

func TestUnmarshal_Name(t *testing.T) {
  ccd, err := UnmarshalFile("testdata/specific/name.xml")
  if err != nil {
    t.Fatal(err)
  }

  name := Name{
    First:    "First",
    Middle:   "Middle",
    Last:     "Last",
    Suffix:   "Suffix",
    Prefix:   "Prefix",
    Type:     "PN",
    NickName: "NickName",
  }

  if !reflect.DeepEqual(name, ccd.Patient.Name) {
    t.Fatalf("Expected:\n%#v, got:\n%#v", name, ccd.Patient.Name)
  }
}
