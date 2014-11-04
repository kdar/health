package hl7v2_3

type Xpn struct {
	// family name
	FamilyName String `position:"XPN.1"`
	// given name
	GivenName String `position:"XPN.2"`
	// middle initial or name
	MiddleInitialOrName String `position:"XPN.3"`
	// suffix (e.g. JR or III)
	Suffix String `position:"XPN.4"`
	// prefix (e.g. DR)
	Prefix String `position:"XPN.5"`
	// degree (e.g. MD)
	Degree String `position:"XPN.6"`
	// name type code
	NameTypeCode String `position:"XPN.7"`
	// Name Representation code
	NameRepresentationCode String `position:"XPN.8"`
}
