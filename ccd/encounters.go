package ccd

import (
	"time"

	"github.com/jteeuwen/go-pkg-xmlx"
)

var (
	TidEncounters    = []string{"2.16.840.1.113883.10.20.22.2.22.1", "2.16.840.1.113883.10.20.22.4.49"}
	EncountersParser = Parser{
		Type:     PARSE_SECTION,
		Values:   TidEncounters,
		Priority: 0,
		Func:     parseEncounters,
	}
	//The following TemplateID constants are used to find specific sections inside of an encounter.
	TidSDL         = "2.16.840.1.113883.10.20.22.4.32" //TemplateID of the Service Delivery Location
	TidEncounterDx = "2.16.840.1.113883.10.20.22.4.80" //TemplateID of the Encounter Diagnosis
)

type Encounter struct {
	Performers []Performer
	Code       Code
	Date       time.Time
	Location   Location
	Diagnosis  Diagnosis
	Complaint  string //todo: this, look for 'RSON'
}

type Location struct {
	Name    string
	Code    Code
	Address Address
	Telecom Telecom
}
type Diagnosis struct {
	Code    Code
	Status  string
	Problem Problem
	//TODO: add support for EffectiveTime
}
type Performer struct{}

func parseEncounters(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		encounter := Encounter{}
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

			switch t.As("*", "root") {
			case TidSDL:
				encounter.Location = decodeLocation(t.Parent)
			case TidEncounterDx:
				encounter.Diagnosis = decodeDiagnosis(t.Parent)
			}
		}

		ccd.Encounters = append(ccd.Encounters, encounter)
	}

	return nil
}

func decodeDiagnosis(node *xmlx.Node) Diagnosis {
	d := Diagnosis{}
	if code := Nget(node, "code"); code != nil {
		d.Code.decode(code)
	}
	d.Status = Nget(node, "statusCode").As("*", "code")

	if problem := Nget(node, "observation"); problem != nil {
		d.Problem = decodeProblem(problem)
	}

	return d
}

func decodeLocation(node *xmlx.Node) Location {
	l := Location{}
	if code := Nget(node, "code"); code != nil {
		l.Code.decode(code)
	}

	if aNode := node.SelectNode("*", "addr"); aNode != nil {
		l.Address = decodeAddress(aNode)
	}
	if tNode := node.SelectNode("*", "telecom"); tNode != nil {
		l.Telecom = decodeTelecom(tNode)
	}
	if peNode := node.SelectNode("*", "playingEntity"); peNode != nil {
		l.Name = peNode.S("*", "name")
	}

	return l
}
