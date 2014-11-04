package hl7v2_3

type CmMsg struct {
	// message type
	MessageType String `position:"CM_MSG.1"`
	// trigger event
	TriggerEvent String `position:"CM_MSG.2"`
}
