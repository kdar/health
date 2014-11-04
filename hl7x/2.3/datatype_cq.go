package hl7v2_3

type Cq struct {
	// quantity
	Quantity String `position:"CQ.1"`
	// units
	Units Ce `position:"CQ.2"`
}
