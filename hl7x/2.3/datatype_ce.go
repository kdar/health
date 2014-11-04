package hl7v2_3

type Ce struct {
	// identifier
	Identifier String `position:"CE.1"`
	// text
	Text String `position:"CE.2"`
	// name of coding system
	NameOfCodingSystem String `position:"CE.3"`
	// alternate identifier
	AlternateIdentifier String `position:"CE.4"`
	// alternate text
	AlternateText String `position:"CE.5"`
	// name of alternate coding system
	NameOfAlternateCodingSystem String `position:"CE.6"`
}
