package hl7v2_3

type Pt struct {
	// processing ID
	ProcessingID String `position:"PT.1"`
	// processing mode
	ProcessingMode String `position:"PT.2"`
}
