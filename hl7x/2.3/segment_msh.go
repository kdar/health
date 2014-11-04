package hl7v2_3

type Msh struct {
	// Field Separator
	FieldSeparator String `position:"MSH.1" require:"true"`
	// Encoding Characters
	EncodingCharacters String `position:"MSH.2" require:"true"`
	// Sending Application
	SendingApplication Hd `position:"MSH.3"`
	// Sending Facility
	SendingFacility Hd `position:"MSH.4"`
	// Receiving Application
	ReceivingApplication Hd `position:"MSH.5"`
	// Receiving Facility
	ReceivingFacility Hd `position:"MSH.6"`
	// Date / Time of Message
	DateTimeOfMessage Ts `position:"MSH.7"`
	// Security
	Security String `position:"MSH.8"`
	// Message Type
	MessageType CmMsg `position:"MSH.9" require:"true"`
	// Message Control ID
	MessageControlID String `position:"MSH.10" require:"true"`
	// Processing ID
	ProcessingID Pt `position:"MSH.11" require:"true"`
	// Version ID
	VersionID String `position:"MSH.12" require:"true"`
	// Sequence Number
	SequenceNumber String `position:"MSH.13"`
	// Continuation Pointer
	ContinuationPointer String `position:"MSH.14"`
	// Accept Acknowledgement Type
	AcceptAcknowledgementType String `position:"MSH.15"`
	// Application Acknowledgement Type
	ApplicationAcknowledgementType String `position:"MSH.16"`
	// Country Code
	CountryCode String `position:"MSH.17"`
	// Character Set
	CharacterSet String `position:"MSH.18"`
	// Principal Language of Message
	PrincipalLanguageOfMessage Ce `position:"MSH.19"`
}
