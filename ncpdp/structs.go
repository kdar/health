package ncpdp

import (
	"github.com/kdar/health/edifact"
	"time"
)

type Name struct {
	Last   string
	First  string
	Middle string
	Suffix string
	Prefix string
}

type Address struct {
	Line1             string
	City              string
	State             string
	Postal            string
	LocationQualifier string
	Location          string // used as line2 sometimes
}

type Phone struct {
	Number    string
	Qualifier string
}

// segment: PTT
type Patient struct {
	Relationship string // 010 C

	// format: CCYYMMDD
	Dob time.Time // 020 C

	// made up of the components:
	// last, first, middle, suffix, prefix
	Name *Name // 030-* M

	// Values:
	// M = Male
	// F = Female
	// U = Unknown
	Gender string // 040 C

	ReferenceNumber string // 050-00 M

	// default: 1D
	ReferenceQualifier string // 050-01 C

	// made up of the components:
	// line1, city, state, postal, location qualifier, location
	Address *Address // 060-* C

	// made up of the components:
	// number, qualifier
	// repeated n times
	Phones []*Phone // 070-* C
}

func (p *Patient) fill(values edifact.Values) {
	p.Relationship = getString(values, 1)

	dob := getString(values, 2)
	if len(dob) > 0 {
		dobTime, err := time.Parse("20060102", dob)
		if err == nil {
			p.Dob = dobTime
		}
	}

	p.Name = getName(values, 3)
	p.Gender = getString(values, 4)

	subValues := getValues(values, 5)
	p.ReferenceNumber = getString(subValues, 0)
	p.ReferenceQualifier = getString(subValues, 1)

	p.Address = getAddress(values, 6)
	p.Phones = getPhones(values, 7)
}

// creates a new patient
func newPatient() *Patient {
	return &Patient{}
}

// type ProviderID struct {
//   Value     string
//   Qualifier string
// }

// segment: PVD
// defines a prescriber, pharmacy/pharmacist, and supervisor
type Provider struct {
	// ----- Provider code 010 M

	ProviderCode string // 010 M

	// ----- Reference number 020 C

	ReferenceNumber string // 020-00 M

	// Values:
	//   ZZ = NPI (not totally sure)
	//   HPI = NPI
	//   DH = DEA number
	//   DE = DEA number
	//   D3 = NCPDP Provider ID Number
	//   0B = State license number
	ReferenceQualifier string // 020-01 C

	// ----- Name 050 C
	// name of prescriber, pharmacist, or supervisor

	// made up of the components:
	// last, first, middle, suffix, prefix
	Name *Name // 050-* C

	// ----- Pharmacy/Clinic name 070 C

	PartyName string // 070 C
}

func newProvider() *Provider {
	return &Provider{}
}

func (p *Provider) fill(values edifact.Values) {
	subValues := getValues(values, 2)
	p.ReferenceNumber = getString(subValues, 0)
	p.ReferenceQualifier = getString(subValues, 1)

	p.Name = getName(values, 5)

	p.PartyName = getString(values, 7)
}

// // returns the value and qualifier of the provider id
// func (p *Provider) ProviderID() *ProviderID {
//   return &ProviderID{p.ReferenceNumber, p.ReferenceQualifier}
// }
