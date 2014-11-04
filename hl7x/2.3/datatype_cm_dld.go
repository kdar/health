package hl7v2_3

type CmDld struct {
	// discharge location
	DischargeLocation String `position:"CM_DLD.1"`
	// effective date
	EffectiveDate String `position:"CM_DLD.2"`
}
