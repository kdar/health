package hl7

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	c_SEGMENT_SEPARATOR    = []byte{0x0d}
	c_ESC_ESCAPE_CHAR      = []byte(`\E\`)
	c_ESC_FIELD_SEP        = []byte(`\F\`)
	c_ESC_REPETITION_SEP   = []byte(`\R\`)
	c_ESC_COMPONENT_SEP    = []byte(`\S\`)
	c_ESC_SUBCOMPONENT_SEP = []byte(`\T\`)
)

// Marshal converts a slice of Segment structs to a byte array containing an HL7
// message in "pipehat" format.
func Marshal(segments []Segment) ([]byte, error) {
	if len(segments) == 0 {
		return nil, errors.New("no data to marshal")
	}

	buf := bytes.Buffer{}
	encoder := NewEncoder(&buf)
	if err := encoder.Encode(segments); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Encoder struct {
	Writer          io.Writer
	segmentSep      []byte
	charsToEscape   string
	fieldSep        []byte
	componentSep    []byte
	repetitionSep   []byte
	escapeChar      []byte
	subcomponentSep []byte
}

// NewEncoder creates and initializes an Encoder which will write an encoded HL7
// message to the provided io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Writer:     w,
		segmentSep: c_SEGMENT_SEPARATOR,
	}
}

// Encode converts the provided slice of Segments stored in the Encoder to an
// HL7 "pipehat" encoded message and writes it to the io.Writer set in its
// Writer field.  The Segment slice must contain a valid message header (MSH) segment
// with fields 1 and 2 populated with the field separator and escape characters;
// if this segment is not present, Encode will return an error.
func (e *Encoder) Encode(segments []Segment) (err error) {
	// TODO find MSH
	var msh_segment *Segment

	// find MSH segment and get delimiters
	for _, s := range segments {
		name, _ := e.fieldDataAt(s, 0)
		if string(name) == "MSH" {
			msh_segment = &s
			break
		}
	}

	if msh_segment == nil {
		return errors.New("Missing required MSH segment")
	}

	if err = e.extractSeparators(*msh_segment); err != nil {
		return err
	}

	for _, s := range segments {
		enc, err := e.encodeSegment(s)
		if err != nil {
			return err
		}
		// each segment must end in segmentSep, which is slightly different than
		// other delimiters which do not appear at the end of a series (e.g.
		// component1^component2).
		e.Writer.Write(enc)
		e.Writer.Write(e.segmentSep)
	}

	return nil
}

func (e *Encoder) extractSeparators(s Segment) error {
	fs, _ := e.fieldDataAt(s, 1)
	if len(fs) >= 1 {
		e.fieldSep = []byte{fs[0]}
	} else {
		return errors.New("missing MSH-1 Field Separator")
	}

	delims, _ := e.fieldDataAt(s, 2)
	if len(delims) >= 4 {
		e.componentSep = []byte{delims[0]}
		e.repetitionSep = []byte{delims[1]}
		e.escapeChar = []byte{delims[2]}
		e.subcomponentSep = []byte{delims[3]}
	} else {
		return errors.New("Missing or truncated MSH-2 Encoding Characters")
	}

	// need to escape everything in the separator chars field plus the
	// field separator character.
	e.charsToEscape = fmt.Sprintf("%s%s", string(delims), string(fs))

	return nil

}

func (e *Encoder) fieldDataAt(segment Segment, index int) ([]byte, error) {
	data, _ := segment.Index(index)
	switch v := data.(type) {
	case Field:
		return data.(Field), nil
	default:
		return nil, fmt.Errorf("cannot get field data for type %T", v)
	}

	return nil, errors.New("unexpected field data decoding error")
}

func (e *Encoder) encodeData(data Data) (buf []byte, err error) {
	switch v := data.(type) {

	case Field:
		// Field is a byte array
		buf = data.(Field)

		// escape any separator characters in the field data, if present
		// http://www.hl7standards.com/blog/2006/11/02/hl7-escape-sequences/

		if bytes.IndexAny(buf, e.charsToEscape) > -1 {
			buf = bytes.Replace(buf, e.escapeChar, c_ESC_ESCAPE_CHAR, -1)
			buf = bytes.Replace(buf, e.fieldSep, c_ESC_FIELD_SEP, -1)
			buf = bytes.Replace(buf, e.componentSep, c_ESC_COMPONENT_SEP, -1)
			buf = bytes.Replace(buf, e.repetitionSep, c_ESC_REPETITION_SEP, -1)
			buf = bytes.Replace(buf, e.subcomponentSep, c_ESC_SUBCOMPONENT_SEP, -1)
		}

	case SubComponent:
		// SubComponent is array of fields delimited by subcomponentSep
		var b [][]byte
		subcomponent := data.(SubComponent)
		l := subcomponent.Len()
		for i := 0; i < l; i++ {
			d, _ := subcomponent.Index(i)
			enc, err := e.encodeData(d)
			if err != nil {
				return nil, err
			}
			b = append(b, enc)
		}
		buf = bytes.Join(b, e.subcomponentSep)
		break

	case Component:
		// Component contains fields and subcomponents
		var b [][]byte
		component := data.(Component)
		l := component.Len()
		for i := 0; i < l; i++ {
			d, _ := component.Index(i)
			enc, err := e.encodeData(d)
			if err != nil {
				return nil, err
			}
			b = append(b, enc)
		}
		buf = bytes.Join(b, e.componentSep)
		break

	case Repeated:
		// Repeated contains fields, subcomponents, and components
		var b [][]byte
		repeated := data.(Repeated)
		l := repeated.Len()
		for i := 0; i < l; i++ {
			d, _ := repeated.Index(i)
			enc, err := e.encodeData(d)
			if err != nil {
				return nil, err
			}
			b = append(b, enc)
		}
		buf = bytes.Join(b, e.repetitionSep)
		break

	default:
		return nil, fmt.Errorf("unrecognized type: %T", v)
	}

	return
}

func (e *Encoder) encodeSegment(segment Segment) (buf []byte, err error) {
	l := segment.Len()

	// TODO add first field
	var b [][]byte

	for i := 0; i < l; i++ {
		data, _ := segment.Index(i)
		enc, err := e.encodeData(data)
		if err != nil {
			return nil, err
		}

		if string(enc) == "MSH" {
			// add the segment name
			b = append(b, enc)

			// next two fields are the separators
			// first, skip the field separator since we parsed it out in NewEncoder
			i++

			// next, get the separator characters without running through encode(),
			// which will escape them.
			i++
			seps, _ := e.fieldDataAt(segment, i)
			b = append(b, seps)
		} else {
			b = append(b, enc)
		}
	}

	return bytes.Join(b, e.fieldSep), nil
}
