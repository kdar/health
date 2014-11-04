package hl7v2_3

type Orc struct {
	// Order Control
	OrderControl String `position:"ORC.1" require:"true"`
	// Placer Order Number
	PlacerOrderNumbers []Ei `position:"ORC.2"`
	// Filler Order Number
	FillerOrderNumber Ei `position:"ORC.3"`
	// Placer Group Number
	PlacerGroupNumber Ei `position:"ORC.4"`
	// Order Status
	OrderStatus String `position:"ORC.5"`
	// Response Flag
	ResponseFlag String `position:"ORC.6"`
	// Quantity/Timing
	QuantityTiming Tq `position:"ORC.7" require:"true"`
	// Parent
	Parent CmEip `position:"ORC.8"`
	// Date/Time of Transaction
	DateTimeOfTransaction Ts `position:"ORC.9"`
	// Entered By
	EnteredBy Xcn `position:"ORC.10"`
	// Verified By
	VerifiedBy Xcn `position:"ORC.11"`
	// Ordering Provider
	OrderingProviders []Xcn `position:"ORC.12"`
	// Enterer's Location
	EntererSLocation Pl `position:"ORC.13"`
	// Call Back Phone Number
	CallBackPhoneNumbers []String `position:"ORC.14"`
	// Order Effective Date/Time
	OrderEffectiveDateTime Ts `position:"ORC.15"`
	// Order Control Code Reason
	OrderControlCodeReason Ce `position:"ORC.16"`
	// Entering Organization
	EnteringOrganization Ce `position:"ORC.17"`
	// Entering Device
	EnteringDevice Ce `position:"ORC.18"`
	// Action By
	ActionBy Xcn `position:"ORC.19"`
}
