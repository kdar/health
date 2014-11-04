package hl7v2_3

type Ei struct {
	// entity identifier
	EntityIdentifier String `position:"EI.1"`
	// namespace ID
	NamespaceID String `position:"EI.2"`
	// universal ID
	UniversalID String `position:"EI.3"`
	// universal String type
	UniversalIdType String `position:"EI.4"`
}
