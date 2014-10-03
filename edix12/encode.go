package edix12

// import (
// 	"bytes"

// 	"github.com/kdar/health/edix12/node"
// )

// const (
// 	isaTpl = []string{
// 		`ISA`,
// 		`00`,              // Authorization Information Qualifier. M. ID. 2/2
// 		`No Auth   `,      // Authorization Information. M. AN. 10/10
// 		`00`,              // Security Information Qualifier. M. ID. 2/2
// 		`No Securit`,      // Security Information. M. AN. 10/10
// 		`ZZ`,              // Interchange ID Qualifier. M. ID. 2/2
// 		`SUBMITTERS.ID  `, // Interchange Sender ID. M. AN. 15/15
// 		`ZZ`,              // Interchange ID Qualifier. M. ID. 2/2
// 		`RECEIVERS.ID   `, // Interchange Receiver ID. M. AN. 15/15
// 		`YYMMDD`,          // Interchange Date. M. DT. 6/6
// 		`HHMM`,            // Interchange Time. M. TM. 4/4
// 		`^`,               // Repetition Separator. M. 1/1
// 		`00501`,           // Interchange Control Version Number. M. ID. 5/5
// 		``,                // Interchange Control Number. M. Nn. 9/9
// 		`0`,               // Acknowledgment Requested. M. ID. 1/1
// 		`T`,               // Usage Indicator. M. ID. 1/1
// 		`:~`,              // Component Element Separator + Segment Terminator. M. 2/2
// 	}
// )

// func Marshal(v interface{}) ([]byte, error) {
// 	e := encodeState{}
// 	err := e.marshal(v)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return e.Bytes(), nil
// }

// type encodeState struct {
// 	elementSeparator    byte
// 	componentSeparator  byte
// 	segmentTerminator   byte
// 	repetitionSeparator byte
// 	bytes.Buffer
// }

// func (e *encodeState) marshal(v interface{}) (err error) {
// 	switch t := v.(type) {
// 	case *node.Segment:
// 		err = e.marshalSegment(t)
// 	case node.Element:
// 	case node.Component:
// 	}

// 	return
// }

// func (e *encodeState) marshalSegment(n *node.Segment) error {
// 	switch n.Tag {
// 	case "ISA":
// 		e.marshalISA(n)
// 	case "GS":
// 	case "ST":
// 	}
// 	return nil
// }

// func (e *encodeState) marshalISA(n *node.Segment) error {
// 	e.elementSeparator
// 	e.componentSeparator
// 	e.repetitionSeparator
// 	e.segmentTerminator
// }
