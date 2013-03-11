package hl7

import (
  "reflect"
  "testing"
)

const (
  UM_IN1 = "MSH|^~\\&|EPIC|EPICADT|SMS|SMSADT|199912271408|CHARRIS|ADT^A04|1817457|D|2.5|\rPID||0493575^^^2^ID 1|454721||DOE^JOHN^^^^|DOE^JOHN^^^^|19480203|M||B|254 MYSTREET AVE^^MYTOWN^OH^44123^USA||(216)123-4567|||M|NON|400003403~1129086|\rNK1||ROE^MARIE^^^^|SPO||(216)123-4567||EC|||||||||||||||||||||||||||\rPV1||O|168 ~219~C~PMA^^^^^^^^^||||277^ALLEN MYLASTNAME^BONNIE^^^^|||||||||| ||2688684|||||||||||||||||||||||||199912271408||||||002376853"
)

var (
  UM_OUT1 = Values{}
)

var unmarshalTests = []struct {
  in  []byte
  out Values
}{
  {[]byte(UM_IN1), UM_OUT1},
}

func TestUnmarshal(t *testing.T) {
  for i, tt := range unmarshalTests {
    out, err := Unmarshal(tt.in)
    if err != nil {
      t.Fatalf("%d. received error: %s", i, err)
    }

    if !reflect.DeepEqual(out, tt.out) {
      t.Fatalf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, out, tt.out)
    }
  }
}
