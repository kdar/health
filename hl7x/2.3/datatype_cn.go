package hl7v2_3

type Cn struct {
	// String number (ST)
	IdNumberSt String `position:"CN.1"`
	// family name
	FamilyName String `position:"CN.2"`
	// given name
	GivenName String `position:"CN.3"`
	// middle initial or name
	MiddleInitialOrName String `position:"CN.4"`
	// suffix (e.g. JR or III)
	Suffix String `position:"CN.5"`
	// prefix (e.g. DR)
	Prefix String `position:"CN.6"`
	// degree (e.g. MD)
	Degree String `position:"CN.7"`
	// source table
	SourceTable String `position:"CN.8"`
	// assigning authority
	AssigningAuthority Hd `position:"CN.9"`
}
