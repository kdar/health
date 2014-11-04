package hl7v2_3

type Tq struct {
	// quantity
	Quantity Cq `position:"TQ.1"`
	// interval
	Interval CmRi `position:"TQ.2"`
	// duration
	Duration String `position:"TQ.3"`
	// start date/time
	StartDateTime Ts `position:"TQ.4"`
	// end date/time
	EndDateTime Ts `position:"TQ.5"`
	// priority
	Priority String `position:"TQ.6"`
	// condition
	Condition String `position:"TQ.7"`
	// text (TX)
	TextTx Tx `position:"TQ.8"`
	// conjunction
	Conjunction String `position:"TQ.9"`
	// order sequencing
	OrderSequencing String `position:"TQ.10"`
}
