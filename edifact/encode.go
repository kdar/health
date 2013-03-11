// Package for generating and parsing EDITFACT (http://en.wikipedia.org/wiki/EDIFACT).
// This package is pretty loose on the standard and doesn't enforce all the mandatory
// sections and such.
// Used as reference:
//   http://www.bic.org.uk/files/pdfs/070322-R2-EDIFACT-Transmission.pdf
package edifact

import (
  "bytes"
  "errors"
  "fmt"
  "reflect"
  "regexp"
)

// A callback used to process a slice while marshalling
type SliceCallback func(reflect.Value) ([]byte, error)

// make the service string advice
// func marshalUNA(e *EDIFACT) []byte {
//   return []byte(fmt.Sprintf("%s%c%c%c%c%c%c", parse.UNA_SEGMENT_NAME,
//     e.ComponentDelimiter, e.DataDelimiter,
//     e.Decimal, e.ReleaseIndicator,
//     e.RepetitionDelimiter, e.SegmentTerminator))
// }

// Marshals a part of the EDIFACT. You can pass a slice callback
// to be called if a slice is found, so you can use different
// delimiters depending on certain factors.
// I don't really like passing in delimiterRegexp, but it saves
// CPU cycles. Could possibly use cache: https://github.com/pmylund/go-cache
func marshalPart(hdr Header, data reflect.Value, delimiter byte, delimiterRegexp *regexp.Regexp, sliceCallback SliceCallback) ([]byte, error) {
  buf := &bytes.Buffer{}

  if data.Kind() == reflect.Interface {
    data = data.Elem()
  }

  switch data.Kind() {
  default:
    return []byte(""), errors.New(fmt.Sprintf("Unknown data type: %s", data.Kind()))
  case reflect.String:
    escapedData := delimiterRegexp.ReplaceAllString(data.String(), string(hdr.ReleaseIndicator())+"$1")
    buf.WriteString(escapedData)
  case reflect.Array, reflect.Slice:
    // Byte slices are special. We treat them just like the string case.
    if data.Type().Elem().Kind() == reflect.Uint8 {
      escapedData := delimiterRegexp.ReplaceAll(data.Bytes(), []byte(string(hdr.ReleaseIndicator())+"$1"))
      buf.Write(escapedData)
      break
    }

    for n := 0; n < data.Len(); n++ {
      cdata := data.Index(n)

      if sliceCallback != nil {
        cbBytes, err := sliceCallback(cdata)
        if err != nil {
          return []byte(""), err
        }
        buf.Write(cbBytes)
      }

      // we don't want to write the delimiter after the last element
      if n+1 < data.Len() {
        buf.WriteByte(delimiter)
      }
    }
  }

  return buf.Bytes(), nil
}

func Marshal(segments Values) ([]byte, error) {
  buf := &bytes.Buffer{}
  //buf.Write(marshalUNA(e))

  if len(segments) < 1 {
    return []byte(""), errors.New("Not enough segments")
  }

  hdr, ok := segments[0].(Header)
  if !ok {
    return []byte(""), errors.New("Could not find header")
  }

  buf.WriteString(hdr[0])
  buf.WriteString(hdr[1])

  dr, err := regexp.Compile(
    fmt.Sprintf(`([%s])`,
      regexp.QuoteMeta(fmt.Sprintf(
        "%c%c%c%c%c",
        hdr.ComponentDelimiter(),
        hdr.DataDelimiter(),
        hdr.ReleaseIndicator(),
        hdr.RepetitionDelimiter(),
        hdr.SegmentTerminator()),
      )))
  if err != nil {
    return []byte(""), err
  }

  for _, segment := range segments[1:] {
    //buf.WriteString(segment[0])
    //buf.WriteByte(e.DataDelimiter)

    vOf := reflect.ValueOf(segment)

    segmentBytes, err := marshalPart(hdr, vOf, hdr.DataDelimiter(), dr, func(data reflect.Value) ([]byte, error) {
      sep := hdr.ComponentDelimiter()
      // if we have a slice and the slice's elements are of type slice as well,
      // then we have a repetition and we need to use the reserved(repetition) delimiter
      if data.Elem().Kind() == reflect.Slice {
        vOftmp := data.Elem().Index(0)
        if vOftmp.Kind() == reflect.Interface {
          vOftmp = vOftmp.Elem()
        }
        if vOftmp.Kind() == reflect.Slice {
          sep = hdr.RepetitionDelimiter()
        }
      }

      return marshalPart(hdr, data, sep, dr, func(data2 reflect.Value) ([]byte, error) {
        return marshalPart(hdr, data2, hdr.ComponentDelimiter(), dr, func(data3 reflect.Value) ([]byte, error) {
          return marshalPart(hdr, data3, hdr.RepetitionDelimiter(), dr, nil)
        })
      })
    })

    if err != nil {
      return []byte(""), err
    }

    buf.Write(segmentBytes)
    buf.WriteByte(hdr.SegmentTerminator())
  }

  return buf.Bytes(), nil
}
