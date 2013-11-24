package ncpdp

import (
	"errors"
	//"log"
	"github.com/kdar/health/edifact"
	//"time"
)

// base struct of all messages.
// required by EDIFACT and NCPDP script
type RXH struct {
	// UIB
	// SyntaxIdentifier            string
	// SyntaxVersionNumber         string
	// TransactionControlReference string
	// SenderIds                   []string
	// ReceiverIds                 []string

	//una
	UIB edifact.Values
	UIH edifact.Values

	//RXHREQ or RXHRES

	UIT edifact.Values
	UIZ edifact.Values
}

// creates a new RXH
func newRXH() *RXH {
	return &RXH{}
}

// func (r RXH) copy(rxh *RXH) {
//   r.UIB = rxh.UIB
//   r.UIH = rxh.UIH
//   r.UIT = rxh.UIT
//   r.UIZ = rxh.UIZ
// }

// base struct for RXHRES and RXHREQ
type RXHREX struct {
	*RXH
	PVD edifact.Values
	PTT edifact.Values
	COO edifact.Values

	*Patient
}

// creates a new RXHREX
func newRXHREX() *RXHREX {
	return &RXHREX{
		RXH:     newRXH(),
		Patient: newPatient(),
	}
}

// Segment: RES
type Response struct {
	// A = Approved
	// D = Denied
	// C = Approved with changes
	// N = Denied, new prescription to
	// follow. Note: Value “N” is used in
	// REFRES transactions only.
	ResponseType string // 010 M

	// AQ= More History
	// Available. There may
	// be less then 50 drugs
	// in this response due to
	// payer processing.
	CodeListQualifier string // 020 C

	// transaction key
	ReferenceNumber string // 030 C

	// free text. could be an error
	Text string // 040 C
}

// creates a new Response
func newResponse() *Response {
	return &Response{}
}

// fills in the values for Response based on the
// parameter passed
func (r *Response) fill(values edifact.Values) {
	r.ResponseType = getString(values, 1)
	r.CodeListQualifier = getString(values, 2)
	r.ReferenceNumber = getString(values, 3)
	r.Text = getString(values, 4)
}

// Response describing a patient’s medication history. Response to RXHREQ.
type RXHRES struct {
	*RXHREX

	Response            *Response
	RequestingPhysician *Provider
	Drugs               []*Drug
}

// creates a new RXHRES
func newRXHRES() *RXHRES {
	return &RXHRES{
		RXHREX:   newRXHREX(),
		Response: newResponse(),
	}
}

// Segment: STS
type Status struct {
	StatusTypeCode    string // 010 M
	CodeListQualifier string // 020 C
	Text              string // 030 C
}

// creates a new status
func newStatus() *Status {
	return &Status{}
}

// fills in the values for Status
func (s *Status) fill(values edifact.Values) {
	s.StatusTypeCode = getString(values, 1)
	s.CodeListQualifier = getString(values, 2)
	s.Text = getString(values, 3)
}

// Message type: ERROR
type ERROR struct {
	*RXH
	Status *Status
}

// creates a new error
func newError() *ERROR {
	return &ERROR{RXH: newRXH(), Status: newStatus()}
}

// unmarshals the values based on the function.
// e.g. if function is RXHRES, values is unmarshaled
// as such.
func unmarshal(values edifact.Values, msgType string) (interface{}, error) {
	// base types
	var rxh *RXH
	// return types
	var error_ *ERROR
	var rxhres *RXHRES
	var ret interface{}

	switch msgType {
	case "RXHRES":
		rxhres = newRXHRES()
		rxh = rxhres.RXH
		ret = rxhres
	case "ERROR":
		error_ := newError()
		rxh = error_.RXH
		ret = error_
	}

	seenDrug := false

	for _, value := range values {
		if vals, ok := value.(edifact.Values); ok {
			if name, ok := vals[0].(string); ok {
				switch name {
				case "UIB":
					rxh.UIB = vals
				case "UIH":
					rxh.UIH = vals
				case "RES":
					rxhres.Response.fill(vals)
				case "PTT":
					patient := newPatient()
					patient.fill(vals)
					rxhres.RXHREX.Patient = patient
				case "COO":
					rxhres.RXHREX.COO = vals
				case "DRU":
					seenDrug = true
					drug := newDrug()
					drug.fill(vals)
					rxhres.Drugs = append(rxhres.Drugs, drug)
				case "PVD":
					provider := newProvider()
					provider.fill(vals)

					if seenDrug {
						typ, ok := vals[1].(string)
						if ok {
							if typ == "PC" { // prescriber
								rxhres.Drugs[len(rxhres.Drugs)-1].Prescriber = provider
							} else if typ == "P2" { // pharmacy
								rxhres.Drugs[len(rxhres.Drugs)-1].Pharmacy = provider
							}
						}
					} else {
						rxhres.RequestingPhysician = provider
					}
				case "STS":
					error_.Status.fill(vals)
				case "UIT":
					rxh.UIT = vals
				case "UIZ":
					rxh.UIZ = vals
				}
			}
		}
	}

	return ret, nil
}

// unmarshals NCPDP script from edifact values
func UnmarshalValues(values edifact.Values) (interface{}, error) {
	// try to determine what kind of message we're dealing
	// with and unmarshal it
	for _, value := range values {
		if vals, ok := value.(edifact.Values); ok {
			if name, ok := vals[0].(string); ok {
				switch name {
				case "UIH":
					if subvals, ok := vals[1].(edifact.Values); ok && len(subvals) > 3 {
						switch getString(subvals, 3) {
						case "RXHRES":
							return unmarshal(values, "RXHRES")
						case "ERROR":
							return unmarshal(values, "ERROR")
						}
					}
				case "RES":
					return unmarshal(values, "RXHRES")
				}
			}
		}
	}

	return nil, errors.New("Could not determine message type")
}

// unmarshals NCPDP script from a byte slice
func Unmarshal(data []byte) (interface{}, error) {
	values, err := edifact.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return UnmarshalValues(values)
}
