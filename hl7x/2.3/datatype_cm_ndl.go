package hl7v2_3

type CmNdl struct {
	// name
	Name Cn `position:"CM_NDL.1"`
	// start date/time
	StartDateTime Ts `position:"CM_NDL.2"`
	// end date/time
	EndDateTime Ts `position:"CM_NDL.3"`
	// point of care (IS)
	PointOfCareIs String `position:"CM_NDL.4"`
	// room
	Room String `position:"CM_NDL.5"`
	// bed
	Bed String `position:"CM_NDL.6"`
	// facility (HD)
	FacilityHd Hd `position:"CM_NDL.7"`
	// location status
	LocationStatus String `position:"CM_NDL.8"`
	// person location type
	PersonLocationType String `position:"CM_NDL.9"`
	// building
	Building String `position:"CM_NDL.10"`
	// floor
	Floor String `position:"CM_NDL.11"`
}
