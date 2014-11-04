package hl7v2_3

type Obr struct {
	// Set String - Observation Request
	SetIDObservationRequest String `position:"OBR.1"`
	// Placer Order Number
	PlacerOrderNumbers []Ei `position:"OBR.2"`
	// Filler Order Number
	FillerOrderNumber Ei `position:"OBR.3"`
	// Universal Service Identifier
	UniversalServiceIdentifier Ce `position:"OBR.4" require:"true"`
	// Priority
	Priority String `position:"OBR.5"`
	// Requested Date/Time
	RequestedDateTime Ts `position:"OBR.6"`
	// Observation Date/Time
	ObservationDateTime Ts `position:"OBR.7"`
	// Observation End Date/Time
	ObservationEndDateTime Ts `position:"OBR.8"`
	// Collection Volume
	CollectionVolume Cq `position:"OBR.9"`
	// Collector Identifier
	CollectorIdentifiers []Xcn `position:"OBR.10"`
	// Specimen Action Code
	SpecimenActionCode String `position:"OBR.11"`
	// Danger Code
	DangerCode Ce `position:"OBR.12"`
	// Relevant Clinical Information
	RelevantClinicalInformation String `position:"OBR.13"`
	// Specimen Received Date/Time
	SpecimenReceivedDateTime Ts `position:"OBR.14"`
	// Specimen Source
	SpecimenSource CmSps `position:"OBR.15"`
	// Ordering Provider
	OrderingProviders []Xcn `position:"OBR.16"`
	// Order Callback Phone Number
	OrderCallbackPhoneNumbers []Xtn `position:"OBR.17"`
	// Placer Field 1
	PlacerField1 String `position:"OBR.18"`
	// Placer Field 2
	PlacerField2 String `position:"OBR.19"`
	// Filler Field 1
	FillerField1 String `position:"OBR.20"`
	// Filler Field 2
	FillerField2 String `position:"OBR.21"`
	// Results Rpt/Status Chng - Date/Time
	ResultsRptStatusChngDateTime Ts `position:"OBR.22"`
	// Charge To Practice
	ChargeToPractice CmMoc `position:"OBR.23"`
	// Diagnostic Service Section ID
	DiagnosticServiceSectionId String `position:"OBR.24"`
	// Result Status
	ResultStatus String `position:"OBR.25"`
	// Parent Result
	ParentResult CmPrl `position:"OBR.26"`
	// Quantity/Timing
	QuantityTiming Tq `position:"OBR.27" require:"true"`
	// Result Copies To
	ResultCopiesTos []Xcn `position:"OBR.28"`
	// Parent Number
	ParentNumber CmEip `position:"OBR.29"`
	// Transportation Mode
	TransportationMode String `position:"OBR.30"`
	// Reason For Study
	ReasonForStudies []Ce `position:"OBR.31"`
	// Principal Result Interpreter
	PrincipalResultInterpreter CmNdl `position:"OBR.32"`
	// Assistant Result Interpreter
	AssistantResultInterpreters []CmNdl `position:"OBR.33"`
	// Technician
	Technicians []CmNdl `position:"OBR.34"`
	// Transcriptionist
	Transcriptionists []CmNdl `position:"OBR.35"`
	// Scheduled Date/Time
	ScheduledDateTime Ts `position:"OBR.36"`
	// Number Of Sample Containers
	NumberOfSampleContainers String `position:"OBR.37"`
	// Transport Logistics Of Collected Sample
	TransportLogisticsOfCollectedSamples []Ce `position:"OBR.38"`
	// Collector’s Comment
	CollectorSComments []Ce `position:"OBR.39"`
	// Transport Arrangement Responsibility
	TransportArrangementResponsibility Ce `position:"OBR.40"`
	// Transport Arranged
	TransportArranged String `position:"OBR.41"`
	// Escort Required
	EscortRequired String `position:"OBR.42"`
	// Planned Patient Transport Comment
	PlannedPatientTransportComments []Ce `position:"OBR.43"`
}
