package hl7v2_3

type CmEip struct {
	// parent´s placer order number
	ParentsPlacerOrderNumber Ei `position:"CM_EIP.1"`
	// parent´s filler order number
	ParentsFillerOrderNumber Ei `position:"CM_EIP.2"`
}
