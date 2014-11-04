package hl7x

import (
	"fmt"
	"strings"
)

type Error struct {
	Errors []string
}

func (e *Error) Error() string {
	points := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d error(s) decoding:\n\n%s",
		len(e.Errors), strings.Join(points, "\n"))
}

func (e *Error) append(err error) {
	switch v := err.(type) {
	case *Error:
		e.Errors = append(e.Errors, v.Errors...)
	default:
		e.Errors = append(e.Errors, v.Error())
	}
}
