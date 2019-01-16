package ccd

import (
	"github.com/mattn/go-pkg-xmlx"
)

var (
	AllergiesTid = []string{"2.16.840.1.113883.10.20.1.2", "2.16.840.1.113883.10.20.22.2.6.1"}

	AllergiesParser = Parser{
		Type:     PARSE_SECTION,
		Values:   AllergiesTid,
		Priority: 0,
		Func:     parseAllergies,
	}
)

type Allergy struct {
	Name         string
	Reaction     string
	Status       string
	SeverityCode string
	SeverityText string

	Type      Code
	Substance Code
	Severity  Code
}

func parseAllergies(node *xmlx.Node, ccd *CCD) []error {
	var errs []error

	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		obvNode := Nget(entryNode, "act", "entryRelationship", "observation")
		if obvNode == nil {
			continue
		}

		allergy := Allergy{}
		allergy.Type.decode(Nget(obvNode, "value"))

		// Sometimes the substance is represented here
		valueNode := Nget(obvNode, "value")
		if valueNode != nil && valueNode.As("*", "codeSystem") == "2.16.840.1.113883.6.88" {
			allergy.Name = valueNode.As("*", "displayName") //TODO: should string format a name later instead of just using the substance name
			allergy.Substance.decode(valueNode)
		}

		// Try to get substance another way
		if len(allergy.Name) == 0 {
			playNode := Nget(obvNode, "participant", "participantRole", "playingEntity")
			if playNode != nil {
				codeNode := Nget(playNode, "code")
				if codeNode != nil {
					allergy.Name = codeNode.As("*", "displayName")
					allergy.Substance.decode(codeNode)
				}

				if len(allergy.Name) == 0 {
					nameNode := Nget(playNode, "name")
					if nameNode != nil {
						allergy.Name = nameNode.S("*", "name")
					}
				}
			}
		}

		// If we still don't have a name, just continue because there's no
		// point in processing this allergy without a name.
		if len(allergy.Name) == 0 {
			continue
		}

		erNodes := obvNode.SelectNodes("*", "entryRelationship")
		for _, erNode := range erNodes {
			oNode := Nget(erNode, "observation")
			codeNode := Nget(oNode, "code")
			valueNode := Nget(oNode, "value")

			switch {
			// Reaction -- "Manifestation"
			case erNode.As("*", "typeCode") == "MFST":
				allergy.Reaction = valueNode.As("*", "displayName")

				// Sometimes severity is a child of this observation
				suboNode := Nget(oNode, "entryRelationship", "observation")
				if suboNode != nil {
					subValueNode := Nget(suboNode, "value")
					allergy.SeverityCode = subValueNode.As("*", "code")
					allergy.SeverityText = subValueNode.As("*", "displayName")
					allergy.Severity.decode(valueNode)
				}

			// Could be a Status or Severity -- Subject
			case erNode.As("*", "typeCode") == "SUBJ":
				if codeNode != nil {
					if codeNode.As("*", "code") == "33999-4" { // Status
						allergy.Status = valueNode.As("*", "displayName")
					} else if codeNode.As("*", "code") == "SEV" { // Severity
						allergy.SeverityCode = valueNode.As("*", "code")
						allergy.SeverityText = valueNode.As("*", "displayName")
						allergy.Severity.decode(valueNode)
					}
				}
			}
		}

		ccd.Allergies = append(ccd.Allergies, allergy)
	}

	return errs
}
