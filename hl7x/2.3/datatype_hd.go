package hl7v2_3

type Hd struct {
	// namespace ID
	NamespaceID String `position:"HD.1"`
	// universal ID
	UniversalID String `position:"HD.2"`
	// universal String type
	UniversalIDType String `position:"HD.3"`
}
