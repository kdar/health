package hl7v2_3

type Xad struct {
	// street address
	StreetAddress String `position:"XAD.1"`
	// other designation
	OtherDesignation String `position:"XAD.2"`
	// city
	City String `position:"XAD.3"`
	// state or province
	StateOrProvince String `position:"XAD.4"`
	// zip or postal code
	ZipOrPostalCode String `position:"XAD.5"`
	// country
	Country String `position:"XAD.6"`
	// address type
	AddressType String `position:"XAD.7"`
	// other geographic designation
	OtherGeographicDesignation String `position:"XAD.8"`
	// county/parish code
	CountyParishCode String `position:"XAD.9"`
	// census tract
	CensusTract String `position:"XAD.10"`
}
