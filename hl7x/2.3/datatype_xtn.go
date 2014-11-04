package hl7v2_3

type Xtn struct {
	// [(999)] 999-9999 [X99999][C any text]
	A_999_999_9999X99999CAnyText String `position:"XTN.1"`
	// telecommunication use code
	TelecommunicationUseCode String `position:"XTN.2"`
	// telecommunication equipment type (ID)
	TelecommunicationEquipmentTypeId String `position:"XTN.3"`
	// Email address
	EmailAddress String `position:"XTN.4"`
	// Country Code
	CountryCode String `position:"XTN.5"`
	// Area/city code
	AreaCityCode String `position:"XTN.6"`
	// Phone number
	PhoneNumber String `position:"XTN.7"`
	// Extension
	Extension String `position:"XTN.8"`
	// any text
	AnyText String `position:"XTN.9"`
}
