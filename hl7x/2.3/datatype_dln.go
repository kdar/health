package hl7v2_3

type Dln struct {
	// Driver´s License Number
	DriversLicenseNumber String `position:"DLN.1"`
	// Issuing State province country
	IssuingStateProvinceCountry String `position:"DLN.2"`
	// expiration date
	ExpirationDate String `position:"DLN.3"`
}
