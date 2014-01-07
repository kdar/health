// Provides a parser to help get more medication info from
// <text> tag of the medications section.

package medtable

import (
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/kdar/health/ccd"
	"time"
)

const (
	MedTableTimeFormat = "Jan 02, 2006"
)

func getDataPositions(node *xmlx.Node) (int, int) {
	thead := ccd.Nget(node, "text", "table", "thead", "tr")
	if thead == nil {
		return -1, -1
	}

	headers := thead.SelectNodes("*", "th")
	if headers == nil {
		return -1, -1
	}

	nameIndex := -1
	dateIndex := -1
	for i, n := range headers {
		switch n.S("*", "th") {
		case "Start Date":
			dateIndex = i
		case "Medication":
			nameIndex = i
		}
	}

	return nameIndex, dateIndex
}

func getData(node *xmlx.Node, nameIndex, dateIndex int) (entries []*xmlx.Node) {
	tbody := ccd.Nget(node, "text", "table", "tbody")
	if tbody == nil {
		return nil
	}

	entries = tbody.SelectNodes("*", "tr")
	return
}

// FIXME: the HTML doesn't always have the same amountof nodes as in the <entry>'s
func parseMedTableMedications(node *xmlx.Node, cc *ccd.CCD) []error {
	var dateIndex int

	nameIndex, dateIndex := getDataPositions(node)
	if nameIndex == -1 && dateIndex == -1 {
		return nil
	}

	entries := getData(node, nameIndex, dateIndex)

	fmt.Println(len(entries))
	fmt.Println(len(cc.Medications))

	// depends on meds not being modified or out of order
	for i, med := range cc.Medications {
		if med.StartDate.IsZero() {
			if dateIndex != -1 && i < len(entries) {
				tds := entries[i].SelectNodes("*", "td")
				if len(tds) > dateIndex {
					cc.Medications[i].StartDate, _ = time.Parse(MedTableTimeFormat, tds[dateIndex].S("*", "td"))
				}
			}
		}
	}

	return nil
}

func Parser() ccd.Parser {
	return ccd.Parser{
		Type:         ccd.PARSE_SECTION,
		Organization: "*",
		Values:       ccd.MedicationsTid,
		Priority:     10,
		Func:         parseMedTableMedications,
	}
}
