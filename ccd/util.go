package ccd

import (
  "github.com/jteeuwen/go-pkg-xmlx"
  "time"
)

const (
  // Found both these formats in the wild
  TimeDecidingIndex = 14
  TimeFormat        = "20060102150405-0700"
  TimeFormat2       = "20060102150405.000-0700"
)

// Dates and times in a CCD can be partial. Meaning they can be:
//   2006, 200601, 20060102, etc...
// This function helps us parse all cases.
func ParseTime(value string) (time.Time, error) {
  l := len(value)
  tmfmt := TimeFormat
  if l > TimeDecidingIndex && value[TimeDecidingIndex] == '.' {
    tmfmt = TimeFormat2
  }
  return time.Parse(tmfmt[:l], value)
}

// Node get.
// helper function to continually transverse down the
// xml nodes in args, and return the last one.
func Nget(node *xmlx.Node, args ...string) *xmlx.Node {
  for _, a := range args {
    if node == nil {
      return nil
    }

    node = node.SelectNode("*", a)
  }

  return node
}

// Node Safe get.
// just like Nget, but returns a node no matter what.
func Nsget(node *xmlx.Node, args ...string) *xmlx.Node {
  n := Nget(node, args...)
  if n == nil {
    return xmlx.NewNode(0)
  }
  return n
}

func insertSortParser(p Parser, parsers Parsers) Parsers {
  i := len(parsers) - 1
  for ; i >= 0; i-- {
    if p.Priority > parsers[i].Priority {
      i += 1
      break
    }
  }

  if i < 0 {
    i = 0
  }

  parsers = append(parsers, p) // this just expands storage.
  copy(parsers[i+1:], parsers[i:])
  parsers[i] = p

  return parsers
}
