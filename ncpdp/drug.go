package ncpdp

import (
  //"log"
  "github.com/kdar/health/edifact"
  "strconv"
  "time"
)

// type DrugId struct {
//   Value     string
//   Qualifier string
// }

// type DrugQuantity struct {
//   Value     string
//   Qualifier string
//   Code      string
// }

type Drug struct {
  // ----- Drug 010 M

  // Definition of the loop of the DRU segment.
  // Values:
  //   P = Prescribed
  //   D = Dispensed
  //   R = Requested
  ItemDescriptionIdentification string // 010-00 M
  // Drug name
  ItemDescription string // 010-01 M
  // Drug id  (NDCs, rxnorms, etc...)
  ItemNumber string // 010-02 C
  // The drug id qualifier.
  // Values:
  //   ND = NDC11
  //   MF = Manufacturing
  CodeListResponsibilityAgency string // 010-03 C

  // more fields are in this section, but unimplemented

  // ----- Quantity 020 C

  // unit of measure
  QuantityQualifier string // 020-00 M
  Quantity          string // 020-01 M
  CodeListQualifier string // 020-02 C

  // ----- Directions 030 C

  DosageId string // 030-00 N/U
  Dosage1  string // 030-01 C
  Dosage2  string // 030-02 C

  // ----- Date 040 M
  // this section is an array of values that contains:
  //  [date period qualifier, date/quantity, date format qualifier], ...
  // can repeat 5 times
  // explanation of date format qualifiers:
  //   102 = CCYYMMDD
  //   108 = Quantity of Days
  //   804 = Quantity of Days

  DaysSupply *time.Duration // qualifier: ZDS
  DateIssued *time.Time     // qualifier: 85
  // last dispensed or last filled
  LastDemand *time.Time // qualifier: LD
  // ExpirationDate *time.Time // qualifier: 36
  // EffectiveDate *time.time // qualifier: 07
  // PeriodEnd *time.Time // qualifier: PE

  // ----- Substitution 050 C

  // Product/service substitution, coded
  // Values:
  // 0 = No Product Selection Indicated
  // 1 = Substitution Not Allowed by Prescriber
  // 2 = Substitution Allowed - Patient Requested Product Dispensed
  // 3 = Substitution Allowed - Pharmacist Selected Product Dispensed
  // 4 = Substitution Allowed - Generic Drug Not in Stock
  // 5 = Substitution Allowed - Brand Drug Dispensed as a Generic
  // 7 = Substitution Not Allowed - Brand Drug Mandated by Law
  // 8 = Substitution Allowed - Generic Drug Not Available in Marketplace
  // (6 was intentionally left off)
  Substitution string // 050 C

  // ----- not part of DRU and not parsed here

  Prescriber *Provider
  Pharmacy   *Provider
}

func newDrug() *Drug {
  return &Drug{}
}

func (d *Drug) fill(values edifact.Values) {
  subValues := getValues(values, 1)
  d.ItemDescriptionIdentification = getString(subValues, 0)
  d.ItemDescription = getString(subValues, 1)
  d.ItemNumber = getString(subValues, 2)
  d.CodeListResponsibilityAgency = getString(subValues, 3)

  subValues = getValues(values, 2)
  d.QuantityQualifier = getString(subValues, 0)
  d.Quantity = getString(subValues, 1)
  d.CodeListQualifier = getString(subValues, 2)

  subValues = getValues(values, 3)
  d.DosageId = getString(subValues, 0)
  d.Dosage1 = getString(subValues, 1)
  d.Dosage2 = getString(subValues, 2)

  subValues = getValues(values, 4)
  for n := range subValues {
    subVals := getValues(subValues, n)
    if len(subVals) == 3 {
      var date time.Time
      var duration time.Duration
      switch subVals[2] {
      case "102":
        date, _ = time.Parse("20060102", getString(subVals, 1))
      case "108":
        fallthrough
      case "804":
        d, _ := strconv.ParseInt(getString(subVals, 1), 10, 64)
        duration = time.Duration(d) * time.Hour * 24
        // default:
        //   log.Printf("unknown format: %s", subVals[2])
      }

      switch subVals[0] {
      case "ZDS":
        d.DaysSupply = &duration
      case "LD":
        d.LastDemand = &date
      case "85":
        d.DateIssued = &date
        // default:
        //   log.Printf("unknown qualifier: %s", subVals[0])
      }
    }
  }

  d.Substitution = getString(values, 5)
}

// return the name of the drug
func (d *Drug) Name() string {
  return d.ItemDescription
}

// // returns the value and qualifier of the drug id
// func (d *Drug) DrugId() *DrugId {
//   return &DrugID{d.ItemNumber, d.CodeListResponsibilityAgency}
// }

// // returns the quantity, qualifier, and the code list qualifier
// func (d *Drug) DrugQuantity() *DrugQuantity {
//   return &DrugQuantity{d.Quantity, d.QuantityQualifier, d.CodeListQualifier}
// }

// returns whether the drug is prescribed
func (d *Drug) Prescribed() bool {
  return d.ItemDescriptionIdentification == "P"
}

// returns whether the drug is dispensed
func (d *Drug) Dispensed() bool {
  return d.ItemDescriptionIdentification == "D"
}

// returns whether the drug is requested
func (d *Drug) Requested() bool {
  return d.ItemDescriptionIdentification == "R"
}
