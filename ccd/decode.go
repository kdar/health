package ccd

import (
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"io"
	"strings"
	"time"
)

// Template IDs
var (
	AllergiesTid     = []string{"2.16.840.1.113883.10.20.1.2", "2.16.840.1.113883.10.20.22.2.6.1"}
	ImmunizationsTid = []string{"2.16.840.1.113883.10.20.1.6", "2.16.840.1.113883.10.20.22.2.2.1"}
	MedicationsTid   = []string{"2.16.840.1.113883.10.20.1.8", "2.16.840.1.113883.10.20.22.2.1.1"}
	ProblemsTid      = []string{"2.16.840.1.113883.10.20.1.11", "2.16.840.1.113883.10.20.22.2.5.1"}
	ResultsTid       = []string{"2.16.840.1.113883.10.20.1.14", "2.16.840.1.113883.10.20.22.2.3.1"}
	VitalSignsTid    = []string{"2.16.840.1.113883.10.20.1.16", "2.16.840.1.113883.10.20.22.2.4.1"}
)

var (
	PatientParser = Parser{
		Type:     PARSE_DOC,
		Priority: 0,
		Func:     parsePatient,
	}

	AllergiesParser = Parser{
		Type:     PARSE_SECTION,
		Values:   AllergiesTid,
		Priority: 0,
		Func:     parseAllergies,
	}

	ImmunizationsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   ImmunizationsTid,
		Priority: 0,
		Func:     parseImmunizations,
	}

	MedicationsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   MedicationsTid,
		Priority: 0,
		Func:     parseMedications,
	}

	ProblemsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   ProblemsTid,
		Priority: 0,
		Func:     parseProblems,
	}

	ResultsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   ResultsTid,
		Priority: 0,
		Func:     parseResults,
	}

	VitalSignsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   VitalSignsTid,
		Priority: 0,
		Func:     parseVitalSigns,
	}

	DefaultParsers = []Parser{
		PatientParser, AllergiesParser, ImmunizationsParser,
		MedicationsParser, ProblemsParser,
		ResultsParser, VitalSignsParser,
	}
)

type CCD struct {
	Patient       Patient
	Immunizations []Immunization
	Medications   []Medication
	Problems      []Problem
	Results       []Result
	VitalSigns    []VitalSign
	Allergies     []Allergy
	Extra         interface{}

	// Right now doc_parsers will only have one map entry "*"
	doc_parsers     map[string]Parsers
	section_parsers map[string]Parsers
}

// New CCD object with no parsers. Use NewDefaultCCD()
// or add your own parsers with AddParsers() if you want
// to actually parse anything.
func NewCCD() *CCD {
	c := &CCD{}
	c.doc_parsers = make(map[string]Parsers)
	c.section_parsers = make(map[string]Parsers)
	return c
}

// New CCD object with all the default parsers
func NewDefaultCCD() *CCD {
	c := NewCCD()
	c.AddParsers(DefaultParsers...)
	return c
}

func (c *CCD) AddParsers(parsers ...Parser) {
	for _, p := range parsers {
		if p.Organization == "" {
			p.Organization = "*"
		}

		p.Organization = strings.ToLower(p.Organization)

		if p.Type == PARSE_DOC {
			if p.Values == nil {
				p.Values = []string{"*"}
			}

			for _, v := range p.Values {
				c.doc_parsers[v] = insertSortParser(p, c.doc_parsers[v])
			}
		} else if p.Type == PARSE_SECTION {
			if p.Values == nil {
				panic("ccd: Section parser cannot have an empty Value.")
			}

			for _, v := range p.Values {
				c.section_parsers[v] = insertSortParser(p, c.section_parsers[v])
			}
		} else {
			panic("ccd: Unknown parser type.")
		}
	}
}

func (c *CCD) ParseFile(filename string) error {
	doc := xmlx.New()
	err := doc.LoadFile(filename, nil)
	if err != nil {
		return err
	}

	return c.ParseDoc(doc)
}

func (c *CCD) ParseStream(r io.Reader) error {
	doc := xmlx.New()
	err := doc.LoadStream(r, nil)
	if err != nil {
		return err
	}

	return c.ParseDoc(doc)
}

func (c *CCD) Parse(data []byte) error {
	doc := xmlx.New()
	err := doc.LoadBytes(data, nil)
	if err != nil {
		return err
	}

	return c.ParseDoc(doc)
}

type ParseType int

const (
	PARSE_DOC ParseType = iota
	PARSE_SECTION
)

type ParseFunc func(root *xmlx.Node, ccd *CCD) []error

type Parser struct {
	Type         ParseType
	Values       []string
	Organization string
	Priority     int
	Func         ParseFunc
}

type Parsers []Parser

// Parses a CCD into a CCD struct.
func (c *CCD) ParseDoc(doc *xmlx.Document) error {
	var errs_ []error
	// var errs []error

	// Reset any data retrieved from another parse
	c.Patient = Patient{}
	c.Medications = nil
	c.Problems = nil
	c.VitalSigns = nil
	c.Immunizations = nil
	c.Extra = nil

	org := Nget(doc.Root, "recordTarget", "providerOrganization", "name")
	orgName := "*"
	if org != nil {
		orgName = strings.ToLower(org.S("*", "name"))
	}

	for _, p := range c.doc_parsers["*"] {
		if orgName == "*" || p.Organization == "*" || orgName == p.Organization {
			errs_ = p.Func(doc.Root, c)
			//errs = append(errs, errs_...)
		}
	}

	componentNode := Nget(doc.Root, "component", "structuredBody")
	if componentNode != nil {
		componentNodes := componentNode.SelectNodes("*", "component")
		for _, componentNode := range componentNodes {
			sectionNode := componentNode.SelectNode("*", "section")

			tid := templateId(sectionNode)

			if parsers, ok := c.section_parsers[tid]; ok {
				for _, p := range parsers {
					if orgName == "*" || p.Organization == "*" || orgName == p.Organization {
						errs_ = p.Func(sectionNode, c)
						//errs = append(errs, errs_...)
					}
				}
			}
		}
	}

	// disabling errors for now. may return "warnings" or something
	_ = errs_

	return nil
}

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
	Type    string // H or HP = home, TMP = temporary, WP = work/office
}

func (a Address) IsZero() bool {
	return a == (Address{})
}

type Patient struct {
	Name          Name
	Dob           time.Time
	Address       Address
	Gender        string
	MaritalStatus string
	RaceCode      string
	EthnicityCode string
}

func (p Patient) IsZero() bool {
	return p == (Patient{})
}

// parses patient information from the CCD and returns
// a Patient struct
func parsePatient(root *xmlx.Node, ccd *CCD) []error {
	anode := Nget(root, "ClinicalDocument", "recordTarget", "patientRole", "addr")
	// address isn't always present
	if anode != nil {
		ccd.Patient.Address.Type = anode.As("*", "use")
		lines := anode.SelectNodes("*", "streetAddressLine")
		if len(lines) > 0 {
			ccd.Patient.Address.Line1 = lines[0].Value
		}
		if len(lines) > 1 {
			ccd.Patient.Address.Line2 = lines[1].Value
		}
		ccd.Patient.Address.City = anode.S("*", "city")
		ccd.Patient.Address.County = anode.S("*", "county")
		ccd.Patient.Address.State = anode.S("*", "state")
		ccd.Patient.Address.Zip = anode.S("*", "postalCode")
		ccd.Patient.Address.Country = anode.S("*", "country")
	}

	pnode := Nget(root, "ClinicalDocument", "recordTarget", "patientRole", "patient")
	if pnode == nil {
		return []error{
			fmt.Errorf("Could not find the node in CCD: ClinicalDocument/recordTarget/patientRole/patient"),
		}
	}

	for n, nameNode := range pnode.SelectNodes("*", "name") {
		given := nameNode.SelectNodes("*", "given")
		// This is a NickName if it's the second <name><given> tag block or the
		// given tag has the qualifier CM.
		if n == 1 || (len(given) > 0 && given[0].As("*", "qualifier") == "CM") {
			ccd.Patient.Name.NickName = given[0].Value
			continue
		}

		ccd.Patient.Name.Type = nameNode.As("*", "use")
		if len(given) > 0 {
			ccd.Patient.Name.First = given[0].Value
		}
		if len(given) > 1 {
			ccd.Patient.Name.Middle = given[1].Value
		}
		ccd.Patient.Name.Last = nameNode.S("*", "family")
		ccd.Patient.Name.Prefix = nameNode.S("*", "prefix")
		suffixes := nameNode.SelectNodes("*", "suffix")
		for n, suffix := range suffixes {
			// if it's the second suffix, or it has the qualifier TITLE
			if n == 1 || (len(ccd.Patient.Name.Prefix) == 0 && suffix.As("*", "qualifier") == "TITLE") {
				ccd.Patient.Name.Prefix = suffix.Value
			} else {
				ccd.Patient.Name.Suffix = suffix.Value
			}
		}
	}

	birthNode := pnode.SelectNode("*", "birthTime")
	if birthNode != nil {
		ccd.Patient.Dob, _ = ParseHL7Time(birthNode.As("*", "value"))
	}

	genderNode := pnode.SelectNode("*", "administrativeGenderCode")
	if genderNode != nil && genderNode.As("*", "codeSystem") == "2.16.840.1.113883.5.1" {
		switch genderNode.As("*", "code") {
		case "M":
			ccd.Patient.Gender = "Male"
		case "F":
			ccd.Patient.Gender = "Female"
		case "UN":
			ccd.Patient.Gender = "Undifferentiated"
		default:
			ccd.Patient.Gender = "Unknown"
		}
	}

	maritalNode := pnode.SelectNode("*", "maritalStatusCode")
	if maritalNode != nil && maritalNode.As("*", "codeSystem") == "2.16.840.1.113883.5.2" {
		ccd.Patient.MaritalStatus = maritalNode.As("*", "code")
	}

	raceNode := pnode.SelectNode("*", "raceCode")
	if raceNode != nil && raceNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
		ccd.Patient.RaceCode = raceNode.As("*", "code")
	}

	ethnicNode := pnode.SelectNode("*", "ethnicGroupCode")
	if ethnicNode != nil && ethnicNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
		ccd.Patient.EthnicityCode = ethnicNode.As("*", "code")
	}

	return nil
}

type Allergy struct {
	Name         string
	Reaction     string
	Status       string
	SeverityCode string
	SeverityText string
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

		// Sometimes the substance is represented here
		valueNode := Nget(obvNode, "value")
		if valueNode != nil && valueNode.As("*", "codeSystem") == "2.16.840.1.113883.6.88" {
			allergy.Name = valueNode.As("*", "displayName")
		}

		// Try to get substance another way
		if len(allergy.Name) == 0 {
			playNode := Nget(obvNode, "participant", "participantRole", "playingEntity")
			if playNode != nil {
				codeNode := Nget(playNode, "code")
				if codeNode != nil {
					allergy.Name = codeNode.As("*", "displayName")
				}

				if len(allergy.Name) == 0 {
					nameNode := Nget(playNode, "name")
					if nameNode != nil {
						allergy.Name = nameNode.S("*", "name")
					}
				}
			}
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
				}

			// Could be a Status or Severity -- Subject
			case erNode.As("*", "typeCode") == "SUBJ":
				if codeNode.As("*", "code") == "33999-4" { // Status
					allergy.Status = valueNode.As("*", "displayName")
				} else if codeNode.As("*", "code") == "SEV" { // Severity
					allergy.SeverityCode = valueNode.As("*", "code")
					allergy.SeverityText = valueNode.As("*", "displayName")
				}
			}
		}

		ccd.Allergies = append(ccd.Allergies, allergy)
	}

	return errs
}

type MedicationId struct {
	Type  string
	Value string
}

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
	DisplayName    string
	Administration string
	//Instructions   string // this is calulated and not specifically in the CCD
	Dose      MedicationDose
	Status    string
	StartDate time.Time
	StopDate  time.Time
	Period    time.Duration
	Id        MedicationId
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
			t := ParseTimeNode(etimeNode)
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
			codeSystem := codeNode.As("*", "codeSystem")
			var err error
			medication.Id.Type, err = codeSystemToMedType(codeSystem)
			if err != nil {
				// Sometimes the attributes for "code" are completely missing.
				// try to see if there is a translation node and get it from there
				transNode := codeNode.SelectNode("*", "translation")
				if transNode != nil {
					codeSystem = transNode.As("*", "codeSystem")
					var err2 error
					medication.Id.Type, err2 = codeSystemToMedType(codeSystem)
					if err2 != nil {
						errs = append(errs, err)
					}
				} else {
					errs = append(errs, err)
				}
			}
		}
		medication.Id.Value = codeNode.As("*", "code")

		if displayName := codeNode.As("*", "displayName"); displayName != "" {
			medication.Name = displayName
			medication.DisplayName = displayName
		}

		if nameNode := manNode.SelectNode("*", "name"); nameNode != nil {
			medication.Name = nameNode.Value
		} else if originalNode := codeNode.SelectNode("*", "originalText"); originalNode != nil {
			medication.Name = originalNode.Value
		}

		ccd.Medications = append(ccd.Medications, medication)
	}

	return errs
}

type Problem struct {
	Name string
	Date time.Time
	// Duration time.Duration
	// Status   string
}

func parseProblems(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		problem := Problem{}

		observationNode := Nget(entryNode, "act", "entryRelationship", "observation")
		problem.Name = Nget(observationNode, "value").As("*", "displayName")

		effectiveTimeNode := Nget(observationNode, "effectiveTime")
		t := ParseTimeNode(effectiveTimeNode)
		problem.Date = t.Value

		// observationNode2 := Nget(observationNode, "entryRelationship", "observation")
		// if observationNode2 != nil {
		//   problem.Status = Nget(observationNode2, "value").As("*", "displayName")
		// }

		ccd.Problems = append(ccd.Problems, problem)
	}

	return nil
}

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

		t := ParseTimeNode(Nget(saNode, "effectiveTime"))
		immunization.Date = t.Value

		routeCodeNode := Nget(saNode, "routeCode")
		if routeCodeNode != nil {
			if routeCodeNode.As("*", "codeSystemName") == "RouteOfAdministration" {
				immunization.Administration = routeCodeNode.As("*", "displayName")
			}
		}

		codeNode := Nget(saNode, "manufacturedProduct", "manufacturedMaterial", "code")
		if codeNode != nil {
			immunization.Name = codeNode.As("*", "displayName")
		}

		ccd.Immunizations = append(ccd.Immunizations, immunization)
	}

	return nil
}

type ResultValue struct {
	Type  string
	Value string
	Unit  string
}

type ResultObservation struct {
	Date                time.Time
	DisplayName         string
	Loinc               string
	Value               ResultValue
	InterpretationCodes []string
}

type Result struct {
	Date         time.Time
	Observations []ResultObservation
}

func parseResults(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		organizerNode := Nget(entryNode, "organizer")
		if organizerNode == nil {
			return nil
		}

		componentNodes := organizerNode.SelectNodes("*", "component")
		if componentNodes == nil {
			return nil
		}

		result := Result{}

		effectiveTimeNode := Nget(organizerNode, "effectiveTime")
		t := ParseTimeNode(effectiveTimeNode)
		result.Date = t.Value

		for _, componentNode := range componentNodes {
			obNode := Nget(componentNode, "observation")
			if obNode == nil {
				continue
			}

			observation := ResultObservation{}

			effectiveTimeNode := Nget(obNode, "effectiveTime")
			t = ParseTimeNode(effectiveTimeNode)
			observation.Date = t.Value

			codeNode := Nget(obNode, "code")
			if codeNode != nil {
				observation.DisplayName = codeNode.As("*", "displayName")
				observation.Loinc = codeNode.As("*", "code")
			}

			valueNode := Nget(obNode, "value")
			if valueNode != nil {
				observation.Value.Type = valueNode.As("*", "type")
				observation.Value.Value = valueNode.As("*", "value")
				observation.Value.Unit = valueNode.As("*", "unit")
			}

			icodeNodes := obNode.SelectNodes("*", "interpretationCode")
			if icodeNodes != nil {
				for _, icodeNode := range icodeNodes {
					observation.InterpretationCodes = append(observation.InterpretationCodes, icodeNode.As("*", "code"))
				}
			}

			result.Observations = append(result.Observations, observation)
		}

		ccd.Results = append(ccd.Results, result)
	}

	return nil
}
