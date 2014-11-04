package hl7v2_3

type Cx struct {
	// ID
	ID String `position:"CX.1"`
	// check digit
	CheckDigit String `position:"CX.2"`
	// code identifying the check digit scheme employed
	CodeIdentifyingTheCheckDigitSchemeEmployed String `position:"CX.3"`
	// assigning authority
	AssigningAuthority String `position:"CX.4"`
	// identifier type code
	IdentifierTypeCode String `position:"CX.5"`
	// assigning facility
	AssigningFacility String `position:"CX.6"`
}
