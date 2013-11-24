package hl7

import (
	"bytes"
	"errors"
	//"regexp"
	"fmt"
)

// unmarshals passed byte data
func Unmarshal(data []byte) (values Values, err error) {
	fmt.Println(string(data))
	if !bytes.HasPrefix(data, []byte("MSH")) {
		return Values{}, errors.New("Could not find MSH header")
	}

	segments := bytes.Split(data, []byte{0x0D})

	if len(segments) < 2 {
		return Values{}, errors.New("Not enough segments in hl7 data")
	}

	hdr := Header{}
	hdr.CompositeDelimiter = data[3]
	hdr.SubCompositeDelimiter = data[4]
	hdr.SubSubCompositeDelimiter = data[5]
	hdr.EscapeCharacter = data[6]
	hdr.RepetitionDelimiter = data[7]
	hdr.Values, err = unmarshalSegment(data[9:])
	if err != nil {
		return Values{}, err
	}

	values = append(values, hdr)

	for _, segment := range segments {
		svalues, err := unmarshalSegment(segment)
		if err != nil {
			return Values{}, err
		}
		values = append(values, svalues)
	}

	return values, nil
}

func unmarshalSegment(data []byte) (values Values, err error) {
	return Values{data}, nil
}
