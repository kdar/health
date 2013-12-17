package ccd

import (
	"github.com/jteeuwen/go-pkg-xmlx"
	"time"
)

var (
	ImmunizationsTid = []string{"2.16.840.1.113883.10.20.1.6", "2.16.840.1.113883.10.20.22.2.2.1"}

	ImmunizationsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   ImmunizationsTid,
		Priority: 0,
		Func:     parseImmunizations,
	}
)

type Immunization struct {
	Name           string
	Administration string
	Date           time.Time
	Status         string
}

func parseImmunizations(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		immunization := Immunization{}

		saNode := Nget(entryNode, "substanceAdministration")
		immunization.Status = Nget(saNode, "statusCode").As("*", "code")

		t := decodeTime(Nget(saNode, "effectiveTime"))
		immunization.Date = t.Value

		routeCodeNode := Nget(saNode, "routeCode")
		if routeCodeNode != nil {
			//if routeCodeNode.As("*", "codeSystemName") == "RouteOfAdministration" {
			immunization.Administration = routeCodeNode.As("*", "displayName")
			//}
		}

		codeNode := Nget(saNode, "manufacturedProduct", "manufacturedMaterial", "code")
		if codeNode != nil {
			immunization.Name = codeNode.As("*", "displayName")
		}

		ccd.Immunizations = append(ccd.Immunizations, immunization)
	}

	return nil
}
