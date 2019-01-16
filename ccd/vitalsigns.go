package ccd

import (
	"time"

	"github.com/mattn/go-pkg-xmlx"
)

var (
	VitalSignsTid = []string{"2.16.840.1.113883.10.20.1.16", "2.16.840.1.113883.10.20.22.2.4.1"}

	VitalSignsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   VitalSignsTid,
		Priority: 0,
		Func:     parseVitalSigns,
	}
)

type VitalSignResult struct {
	Type  string
	Value string
	Unit  string
}

type VitalSignObservation struct {
	Date   time.Time
	Code   Code
	Result VitalSignResult
}

type VitalSign struct {
	Date         time.Time
	Observations []VitalSignObservation
}

func parseVitalSigns(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		organizerNode := Nget(entryNode, "organizer")
		if organizerNode == nil {
			continue
		}

		vitalSign := VitalSign{}

		effectiveTimeNode := Nget(organizerNode, "effectiveTime")
		t := decodeTime(effectiveTimeNode)
		vitalSign.Date = t.Value

		componentNodes := organizerNode.SelectNodes("*", "component")
		for _, componentNode := range componentNodes {
			observation := VitalSignObservation{}

			codeNode := Nget(componentNode, "code")
			if codeNode == nil {
				continue
			}
			observation.Code.decode(codeNode)

			effectiveTimeNode := Nget(componentNode, "effectiveTime")
			t := decodeTime(effectiveTimeNode)
			observation.Date = t.Value

			valueNode := Nget(componentNode, "value")
			// fmt.Println(valueNode)

			observation.Result = VitalSignResult{
				Type:  valueNode.As("*", "type"),
				Value: valueNode.As("*", "value"),
				Unit:  valueNode.As("*", "unit"),
			}

			vitalSign.Observations = append(vitalSign.Observations, observation)
		}

		ccd.VitalSigns = append(ccd.VitalSigns, vitalSign)
	}

	return nil
}
