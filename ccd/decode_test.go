package ccd

import (
  "bytes"
  "github.com/davecgh/go-spew/spew"
  "runtime/debug"
  "strings"
  //"github.com/jteeuwen/go-pkg-xmlx"
  "os"
  "path/filepath"
  "reflect"
  "testing"
)

func unmarshalAndRecover(t *testing.T, c *CCD, path string) (err error) {
  defer func() {
    if e := recover(); e != nil {
      lines := bytes.Split(debug.Stack(), []byte{'\n'})
      for i, _ := range lines {
        if lines[i][0] == '\t' {
          lines[i] = lines[i][1:]
        }
      }
      t.Fatalf("Error processing: %s\n\n%s", path, bytes.Join(lines, []byte{'\n'}))
    }
  }()

  err = c.ParseFile(path)
  return
}

func TestParseAllCCDs(t *testing.T) {
  filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
    if info.IsDir() {
      if info.Name()[0] == '.' {
        return filepath.SkipDir
      }

      return nil
    }

    if !strings.HasSuffix(path, "xml") && !strings.HasSuffix(path, "ccd") {
      return nil
    }

    c := NewDefaultCCD()
    err = unmarshalAndRecover(t, c, path)
    shouldfail := strings.HasPrefix(info.Name(), "fail_")
    if shouldfail && err == nil {
      t.Fatalf("%s: Expected failure, instead received success.", path)
    } else if !shouldfail && err != nil {
      t.Fatalf("%s: Failed: %v", path, err)
    }

    return nil
  })
}

func TestNewStuff(t *testing.T) {
  c := NewDefaultCCD()
  err := unmarshalAndRecover(t, c, "testdata/private/2013-08-26T04_03_24 - 0b7fddbdc631aecc6c96090043f690204f7d0d9d.xml")
  if err != nil {
    t.Fatal(err)
  }

  _ = spew.Dump

  spew.Dump(c.Medications)
}

func TestParse_Address(t *testing.T) {
  c := NewDefaultCCD()
  err := unmarshalAndRecover(t, c, "testdata/specific/address.xml")
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

  if !reflect.DeepEqual(addr, c.Patient.Address) {
    t.Fatalf("Expected:\n%#v, got:\n%#v", addr, c.Patient.Address)
  }

  if !c.Patient.Name.IsZero() {
    t.Fatalf("Patient.Name was suppose to be empty, but it's not")
  }
}

func TestParse_Name(t *testing.T) {
  c := NewDefaultCCD()
  err := unmarshalAndRecover(t, c, "testdata/specific/name.xml")
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

  if !reflect.DeepEqual(name, c.Patient.Name) {
    t.Fatalf("Expected:\n%#v, got:\n%#v", name, c.Patient.Name)
  }
}
