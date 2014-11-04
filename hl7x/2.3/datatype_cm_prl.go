package hl7v2_3

type CmPrl struct {
	// OBX-3 observation identifier of parent result
	Obx_3ObservationIdentifierOfParentResult Ce `position:"CM_PRL.1"`
	// OBX-4 sub-ID of parent result
	Obx_4SubIdOfParentResult String `position:"CM_PRL.2"`
	// part of OBX-5 observation result from parent
	PartOfObx_5ObservationResultFromParent Tx `position:"CM_PRL.3"`
}
