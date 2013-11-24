// Provides a parser to help get more medication info from
// <text> tag of the medications section.

package medtable

import (
  "github.com/jteeuwen/go-pkg-xmlx"
  "github.com/kdar/health/ccd"
  "time"
)

const (
  MedTableTimeFormat = "Jan 02, 2006"
)

func medDateGetData(node *xmlx.Node) (entries []*xmlx.Node, dateIndex int) {
  thead := ccd.Nget(node, "text", "table", "thead", "tr")
  if thead == nil {
    return nil, -1
  }

  headers := thead.SelectNodes("*", "th")
  if headers == nil {
    return nil, -1
  }

  dateIndex = -1
  for i, n := range headers {
    if n.S("*", "th") == "Start Date" {
      dateIndex = i
      break
    }
  }

  if dateIndex == -1 {
    return nil, -1
  }

  tbody := ccd.Nget(node, "text", "table", "tbody")
  if thead == nil {
    return nil, -1
  }

  entries = tbody.SelectNodes("*", "tr")
  return
}

func parseMedTableMedications(node *xmlx.Node, cc *ccd.CCD) []error {
  var dateIndex int
  var entries []*xmlx.Node

  // depends on meds not being modified or out of order
  for i, med := range cc.Medications {
    if med.StartDate.IsZero() {
      // initalize stuff we need
      if entries == nil {
        entries, dateIndex = medDateGetData(node)
        if dateIndex == -1 {
          return nil
        }
      }

      tds := entries[i].SelectNodes("*", "td")
      if len(tds) > dateIndex {
        cc.Medications[i].StartDate, _ = time.Parse(MedTableTimeFormat, tds[dateIndex].S("*", "td"))
      }
    }
  }

  return nil
}

func Parser() ccd.Parser {
  return ccd.Parser{
    Type:         ccd.PARSE_SECTION,
    Organization: "Good Health Clinic",
    Value:        "2.16.840.1.113883.10.20.1.8",
    Priority:     10,
    Func:         parseMedTableMedications,
  }
}
