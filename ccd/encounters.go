package ccd

import (
	"time"

	"github.com/jteeuwen/go-pkg-xmlx"
)

var (
	EncountersTid = []string{"2.16.840.1.113883.10.20.22.2.22.1", "2.16.840.1.113883.10.20.22.4.49"}

	EncountersParser = Parser{
		Type:     PARSE_SECTION,
		Values:   EncountersTid,
		Priority: 0,
		Func:     parseEncounters,
	}
	SDLTid = "2.16.840.1.113883.10.20.22.4.32" //Service Delivery Location TemplateID
)

type Encounter struct {
	Name       string
	Performers []Performer
	Code       Code
	Date       time.Time
	Location   Location
	Complaint  string
	Diagnosis  string //todo: maybe support these
}

type Location struct {
	Name    string
	Code    Code
	Address Address
	Telecom Telecom
}
type Performer struct{}

func parseEncounters(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		encounter := Encounter{}
		encounter.Name = entryNode.Value
		if code := Nget(entryNode, "code"); code != nil {
			encounter.Code.decode(code)
		}

		if effectiveTimeNode := Nget(entryNode, "effectiveTime"); effectiveTimeNode != nil {
			t := decodeTime(effectiveTimeNode)
			encounter.Date = t.Value
		}

		//we loop through anything with a templateID inside of our main entry node
		templates := entryNode.SelectNodes("*", "templateId")
		for _, t := range templates {
			if t.As("*", "root") == SDLTid {
				encounter.Location = parseLocation(t.Parent)
			}
		}

		ccd.Encounters = append(ccd.Encounters, encounter)
	}

	return nil
}

func parseLocation(node *xmlx.Node) Location {
	l := Location{}
	if code := Nget(node, "code"); code != nil {
		l.Code.decode(code)
	}

	if aNode := node.SelectNode("*", "addr"); aNode != nil {
		l.Address = decodeAddress(aNode)
	}

	return l
}
