package ccd

import (
  "github.com/jteeuwen/go-pkg-xmlx"
  "strings"
  "time"
)

type Problem struct {
}

type CCD struct {
  Patient  *Patient
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

func Unmarshal(data []byte) (*CCD, error) {
  doc := xmlx.New()
  err := doc.LoadBytes(data, nil)
  if err != nil {
    return nil, err
  }

  ccd := &CCD{}
  ccd.Patient = parsePatient(doc.Root)

  componentNode := nget(doc.Root, "component", "structuredBody")
  if componentNode != nil {
    componentNodes := componentNode.SelectNodes("*", "component")
    for _, componentNode := range componentNodes {
      sectionNode := componentNode.SelectNode("*", "section")
      switch templateId(sectionNode) {
      case "2.16.840.1.113883.10.20.1.8":
        ccd.Meds = parseMeds(sectionNode)
      case "2.16.840.1.113883.10.20.1.11":
        ccd.Problems = parseProblems(sectionNode)
      }
    }
  }

  return ccd, nil
}

type PatientName struct {
  Last   string
  First  string
  Middle string
  Suffix string
}

type Patient struct {
  Name          PatientName
  Dob           time.Time
  Gender        string
  MaritalStatus string
  Race          string
  Ethnicity     string
}

// parses patient information from the CCD and returns
// a Patient struct
func parsePatient(root *xmlx.Node) *Patient {
  node := nget(root, "ClinicalDocument", "recordTarget", "patientRole", "patient")
  if node == nil {
    return nil
  }

  patient := &Patient{}

  nameNode := node.SelectNode("*", "name")
  given := nameNode.SelectNodes("*", "given")
  patient.Name.First = given[0].Value
  if len(given) > 1 {
    patient.Name.Middle = given[1].Value
  }
  patient.Name.Last = nameNode.S("*", "family")
  patient.Name.Suffix = nameNode.S("*", "suffix")

  birthNode := node.SelectNode("*", "birthTime")
  if birthNode != nil {
    patient.Dob, _ = time.Parse("20060102", birthNode.As("*", "value"))
  }

  genderNode := node.SelectNode("*", "administrativeGenderCode")
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

  maritalNode := node.SelectNode("*", "maritalStatusCode")
  if maritalNode != nil && maritalNode.As("*", "codeSystem") == "2.16.840.1.113883.5.2" {
    patient.MaritalStatus = maritalNode.As("*", "code")
  }

  raceNode := node.SelectNode("*", "raceCode")
  if raceNode != nil && raceNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
    patient.Race = raceNode.As("*", "code")
  }

  ethnicNode := node.SelectNode("*", "ethnicGroupCode")
  if ethnicNode != nil && ethnicNode.As("*", "codeSystem") == "2.16.840.1.113883.6.238" {
    patient.Ethnicity = ethnicNode.As("*", "code")
  }

  return patient
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

func parseMeds(node *xmlx.Node) []Med {
  var meds []Med

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
      switch codeNode.As("*", "codeSystem") {
      case "2.16.840.1.113883.6.69": // NDC
        med.Id.Type = "NDC"
      case "2.16.840.1.113883.6.88": // RxNorm
        med.Id.Type = "RxNorm"
      default:
        // "Warning: Unknown codeSystem value of '".
        //              $mp->{manufacturedMaterial}{code}{codeSystem}
        //              ."'\n"
      }
    } else {
      // "Warning: No codeSystem supplied for med!\n"
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

  return meds
}

func parseProblems(node *xmlx.Node) []Problem {
  return nil
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
