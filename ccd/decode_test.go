package ccd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"testing"

	"golang.org/x/net/html/charset"

	"github.com/davecgh/go-spew/spew"
	"github.com/mattn/go-pkg-xmlx"
	"github.com/kdar/health/ccd"
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
		err = doc.LoadFile(path, charset.NewReaderLabel)
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
	//successful here means that it only failed if the name started with fail_
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

func TestNulls(t *testing.T) {
	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/specific/nulls.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	if c.Allergies != nil {
		t.Fatal("expected allergies to be nil")
	}
	if c.Encounters != nil {
		t.Fatal("expected encounters to be nil")
	}
	if c.Immunizations != nil {
		t.Fatal("expected immunizations to be nil")
	}
	if c.Medications != nil {
		t.Fatal("expected medications to be nil")
	}
	if c.Problems != nil {
		t.Fatal("expected problems to be nil")
	}
	if c.Results != nil {
		t.Fatal("expected results to be nil")
	}
}

func TestNewStuff(t *testing.T) {
	//t.Skip("just for my own needs")

	c := ccd.NewDefaultCCD()
	err := parseAndRecover(t, c, "testdata/public/CCD.sample.xml", nil)
	if err != nil {
		t.Fatal(err)
	}

	_ = spew.Dump

	spew.Dump(c.Allergies[0])
}
