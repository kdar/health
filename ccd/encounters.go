package ccd

import "github.com/mattn/go-pkg-xmlx"

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
	TidIndication  = "2.16.840.1.113883.10.20.22.4.19"
)

//Encounter describes any interaction with the patient and healthcare provider.
//See: http://www.cdatools.org/infocenter/index.jsp?topic=%2Forg.openhealthtools.mdht.uml.cda.consol.doc%2Fclasses%2FEncounterDiagnosis.html
type Encounter struct {
	Performers []Performer
	Code       Code
	Time       Time

	//CCD calls this a 'participant', but we only use it for Service Delivery Locations
	Locations []Location

	/*List of problems found as a result of this visit*/
	Diagnosis []Diagnosis

	/*List of problems that caused the visit to happen*/
	Indications []Problem
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
}

type Performer struct {
	Address Address
	Telecom Telecom
	Name    Name
	Code    Code
}

func parseEncounters(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		encounter := Encounter{}
		if code := Nget(entryNode, "code"); code != nil {
			encounter.Code.decode(code)
		}

		// Continue if this has no display name.
		if len(encounter.Code.DisplayName) == 0 {
			continue
		}

		if effectiveTimeNode := Nget(entryNode, "effectiveTime"); effectiveTimeNode != nil {
			encounter.Time = decodeTime(effectiveTimeNode)
		}

		//we loop through anything with a templateID inside of our main entry node
		templates := entryNode.SelectNodes("*", "templateId")
		for _, t := range templates {

			switch t.As("*", "root") {
			case TidSDL:
				encounter.Locations = append(encounter.Locations, decodeLocation(t.Parent))
			case TidEncounterDx:
				encounter.Diagnosis = append(encounter.Diagnosis, decodeDiagnosis(t.Parent))
			case TidIndication:
				//TODO: should we check for @typeCode="RSON"?
				problem := decodeProblem(t.Parent)
				if problem != nil {
					encounter.Indications = append(encounter.Indications, *problem)
				}
			}
		}

		performers := entryNode.SelectNodes("*", "performer")
		for _, p := range performers {
			performer := decodePerformer(p)
			encounter.Performers = append(encounter.Performers, performer)
		}

		ccd.Encounters = append(ccd.Encounters, encounter)
	}

	return nil
}

func decodePerformer(node *xmlx.Node) Performer {
	p := Performer{}
	if code := Nget(node, "assignedEntity", "code"); code != nil {
		p.Code.decode(code)
	}
	if aNode := node.SelectNode("*", "addr"); aNode != nil {
		p.Address = decodeAddress(aNode)
	}
	if tNode := node.SelectNode("*", "telecom"); tNode != nil {
		p.Telecom = decodeTelecom(tNode)
	}

	//Parse name -- lifted from our patient parsing code.  Might refactor this into its own function later.

	if peNode := node.SelectNode("*", "assignedPerson"); peNode != nil {
		for n, nameNode := range peNode.SelectNodesDirect("*", "name") {
			given := nameNode.SelectNodesDirect("*", "given")
			// This is a NickName if it's the second <name><given> tag block or the
			// given tag has the qualifier CM.
			if n == 1 || (len(given) > 0 && given[0].As("*", "qualifier") == "CM") {
				p.Name.NickName = given[0].GetValue()
				continue
			}

			p.Name.Type = nameNode.As("*", "use")
			if len(given) > 0 {
				p.Name.First = given[0].GetValue()
			}
			if len(given) > 1 {
				p.Name.Middle = given[1].GetValue()
			}
			p.Name.Last = nameNode.S("*", "family")
			p.Name.Prefix = nameNode.S("*", "prefix")
			suffixes := nameNode.SelectNodes("*", "suffix")
			for n, suffix := range suffixes {
				// if it's the second suffix, or it has the qualifier TITLE
				if n == 1 || (len(p.Name.Prefix) == 0 && suffix.As("*", "qualifier") == "TITLE") {
					p.Name.Prefix = suffix.GetValue()
				} else {
					p.Name.Suffix = suffix.GetValue()
				}
			}
		}

	}
	return p
}

func decodeDiagnosis(node *xmlx.Node) Diagnosis {
	d := Diagnosis{}
	if code := Nget(node, "code"); code != nil {
		d.Code.decode(code)
	}
	d.Status = Nget(node, "statusCode").As("*", "code")

	if problem := Nget(node, "observation"); problem != nil {
		problem := decodeProblem(problem)
		if problem != nil {
			d.Problem = *problem
		}
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
