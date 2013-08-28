package ccd

import (
  "fmt"
  "github.com/jteeuwen/go-pkg-xmlx"
  "io"
  "strings"
  "time"
)

type Problem struct {
}

type CCD struct {
  Patient  Patient
  Meds     []Med
  Problems []Problem
}

// Node get.
// helper function to continually transverse down the
// xml nodes in args, and return the last one.
func nget(node *xmlx.Node, args ...string) *xmlx.Node {
  for _, a := range args {
    if node == nil {
      return nil
    }

    node = node.SelectNode("*", a)
  }

  return node
}

// Node Safe get.
// just like nget, but returns a node no matter what.
func nsget(node *xmlx.Node, args ...string) *xmlx.Node {
  n := nget(node, args...)
  if n == nil {
    return xmlx.NewNode(0)
  }
  return n
}

func UnmarshalFile(filename string) (*CCD, error) {
  doc := xmlx.New()
  err := doc.LoadFile(filename, nil)
  if err != nil {
    return nil, err
  }

  return UnmarshalDoc(doc)
}

func UnmarshalStream(r io.Reader) (*CCD, error) {
  doc := xmlx.New()
  err := doc.LoadStream(r, nil)
  if err != nil {
    return nil, err
  }

  return UnmarshalDoc(doc)
}

func Unmarshal(data []byte) (*CCD, error) {
  doc := xmlx.New()
  err := doc.LoadBytes(data, nil)
  if err != nil {
    return nil, err
  }

  return UnmarshalDoc(doc)
}

// Unmarshals a CCD into a CCD struct.
func UnmarshalDoc(doc *xmlx.Document) (*CCD, error) {
  var errs_ []error
  // var errs []error

  ccd := &CCD{}
  ccd.Patient, errs_ = parsePatient(doc.Root)
  //errs = append(errs, errs_...)

  componentNode := nget(doc.Root, "component", "structuredBody")
  if componentNode != nil {
    componentNodes := componentNode.SelectNodes("*", "component")
    for _, componentNode := range componentNodes {
      sectionNode := componentNode.SelectNode("*", "section")
      switch templateId(sectionNode) {
      case "2.16.840.1.113883.10.20.1.8":
        ccd.Meds, errs_ = parseMeds(sectionNode)
      case "2.16.840.1.113883.10.20.1.11":
        ccd.Problems, errs_ = parseProblems(sectionNode)
      }
      //errs = append(errs, errs_...)
    }
  }

  // disabling errors for now. may return "warnings" or something
  _ = errs_

  return ccd, nil
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
  Race          string
  Ethnicity     string
}

func (p Patient) IsZero() bool {
  return p == (Patient{})
}

// parses patient information from the CCD and returns
// a Patient struct
func parsePatient(root *xmlx.Node) (Patient, []error) {
  patient := Patient{}

  anode := nget(root, "ClinicalDocument", "recordTarget", "patientRole", "addr")
  // address isn't always present
  if anode != nil {
    patient.Address.Type = anode.As("*", "use")
    lines := anode.SelectNodes("*", "streetAddressLine")
    if len(lines) > 0 {
      patient.Address.Line1 = lines[0].Value
    }
    if len(lines) > 1 {
      patient.Address.Line2 = lines[1].Value
    }
    patient.Address.City = anode.S("*", "city")
    patient.Address.County = anode.S("*", "county")
    patient.Address.State = anode.S("*", "state")
    patient.Address.Zip = anode.S("*", "postalCode")
    patient.Address.Country = anode.S("*", "country")
  }

  pnode := nget(root, "ClinicalDocument", "recordTarget", "patientRole", "patient")
  if pnode == nil {
    return patient, []error{
      fmt.Errorf("Could not find the node in CCD: ClinicalDocument/recordTarget/patientRole/patient"),
    }
  }

  for n, nameNode := range pnode.SelectNodes("*", "name") {
    given := nameNode.SelectNodes("*", "given")
    // This is a NickName if it's the second <name><given> tag block or the
    // given tag has the qualifier CM.
    if n == 1 || (len(given) > 0 && given[0].As("*", "qualifier") == "CM") {
      patient.Name.NickName = given[0].Value
      continue
    }

    patient.Name.Type = nameNode.As("*", "use")
    if len(given) > 0 {
      patient.Name.First = given[0].Value
    }
    if len(given) > 1 {
      patient.Name.Middle = given[1].Value
    }
    patient.Name.Last = nameNode.S("*", "family")
    patient.Name.Prefix = nameNode.S("*", "prefix")
    suffixes := nameNode.SelectNodes("*", "suffix")
    for n, suffix := range suffixes {
      // if it's the second suffix, or it has the qualifier TITLE
      if n == 1 || (len(patient.Name.Prefix) == 0 && suffix.As("*", "qualifier") == "TITLE") {
        patient.Name.Prefix = suffix.Value
      } else {
        patient.Name.Suffix = suffix.Value
      }
    }
  }

  birthNode := pnode.SelectNode("*", "birthTime")
  if birthNode != nil {
    patient.Dob, _ = time.Parse("20060102", birthNode.As("*", "value"))
  }

  genderNode := pnode.SelectNode("*", "administrativeGenderCode")
  if genderNode != nil && genderNode.As("*", "codeSystem") == "2.16.840.1.113883.5.1" {
    switch genderNode.As("*", "code") {
    case "M":
      patient.Gender = "Male"
    case "F":
      patient.Gender = "Female"
    case "UN":
      patient.Gender = "Undifferentiated"
    default:
      patient.Gender = "Unknown"
    }
  }

  maritalNode := pnode.SelectNode("*", "maritalStatusCode")
  if maritalNode != nil && maritalNode.As("*", "codeSystem") == "2.16.840.1.113883.5.2" {
    patient.MaritalStatus = maritalNode.As("*", "code")
  }

  raceNode := pnode.SelectNode("*", "raceCode")
  if raceNode != nil && raceNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
    patient.Race = raceNode.As("*", "code")
  }

  ethnicNode := pnode.SelectNode("*", "ethnicGroupCode")
  if ethnicNode != nil && ethnicNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
    patient.Ethnicity = ethnicNode.As("*", "code")
  }

  return patient, nil
}

type Med struct {
  Name      string
  Status    string
  StartDate time.Time
  StopDate  time.Time

  Id struct {
    Type  string
    Value string
  }
}

func parseMeds(node *xmlx.Node) ([]Med, []error) {
  var meds []Med
  var errs []error

  entryNodes := node.SelectNodes("*", "entry")
  for _, entryNode := range entryNodes {
    if templateId(entryNode) != "2.16.840.1.113883.10.20.1.24" {
      continue
    }

    mpNode := nget(entryNode, "substanceAdministration", "consumable", "manufacturedProduct")
    if mpNode == nil {
      continue
    }

    med := Med{}
    med.Status = nsget(entryNode, "substanceAdministration", "statusCode").As("*", "code")

    etimeNodes := nget(entryNode, "substanceAdministration").SelectNodes("*", "effectiveTime")
    for _, etimeNode := range etimeNodes {
      if strings.ToLower(etimeNode.As("*", "type")) == "ivl_ts" {
        lowNode := etimeNode.SelectNode("*", "low")
        if lowNode != nil {
          med.StartDate, _ = time.Parse("20060102", lowNode.As("*", "value"))
        }

        highNode := etimeNode.SelectNode("*", "high")
        if highNode != nil {
          med.StopDate, _ = time.Parse("20060102", highNode.As("*", "value"))
        }
      }
    }

    manNode := nget(mpNode, "manufacturedMaterial")

    codeNode := nget(manNode, "code")
    if codeNode != nil {
      codeSystem := codeNode.As("*", "codeSystem")
      var err error
      med.Id.Type, err = codeSystemToMedType(codeSystem)
      if err != nil {
        // Sometimes the attributes for "code" are completely missing.
        // try to see if there is a translation node and get it from there
        transNode := codeNode.SelectNode("*", "translation")
        if transNode != nil {
          codeSystem = transNode.As("*", "codeSystem")
          var err2 error
          med.Id.Type, err2 = codeSystemToMedType(codeSystem)
          if err2 != nil {
            errs = append(errs, err)
          }
        } else {
          errs = append(errs, err)
        }
      }
    }
    med.Id.Value = codeNode.As("*", "code")

    if displayName := codeNode.As("*", "displayName"); displayName != "" {
      med.Name = displayName
    } else if nameNode := manNode.SelectNode("*", "name"); nameNode != nil {
      med.Name = nameNode.Value
    } else if originalNode := codeNode.SelectNode("*", "originalText"); originalNode != nil {
      med.Name = originalNode.Value
    }

    meds = append(meds, med)
  }

  return meds, errs
}

func parseProblems(node *xmlx.Node) ([]Problem, []error) {
  return nil, nil
}

func templateId(node *xmlx.Node) string {
  idNodes := node.SelectNodes("*", "templateId")
  id := ""
  for _, idNode := range idNodes {
    id = idNode.As("*", "root")
    if strings.HasPrefix(id, "2.16.840.1.113883.10.20.1.") {
      return id
    }
  }

  return id
}

func codeSystemToMedType(codeSystem string) (string, error) {
  switch codeSystem {
  case "2.16.840.1.113883.6.69": // NDC
    return "NDC", nil
  case "2.16.840.1.113883.6.88": // RxNorm
    return "RxNorm", nil
  }
  return "", fmt.Errorf(`Unknown med codeSystem value of "%s"`, codeSystem)
}
