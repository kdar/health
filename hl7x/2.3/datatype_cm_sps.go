package hl7v2_3

type CmSps struct {
	// specimen source name or code
	SpecimenSourceNameOrCode Ce `position:"CM_SPS.1"`
	// additives
	Additives Tx `position:"CM_SPS.2"`
	// freetext
	Freetext Tx `position:"CM_SPS.3"`
	// body site
	BodySite Ce `position:"CM_SPS.4"`
	// site modifier
	SiteModifier Ce `position:"CM_SPS.5"`
	// collection modifier method code
	CollectionModifierMethodCode Ce `position:"CM_SPS.6"`
}
