package hl7v2_3

type Pv1 struct {
	// Set String - Patient Visit
	SetIDPatientVisit String `position:"PV1.1"`
	// Patient Class
	PatientClass String `position:"PV1.2" require:"true"`
	// Assigned Patient Location
	AssignedPatientLocation Pl `position:"PV1.3"`
	// Admission Type
	AdmissionType String `position:"PV1.4"`
	// Preadmit Number
	PreadmitNumber Cx `position:"PV1.5"`
	// Prior Patient Location
	PriorPatientLocation Pl `position:"PV1.6"`
	// Attending Doctor
	AttendingDoctor Xcn `position:"PV1.7"`
	// Referring Doctor
	ReferringDoctor Xcn `position:"PV1.8"`
	// Consulting Doctor
	ConsultingDoctors []Xcn `position:"PV1.9"`
	// Hospital Service
	HospitalService String `position:"PV1.10"`
	// Temporary Location
	TemporaryLocation Pl `position:"PV1.11"`
	// Preadmit Test Indicator
	PreadmitTestIndicator String `position:"PV1.12"`
	// Readmission Indicator
	ReadmissionIndicator String `position:"PV1.13"`
	// Admit Source
	AdmitSource String `position:"PV1.14"`
	// Ambulatory Status
	AmbulatoryStatus String `position:"PV1.15"`
	// VIP Indicator
	VipIndicator String `position:"PV1.16"`
	// Admitting Doctor
	AdmittingDoctor Xcn `position:"PV1.17"`
	// Patient Type
	PatientType String `position:"PV1.18"`
	// Visit Number
	VisitNumber Cx `position:"PV1.19"`
	// Financial Class
	FinancialClasses []Fc `position:"PV1.20"`
	// Charge Price Indicator
	ChargePriceIndicator String `position:"PV1.21"`
	// Courtesy Code
	CourtesyCode String `position:"PV1.22"`
	// Credit Rating
	CreditRating String `position:"PV1.23"`
	// Contract Code
	ContractCodes []String `position:"PV1.24"`
	// Contract Effective Date
	ContractEffectiveDates []String `position:"PV1.25"`
	// Contract Amount
	ContractAmounts []String `position:"PV1.26"`
	// Contract Period
	ContractPeriods []String `position:"PV1.27"`
	// Interest Code
	InterestCode String `position:"PV1.28"`
	// Transfer to Bad Debt Code
	TransferToBadDebtCode String `position:"PV1.29"`
	// Transfer to Bad Debt Date
	TransferToBadDebtDate String `position:"PV1.30"`
	// Bad Debt Agency Code
	BadDebtAgencyCode String `position:"PV1.31"`
	// Bad Debt Transfer Amount
	BadDebtTransferAmount String `position:"PV1.32"`
	// Bad Debt Recovery Amount
	BadDebtRecoveryAmount String `position:"PV1.33"`
	// Delete Account Indicator
	DeleteAccountIndicator String `position:"PV1.34"`
	// Delete Account Date
	DeleteAccountDate String `position:"PV1.35"`
	// Discharge Disposition
	DischargeDisposition String `position:"PV1.36"`
	// Discharged to Location
	DischargedToLocation CmDld `position:"PV1.37"`
	// Diet Type
	DietType String `position:"PV1.38"`
	// Servicing Facility
	ServicingFacility String `position:"PV1.39"`
	// Bed Status
	BedStatus String `position:"PV1.40"`
	// Account Status
	AccountStatus String `position:"PV1.41"`
	// Pending Location
	PendingLocation Pl `position:"PV1.42"`
	// Prior Temporary Location
	PriorTemporaryLocation Pl `position:"PV1.43"`
	// Admit Date/Time
	AdmitDateTime Ts `position:"PV1.44"`
	// Discharge Date/Time
	DischargeDateTime Ts `position:"PV1.45"`
	// Current Patient Balance
	CurrentPatientBalance String `position:"PV1.46"`
	// Total Charges
	TotalCharges String `position:"PV1.47"`
	// Total Adjustments
	TotalAdjustments String `position:"PV1.48"`
	// Total Payments
	TotalPayments String `position:"PV1.49"`
	// Alternate Visit ID
	AlternateVisitID Cx `position:"PV1.50"`
	// Visit Indicator
	VisitIndicator String `position:"PV1.51"`
	// Other Healthcare Provider
	OtherHealthcareProviders []Xcn `position:"PV1.52"`
}
