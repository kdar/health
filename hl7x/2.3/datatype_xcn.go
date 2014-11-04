package hl7v2_3

type Xcn struct {
	// String number (ST)
	IDNumberST String `position:"XCN.1"`
	// family name
	FamilyName String `position:"XCN.2"`
	// given name
	GivenName String `position:"XCN.3"`
	// middle initial or name
	MiddleInitialOrName String `position:"XCN.4"`
	// suffix (e.g. JR or III)
	Suffix String `position:"XCN.5"`
	// prefix (e.g. DR)
	Prefix String `position:"XCN.6"`
	// degree (e.g. MD)
	Degree String `position:"XCN.7"`
	// source table
	SourceTable String `position:"XCN.8"`
	// assigning authority
	AssigningAuthority Hd `position:"XCN.9"`
	// name type
	NameType String `position:"XCN.10"`
	// identifier check digit
	IdentifierCheckDigit String `position:"XCN.11"`
	// code identifying the check digit scheme employed
	CodeIdentifyingTheCheckDigitSchemeEmployed String `position:"XCN.12"`
	// identifier type code
	IdentifierTypeCode String `position:"XCN.13"`
	// assigning facility ID
	AssigningFacilityID Hd `position:"XCN.14"`
}
