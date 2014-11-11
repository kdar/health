package hl7v2_3

import (
	"strings"
)

// it really should be []interface{}, but for now
// just doing string until I run into something
// that breaks
type Varies []string

func (v Varies) String() string {
	return strings.Join(v, "")
}

// Is
// Id
// St
// Si
// Nm
// CmOsd
// Tn
// Dt
type String string

func (s String) String() string {
	return string(s)
}
