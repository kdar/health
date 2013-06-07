package ccd

import (
  "strings"

  //"github.com/davecgh/go-spew/spew"
  //"github.com/jteeuwen/go-pkg-xmlx"
  "io/ioutil"
  "os"
  "path/filepath"
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

    file, err := os.Open(path)
    if err != nil {
      t.Fatal(err)
    }

    data, err := ioutil.ReadAll(file)
    if err != nil {
      t.Fatal(err)
    }
    file.Close()

    _, errs := Unmarshal(data)
    shouldfail := strings.HasPrefix(info.Name(), "fail_")
    if shouldfail && errs == nil {
      t.Fatalf("%s: Expected failure, instead received success.", path)
    } else if !shouldfail && errs != nil {
      t.Fatalf("%s: Failed: %v", path, errs)
    }

    return nil
  })
}
