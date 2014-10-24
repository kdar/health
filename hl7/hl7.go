package hl7

import (
	"errors"
)

// Data is an interface the represents a piece of data that
// is indexable.
type Data interface {
	Index(index int) (Data, bool)
}

// Composer is an interface that represents a piece of data that
// is indexable and composable.
type Composer interface {
	Data
	Append(values ...Data) error
}

// Field represents a slice of bytes. This is the most basic type in
// a HL7 message.
// e.g., in: MSH|^~\&|field|component1^component2
// "MSH", "^~\&", "field", "component1", "component2" are all Fields.
type Field []byte

// Index just returns itself at index 0.
func (f Field) Index(index int) (Data, bool) {
	if index != 0 {
		return nil, false
	}

	return f, true
}

func (f Field) String() string {
	return string(f)
}

// SubComponent is a type that contains Fields. Usually separated by '&'.
// e.g. subcomponent1&subcomponent2
type SubComponent []Field

// Index returns the Data at a certain index inside this SubComponent.
func (s SubComponent) Index(index int) (Data, bool) {
	if index < 0 || index >= len(s) {
		return nil, false
	}

	return s[index], true
}

// Append appends Data to this SubComponent.
func (s *SubComponent) Append(values ...Data) error {
	for _, v := range values {
		value, ok := v.(Field)

		if !ok {
			return errors.New(
				"subcomponents may only contain Fields",
			)
		}

		*s = append(*s, value)
	}

	return nil
}

// Component is a type that contains Fields and SubComponents. Usually separated by '^'.
// e.g. component1^component2
type Component []Data

// Index returns the Data at a certain index inside this Component.
func (c Component) Index(index int) (Data, bool) {
	if index < 0 || index >= len(c) {
		return nil, false
	}

	return c[index], true
}

// Append appends Data to this Component.
func (c *Component) Append(values ...Data) error {
	for _, v := range values {
		switch v.(type) {
		case Field, SubComponent:
		default:
			return errors.New(
				"components may only contain Fields and SubComponents",
			)
		}

		*c = append(*c, v)
	}

	return nil
}

// Repeated is a type that contains Fields, SubComponents, and Components.
// Usually separated by '~'.
// e.g. component1a^component2a~component1b^component2b
type Repeated []Data

// Index returns the Data at a certain index inside this Repeated.
func (r Repeated) Index(index int) (Data, bool) {
	if index < 0 || index >= len(r) {
		return nil, false
	}

	return r[index], true
}

// Append appends Data to this Repeated.
func (r *Repeated) Append(values ...Data) error {
	for _, v := range values {
		switch v.(type) {
		case Field, SubComponent, Component:
		default:
			return errors.New(
				"repeated fields may only contain Fields, SubComponents, and Components",
			)
		}

		*r = append(*r, v)
	}

	return nil
}

// Segment is a type that contains Fields, SubComponents, Components, and Repeated.
// Usually separated by '|'.
// e.g. MSH|field1
type Segment []Data

// Index returns the Data at a certain index inside this Segment.
func (s Segment) Index(index int) (Data, bool) {
	if index < 0 || index >= len(s) {
		return nil, false
	}

	return s[index], true
}

// Append appends Data to this Segment.
func (s *Segment) Append(values ...Data) (err error) {
	for _, v := range values {
		switch v.(type) {
		case Field, SubComponent, Component, Repeated:
		default:
			return errors.New(
				"segments may only contain Fields, SubComponents, Components, and Repeated",
			)
		}

		*s = append(*s, v)
	}

	return nil
}
