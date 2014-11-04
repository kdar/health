package hl7v2_3

type Pl struct {
	// point of care (ID)
	PointOfCareID String `position:"PL.1"`
	// room
	Room String `position:"PL.2"`
	// bed
	Bed String `position:"PL.3"`
	// facility (HD)
	FacilityHd Hd `position:"PL.4"`
	// location status
	LocationStatus String `position:"PL.5"`
	// person location type
	PersonLocationType String `position:"PL.6"`
	// building
	Building String `position:"PL.7"`
	// floor
	Floor String `position:"PL.8"`
	// Location type
	LocationType String `position:"PL.9"`
}
