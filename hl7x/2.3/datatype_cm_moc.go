package hl7v2_3

type CmMoc struct {
	// dollar amount
	DollarAmount Mo `position:"CM_MOC.1"`
	// charge code
	ChargeCode Ce `position:"CM_MOC.2"`
}
