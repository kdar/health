package ccd

import (
	"fmt"
	"time"

	"github.com/mattn/go-pkg-xmlx"
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

type MedicationReason struct {
	Value Code
	Date  time.Time
}

type Medication struct {
	Name           string
	Administration string
	Dose           MedicationDose
	// Active, On Hold, Prior History, No Longer Active
	// http://motorcycleguy.blogspot.com/2011/03/medication-status-in-ccd.html
	Status string
	// Document HL7 ActStatus: aborted / active / cancelled / completed / held / new / suspended
	StatusCode string
	StartDate  time.Time
	StopDate   time.Time
	Period     time.Duration
	Code       Code
	Reason     *MedicationReason

	//Instructions   string // this is calulated and not specifically in the CCD
}

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

		medication.StatusCode = Nsget(saNode, "statusCode").As("*", "code")

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
		}

		medication.Name = medication.Code.DisplayName

		if nameNode := manNode.SelectNode("*", "name"); nameNode != nil {
			medication.Name = nameNode.GetValue()
		} else if medication.Name == "" {
			if originalNode := codeNode.SelectNode("*", "originalText"); originalNode != nil {
				medication.Name = originalNode.GetValue()
			}
		}

		// If we still don't have a name, just continue because there's no
		// point in processing this medication without a name.
		if len(medication.Name) == 0 {
			continue
		}

		entryRelationshipNodes := saNode.SelectNodes("*", "entryRelationship")
		if entryRelationshipNodes != nil {
			for _, entryRelationshipNode := range entryRelationshipNodes {
				switch entryRelationshipNode.As("*", "typeCode") {
				case "RSON":
					medication.Reason = parseMedicationReason(entryRelationshipNode)
				case "REFR":
					obvNode := Nget(entryRelationshipNode, "observation")
					if obvNode != nil && Nget(obvNode, "statusCode") != nil {
						valueNode := Nget(obvNode, "value")
						if valueNode != nil {
							medication.Status = valueNode.As("*", "displayName")
						}
					}
				}
			}
		}

		ccd.Medications = append(ccd.Medications, medication)
	}

	return errs
}

func parseMedicationReason(node *xmlx.Node) *MedicationReason {
	observationNode := Nget(node, "observation")
	if observationNode == nil {
		return nil
	}

	reason := &MedicationReason{}

	effectiveTimeNode := Nget(observationNode, "effectiveTime")
	t := decodeTime(effectiveTimeNode)
	reason.Date = t.Value

	valueNode := Nget(observationNode, "value")
	reason.Value.decode(valueNode)

	return reason
}
