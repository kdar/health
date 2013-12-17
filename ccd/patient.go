package ccd

import (
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"strings"
	"time"
)

var (
	PatientParser = Parser{
		Type:     PARSE_DOC,
		Priority: 0,
		Func:     parsePatient,
	}
)

type Name struct {
	Last     string
	First    string
	Middle   string
	Suffix   string
	Prefix   string // title
	Type     string // L = legal name, PN = patient name (not sure)
	NickName string
}

func (n Name) IsZero() bool {
	return n == (Name{})
}

type Address struct {
	Line1   string
	Line2   string
	City    string
	County  string
	State   string
	Zip     string
	Country string
	Use     string // H or HP = home, TMP = temporary, WP = work/office
}

func (a Address) IsZero() bool {
	return a == (Address{})
}

type Telecom struct {
	Type          string
	Use           string
	Value         string
	OriginalValue string
}

func decodeTelecom(n *xmlx.Node) (t Telecom) {
	if n == nil {
		return t
	}

	t.OriginalValue = n.As("*", "value")
	t.Use = n.As("*", "use")

	coloni := strings.Index(t.OriginalValue, ":")
	if coloni > 0 {
		switch t.OriginalValue[:coloni] {
		case "tel":
			t.Type = "phone"
		case "http":
			t.Type = "url"
		case "mailto":
			t.Type = "email"
		}
		t.Value = t.OriginalValue[coloni+1:]
	}

	return t
}

type Patient struct {
	Name          Name
	Dob           time.Time
	Addresses     []Address
	Telecoms      []Telecom
	LanguageCode  string
	Gender        Code
	MaritalStatus Code
	Race          Code
	Ethnicity     Code
	Religion      Code
}

// func (p Patient) IsZero() bool {
// 	return p == (Patient{})
// }

// parses patient information from the CCD and returns
// a Patient struct
func parsePatient(root *xmlx.Node, ccd *CCD) []error {
	prNode := Nget(root, "ClinicalDocument", "recordTarget", "patientRole")
	if prNode == nil {
		return []error{
			fmt.Errorf("Could not find the node in CCD: ClinicalDocument/recordTarget/patientRole"),
		}
	}

	aNodes := prNode.SelectNodes("*", "addr")
	if aNodes != nil {
		for _, anode := range aNodes {
			address := Address{}

			address.Use = anode.As("*", "use")
			lines := anode.SelectNodes("*", "streetAddressLine")
			if len(lines) > 0 {
				address.Line1 = lines[0].GetValue()
			}
			if len(lines) > 1 {
				address.Line2 = lines[1].GetValue()
			}
			address.City = anode.S("*", "city")
			address.County = anode.S("*", "county")
			address.State = anode.S("*", "state")
			address.Zip = anode.S("*", "postalCode")
			address.Country = anode.S("*", "country")

			ccd.Patient.Addresses = append(ccd.Patient.Addresses, address)
		}
	}

	teleNodes := prNode.SelectNodes("*", "telecom")
	if teleNodes != nil {
		for _, tnode := range teleNodes {
			ccd.Patient.Telecoms = append(ccd.Patient.Telecoms, decodeTelecom(tnode))
		}
	}

	pnode := Nget(prNode, "patient")
	if pnode == nil {
		return []error{
			fmt.Errorf("Could not find the node in CCD: ClinicalDocument/recordTarget/patientRole/patient"),
		}
	}

	languageNode := Nget(pnode, "languageCommunication", "languageCode")
	if languageNode != nil {
		ccd.Patient.LanguageCode = languageNode.As("*", "code")
	}

	for n, nameNode := range pnode.SelectNodes("*", "name") {
		given := nameNode.SelectNodes("*", "given")
		// This is a NickName if it's the second <name><given> tag block or the
		// given tag has the qualifier CM.
		if n == 1 || (len(given) > 0 && given[0].As("*", "qualifier") == "CM") {
			ccd.Patient.Name.NickName = given[0].GetValue()
			continue
		}

		ccd.Patient.Name.Type = nameNode.As("*", "use")
		if len(given) > 0 {
			ccd.Patient.Name.First = given[0].GetValue()
		}
		if len(given) > 1 {
			ccd.Patient.Name.Middle = given[1].GetValue()
		}
		ccd.Patient.Name.Last = nameNode.S("*", "family")
		ccd.Patient.Name.Prefix = nameNode.S("*", "prefix")
		suffixes := nameNode.SelectNodes("*", "suffix")
		for n, suffix := range suffixes {
			// if it's the second suffix, or it has the qualifier TITLE
			if n == 1 || (len(ccd.Patient.Name.Prefix) == 0 && suffix.As("*", "qualifier") == "TITLE") {
				ccd.Patient.Name.Prefix = suffix.GetValue()
			} else {
				ccd.Patient.Name.Suffix = suffix.GetValue()
			}
		}
	}

	birthNode := pnode.SelectNode("*", "birthTime")
	if birthNode != nil {
		ccd.Patient.Dob, _ = ParseHL7Time(birthNode.As("*", "value"))
	}

	genderNode := pnode.SelectNode("*", "administrativeGenderCode")
	if genderNode != nil && genderNode.As("*", "codeSystem") == "2.16.840.1.113883.5.1" {
		ccd.Patient.Gender.decode(genderNode)
	}

	maritalNode := pnode.SelectNode("*", "maritalStatusCode")
	if maritalNode != nil && maritalNode.As("*", "codeSystem") == "2.16.840.1.113883.5.2" {
		ccd.Patient.MaritalStatus.decode(maritalNode)
	}

	raceNode := pnode.SelectNode("*", "raceCode")
	if raceNode != nil && raceNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
		ccd.Patient.Race.decode(raceNode)
	}

	ethnicNode := pnode.SelectNode("*", "ethnicGroupCode")
	if ethnicNode != nil && ethnicNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
		ccd.Patient.Ethnicity.decode(ethnicNode)
	}

	religionNode := pnode.SelectNode("*", "religiousAffiliationCode")
	if religionNode != nil {
		ccd.Patient.Religion.decode(religionNode)
	}

	// spew.Dump(ccd.Patient)

	return nil
}
