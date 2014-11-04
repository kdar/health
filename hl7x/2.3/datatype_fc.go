package hl7v2_3

type Fc struct {
	// Financial Class
	FinancialClass String `position:"FC.1"`
	// Effective Date
	EffectiveDate Ts `position:"FC.2"`
}
