package hl7v2_3

type Obx struct {
	// Set String - OBX
	SetIDOBX String `position:"OBX.1"`
	// Value Type
	ValueType String `position:"OBX.2" require:"true"`
	// Observation Identifier
	ObservationIdentifier Ce `position:"OBX.3" require:"true"`
	// Observation Sub-ID
	ObservationSubID String `position:"OBX.4"`
	// Observation Value
	ObservationValues Varies `position:"OBX.5"`
	// Units
	Units Ce `position:"OBX.6"`
	// References Range
	ReferencesRange String `position:"OBX.7"`
	// Abnormal Flags
	AbnormalFlags []String `position:"OBX.8"`
	// Probability
	Probability String `position:"OBX.9"`
	// Nature of Abnormal Test
	NatureOfAbnormalTest String `position:"OBX.10"`
	// Observ Result Status
	ObservResultStatus String `position:"OBX.11" require:"true"`
	// Date Last Obs Normal Values
	DateLastObsNormalValues Ts `position:"OBX.12"`
	// User Defined Access Checks
	UserDefinedAccessChecks String `position:"OBX.13"`
	// Date/Time of the Observation
	DateTimeOfTheObservation Ts `position:"OBX.14"`
	// Producer's ID
	ProducerID Ce `position:"OBX.15"`
	// Responsible Observer
	ResponsibleObserver Xcn `position:"OBX.16"`
	// Observation Method
	ObservationMethods []Ce `position:"OBX.17"`
}

// func (o *Obx) UnmarshalHL7(data hl7.Data) error {
// 	fmt.Println(data)
// 	return nil
// }
