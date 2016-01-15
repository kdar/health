package ccd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/kdar/health/ccd"
	"github.com/kdar/health/ccd/parsers/medtable"
)

var successfulCCDs int64   //Keeps track of how many CCDs we successfully parsed
var unsuccessfulCCDs int64 //And the ones that do not.

func parseAndRecover(t *testing.T, c *ccd.CCD, path string, doc *xmlx.Document) (err error) {
	defer func() {
		if e := recover(); e != nil {
			lines := bytes.Split(debug.Stack(), []byte{'\n'})
			for i, _ := range lines {
				if len(lines[i]) >= 1 && lines[i][0] == '\t' {
					lines[i] = lines[i][1:]
				}
			}
			t.Errorf("Error processing: %s\n\n%s", path, bytes.Join(lines, []byte{'\n'}))
		}
	}()

	if doc != nil {
		err = c.ParseDoc(doc)
	} else {
		err = c.ParseFile(path)
	}
	return
}

func walkAllCCDs(f filepath.WalkFunc) {
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

		return f(path, info, err)
	})
}

func TestParseAllCCDs(t *testing.T) {
	walkAllCCDs(func(path string, info os.FileInfo, err error) error {
		shouldfail := strings.HasPrefix(info.Name(), "fail_")

		doc := xmlx.New()
		err = doc.LoadFile(path, nil)
		if shouldfail && err != nil {
			return nil
		}

		c := ccd.NewDefaultCCD()
		err = parseAndRecover(t, c, path, doc)
		if shouldfail && err == nil {
			t.Errorf("%s: Expected failure, instead received success.", path)
			atomic.AddInt64(&unsuccessfulCCDs, 1)
		} else if !shouldfail && err != nil {
			t.Errorf("%s: Failed: %v", path, err)
			atomic.AddInt64(&unsuccessfulCCDs, 1)
		} else {
			atomic.AddInt64(&successfulCCDs, 1)
		}

		return nil
	})
	//Note that by successful, failing counts as success if the ccd is named _fail
	t.Logf("parsed %d CCDS. %d successful, %d unsuccessful\n", (successfulCCDs + unsuccessfulCCDs), successfulCCDs, unsuccessfulCCDs)
}

func TestInvalidCCD(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/invalid_ccd.xml", nil)
	if err == nil {
		t.Fatal("Expected parsing of CCD to fail and throw and error.")
	}

	err = parseAndRecover(t, c, "testdata/specific/valid_ccd.xml", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewStuff(t *testing.T) {
	t.Skip("just for my own needs")

	c := ccd.NewDefaultCCD()
	c.AddParsers(medtable.Parser())
	err := parseAndRecover(t, c, "testdata/private/2013-08-26T04_03_24 - 0b7fddbdc631aecc6c96090043f690204f7d0d9d.xml", nil)
	//err := parseAndRecover(t, c, "testdata/public/ToC_CCDA_CCD_CompGuideSample_FullXML_v01a.xml", nil)
	//err := parseAndRecover(t, c, "testdata/public/SampleCCDDocument.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	_ = spew.Dump

	//spew.Dump(c.Patient)
}
