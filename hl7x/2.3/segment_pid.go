package hl7v2_3

type Pid struct {
	// Set String - Patient ID
	SetIDPatientID String `position:"PID.1"`
	// Patient String (External ID)
	PatientIDExternalID Cx `position:"PID.2"`
	// Patient String (Internal ID)
	PatientIDInternalIDs []Cx `position:"PID.3" require:"true"`
	// Alternate Patient ID
	AlternatePatientID Cx `position:"PID.4"`
	// Patient Name
	PatientName Xpn `position:"PID.5" require:"true"`
	// Mother's Maiden Name
	MotherSMaidenName Xpn `position:"PID.6"`
	// Date of Birth
	DateOfBirth Ts `position:"PID.7"`
	// Sex
	Sex String `position:"PID.8"`
	// Patient Alias
	PatientAliases []Xpn `position:"PID.9"`
	// Race
	Race String `position:"PID.10"`
	// Patient Address
	PatientAddresses []Xad `position:"PID.11"`
	// County Code
	CountyCode String `position:"PID.12"`
	// Phone Number - Home
	PhoneNumberHomes []Xtn `position:"PID.13"`
	// Phone Number - Business
	PhoneNumberBusinesses []Xtn `position:"PID.14"`
	// Primary Language
	PrimaryLanguage Ce `position:"PID.15"`
	// Marital Status
	MaritalStatuses []String `position:"PID.16"`
	// Religion
	Religion String `position:"PID.17"`
	// Patient Account Number
	PatientAccountNumber Cx `position:"PID.18"`
	// SSN Number - Patient
	SsnNumberPatient String `position:"PID.19"`
	// Driver's License Number
	DriverSLicenseNumber Dln `position:"PID.20"`
	// Mother's Identifier
	MotherSIdentifier Cx `position:"PID.21"`
	// Ethnic Group
	EthnicGroup String `position:"PID.22"`
	// Birth Place
	BirthPlace String `position:"PID.23"`
	// Multiple Birth Indicator
	MultipleBirthIndicator String `position:"PID.24"`
	// Birth Order
	BirthOrder String `position:"PID.25"`
	// Citizenship
	Citizenship String `position:"PID.26"`
	// Veterans Military Status
	VeteransMilitaryStatus Ce `position:"PID.27"`
	// Nationality Code
	NationalityCode Ce `position:"PID.28"`
	// Patient Death Date and Time
	PatientDeathDateAndTime Ts `position:"PID.29"`
	// Patient Death Indicator
	PatientDeathIndicator String `position:"PID.30"`
}
