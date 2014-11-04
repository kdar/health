package hl7v2_3

type CmRi struct {
	// repeat pattern
	RepeatPattern String `position:"CM_RI.1"`
	// explicit time interval
	ExplicitTimeInterval String `position:"CM_RI.2"`
}
