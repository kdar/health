package ccd

import (
	"github.com/jteeuwen/go-pkg-xmlx"
	"time"
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

type VitalSign struct {
	Name   string
	Result VitalSignResult
	Date   time.Time
}

func parseVitalSigns(node *xmlx.Node, ccd *CCD) []error {
	orgNode := Nget(node, "entry", "organizer")
	if orgNode == nil {
		return nil
	}

	componentNodes := orgNode.SelectNodes("*", "component")
	for _, componentNode := range componentNodes {
		vitalsign := VitalSign{}

		codeNode := Nget(componentNode, "code")
		vitalsign.Name = codeNode.As("*", "displayName")

		effectiveTimeNode := Nget(componentNode, "effectiveTime")
		t := ParseTimeNode(effectiveTimeNode)
		vitalsign.Date = t.Value

		valueNode := Nget(componentNode, "value")
		vitalsign.Result = VitalSignResult{
			Type:  valueNode.As("*", "type"),
			Value: valueNode.As("*", "value"),
			Unit:  valueNode.As("*", "unit"),
		}

		ccd.VitalSigns = append(ccd.VitalSigns, vitalsign)
	}

	return nil
}
