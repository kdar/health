package edifact

import (
	"bytes"
	"errors"
	"fmt"
	"unicode/utf8"
)

// Unmarshals a byte slice like this: UNA~|.?^'
func unmarshalUNA(data []byte) (int, *EDIFACT, error) {
	if len(data) < 3 {
		return 0, nil, errors.New("Length of data is too small")
	}

	// Return the default UNA delimiters
	if !bytes.Equal(data[:3], []byte(UNA_SEGMENT_CODE)) {
		return 0, NewDefault(), nil
	}

	return 9, &EDIFACT{
		ComponentDelimiter: data[3],
		DataDelimiter:      data[4],
		Decimal:            data[5],
		ReleaseIndicator:   data[6],
		ReservedDelimiter:  data[7],
		SegmentTerminator:  data[8],
	}, nil
}

func Unmarshal(data []byte) ([]Segment, error) {
	var segments []Segment

	n, ef, err := unmarshalUNA(data)
	if err != nil {
		return nil, err
	}

	// Skip UNA header if there is one
	data = data[n:]

	if data[len(data)-1] != ef.SegmentTerminator {
		return nil, errors.New(fmt.Sprintf("Data is not terminated by the segment terminator: %c", ef.SegmentTerminator))
	}

	byteSegments := splitWithEscape(data, []byte{ef.SegmentTerminator}, []byte{ef.ReleaseIndicator})

	// Because of the way split works, and a segment is terminated
	// and not delimited, we will have a blank segment at the end.
	// Remove it.
	byteSegments = byteSegments[:len(byteSegments)-1]

	// By the spec, we 4 segments are mandatory, but we'll be
	// forgiving and test for only one, since we can still parse
	// it.
	if len(byteSegments) < 1 {
		return nil, errors.New("Data does not contain any segments")
	}

	for _, byteSegment := range byteSegments {
		var segment Segment
		byteDataSlice := splitWithEscape(byteSegment, []byte{ef.DataDelimiter}, []byte{ef.ReleaseIndicator})

		for _, byteData := range byteDataSlice {
			//var components []string
			//fmt.Printf("%q\n", byteData)

			delimiter := ef.ComponentDelimiter

			// rn := bytes.IndexByte(byteData, ef.ReservedDelimiter)
			// if rn != -1 && rn > 0 && byteData[rn-1] != ef.ReleaseIndicator {
			//   delimiter = ef.ReservedDelimiter
			// }

			components := splitWithEscape(byteData, []byte{delimiter}, []byte{ef.ReleaseIndicator})
			if len(components) == 0 {
				segment = append(segment, "")
			} else if len(components) == 1 {
				segment = append(segment, string(removeEscape(components[0], []byte{ef.ReleaseIndicator})))
			}

			//fmt.Printf("%q\n", components)
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

// Removes an escape character. If the escape character
// is itself escaped, then one is left.
func removeEscape(s, escape []byte) []byte {
	if len(s) == 0 {
		return []byte("")
	}

	ec := bytes.Count(s, escape)
	ec = ec - bytes.Count(s, append(escape, escape...))
	s2 := make([]byte, len(s)-(len(escape)*ec))

	for i, i2 := 0, 0; i < len(s) && i2 < len(s2); i, i2 = i+1, i2+1 {
		if s[i] == escape[0] && (len(escape) == 1 || bytes.Equal(s[i:i+len(escape)], escape)) {
			if len(s) > (i+2*len(escape)) && bytes.Equal(s[i+len(escape):i+2*len(escape)], escape) {
				s2[i] = s[i]
				continue
			}
			i2 -= 1
		} else {
			s2[i2] = s[i]
		}
	}

	return s2
}

// explode splits s into an array of UTF-8 sequences, one per Unicode character (still arrays of bytes),
// up to a maximum of n byte arrays. Invalid UTF-8 sequences are chopped into individual bytes.
func explode(s []byte, n int) [][]byte {
	if n <= 0 {
		n = len(s)
	}
	a := make([][]byte, n)
	var size int
	na := 0
	for len(s) > 0 {
		if na+1 >= n {
			a[na] = s
			na++
			break
		}
		_, size = utf8.DecodeRune(s)
		a[na] = s[0:size]
		s = s[size:]
		na++
	}
	return a[0:na]
}

// Splits a []byte by the given seperator, and honoring
// the escape []byte. If an escape comes right before the
// separator, then it is not split at that junction.
// You can use escape to escape itself, in which case if an
// escape is escaped right before a separator, then it will
// be split at that junction.
func splitWithEscape(s, sep []byte, escape []byte) [][]byte {
	sepSave := 0
	n := -1

	if n == 0 {
		return nil
	}
	if len(sep) == 0 {
		return explode(s, n)
	}
	if n < 0 {
		n = bytes.Count(s, sep) + 1
		// subtract separators that are escaped
		n = n - bytes.Count(s, append(escape, sep...)) + 1
	}
	c := sep[0]
	start := 0
	a := make([][]byte, n)
	na := 0

	escapeCount := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		foundEscape := i+len(escape) < len(s) && bytes.Equal(s[i:i+len(escape)], escape)
		if foundEscape {
			escapeCount++
		}

		if s[i] == c && (len(sep) == 1 || bytes.Equal(s[i:i+len(sep)], sep)) && escapeCount == 0 {
			a[na] = s[start : i+sepSave]
			na++
			start = i + len(sep)
			i += len(sep) - 1
		}

		// Reset if we didn't find an escape.
		// Reset if the escapeCount is 2, because this means
		// the escape escaped itself.
		if !foundEscape || escapeCount == 2 {
			escapeCount = 0
		}
	}
	a[na] = s[start:]
	return a[0 : na+1]
}
