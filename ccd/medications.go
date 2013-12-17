package ccd

import (
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"time"
)

var (
	MedicationsTid = []string{"2.16.840.1.113883.10.20.1.8", "2.16.840.1.113883.10.20.22.2.1.1"}

	MedicationsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   MedicationsTid,
		Priority: 0,
		Func:     parseMedications,
	}
)

type MedicationDose struct {
	LowValue  string
	LowUnit   string
	HighValue string
	HighUnit  string
}

func (m MedicationDose) ValueUnit() (string, string) {
	unit := m.LowUnit
	if len(unit) == 0 {
		unit = m.HighUnit
	}

	if len(m.HighValue) == 0 {
		return m.LowValue, unit
	}

	return fmt.Sprintf("%s-%s", m.LowValue, m.HighValue), unit
}

func (m MedicationDose) String() string {
	unit := m.LowUnit
	if len(unit) == 0 {
		unit = m.HighUnit
	}

	if len(unit) > 0 {
		unit = " " + unit
	}

	if len(m.HighValue) == 0 {
		return fmt.Sprintf("%s%s", m.LowValue, m.LowUnit)
	}

	return fmt.Sprintf("%s-%s%s", m.LowValue, m.HighValue, unit)
}

type Medication struct {
	Name           string
	Administration string
	//Instructions   string // this is calulated and not specifically in the CCD
	Dose      MedicationDose
	Status    string
	StartDate time.Time
	StopDate  time.Time
	Period    time.Duration
	Code      Code
}

// http://wiki.ihe.net/index.php?title=1.3.6.1.4.1.19376.1.5.3.1.4.7
func parseMedications(node *xmlx.Node, ccd *CCD) []error {
	var errs []error

	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		// if templateId(entryNode) != "2.16.840.1.113883.10.20.1.24" {
		//   continue
		// }

		medication := Medication{}

		saNode := Nget(entryNode, "substanceAdministration")

		doseNode := Nget(saNode, "doseQuantity")
		if doseNode != nil {
			medication.Dose.LowValue = doseNode.As("*", "value")
			medication.Dose.LowUnit = doseNode.As("*", "unit")
			doseLowNode := Nget(doseNode, "low")
			doseHighNode := Nget(doseNode, "high")
			if doseLowNode != nil {
				medication.Dose.LowValue = doseLowNode.As("*", "value")
				medication.Dose.LowUnit = doseLowNode.As("*", "unit")
			}
			if doseHighNode != nil {
				medication.Dose.HighValue = doseHighNode.As("*", "value")
				medication.Dose.HighUnit = doseHighNode.As("*", "unit")
			}
		}

		routeCodeNode := Nget(saNode, "routeCode")
		if routeCodeNode != nil {
			if routeCodeNode.As("*", "codeSystemName") == "RouteOfAdministration" {
				medication.Administration = routeCodeNode.As("*", "displayName")
			}
		}

		mpNode := Nget(saNode, "consumable", "manufacturedProduct")
		if mpNode == nil {
			continue
		}

		medication.Status = Nsget(saNode, "statusCode").As("*", "code")

		etimeNodes := saNode.SelectNodes("*", "effectiveTime")
		for _, etimeNode := range etimeNodes {
			t := decodeTime(etimeNode)
			if t.Type == TIME_INTERVAL {
				medication.StartDate = t.Low
				medication.StopDate = t.High
			} else if t.Type == TIME_PERIODIC {
				medication.Period = t.Period
			}
		}

		manNode := Nget(mpNode, "manufacturedMaterial")

		codeNode := Nget(manNode, "code")
		if codeNode != nil {
			medication.Code.decode(codeNode)
			if medication.Code.Code == "" {
				// Sometimes the attributes for "code" are completely missing.
				// try to see if there is a translation node and get it from there
				transNode := codeNode.SelectNode("*", "translation")
				medication.Code.decode(transNode)
			}
		}

		medication.Name = medication.Code.DisplayName

		if nameNode := manNode.SelectNode("*", "name"); nameNode != nil {
			medication.Name = nameNode.GetValue()
		} else if medication.Name == "" {
			if originalNode := codeNode.SelectNode("*", "originalText"); originalNode != nil {
				medication.Name = originalNode.GetValue()
			}
		}

		ccd.Medications = append(ccd.Medications, medication)
	}

	return errs
}
