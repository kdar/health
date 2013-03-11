package hl7

type Values []interface{}

type Header struct {
  Name                     string
  CompositeDelimiter       byte
  SubCompositeDelimiter    byte
  SubSubCompositeDelimiter byte
  RepetitionDelimiter      byte
  EscapeCharacter          byte
  Values                   Values
}
