// Parse nodes

package parse

import (
  "bytes"
  "fmt"
)

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
  Type() NodeType
  String() string
}

// NodeType identifies the type of a parse tree node.
type NodeType int

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
  return t
}

const (
  NodeList       NodeType = iota // A list of Nodes.
  NodeHeader                     // The UNA header node.
  NodeSegment                    // A segment node.
  NodeData                       // A data node.
  NodeText                       // Plain text.
  NodeComponent                  // A component node.
  NodeRepetition                 // A repetition node.
)

// ListNode holds a sequence of nodes.
type ListNode struct {
  NodeType
  Nodes []Node // The element nodes in lexical order.
}

func newList() *ListNode {
  return &ListNode{NodeType: NodeList}
}

func (l *ListNode) append(n Node) {
  l.Nodes = append(l.Nodes, n)
}

// Returns a flat string of the ListNode.
func (l *ListNode) String() string {
  b := new(bytes.Buffer)
  for _, n := range l.Nodes {
    fmt.Fprint(b, n)
  }
  return b.String()
}

// Prints out the nodes with a delimiter.
func (l *ListNode) StringDelim(delim string) string {
  b := new(bytes.Buffer)
  for c, n := range l.Nodes {
    fmt.Fprint(b, n)

    if c < len(l.Nodes)-1 {
      fmt.Fprint(b, delim)
    }
  }
  return b.String()
}

// A header node holds the UNA header and the
// configuration text.
type HeaderNode struct {
  NodeType
  SegmentName Node
  Text        Node
}

func newHeader() *HeaderNode {
  return &HeaderNode{NodeType: NodeSegment}
}

func (h *HeaderNode) String() string {
  return fmt.Sprint(h.SegmentName, h.Text)
}

// SegmentNode holds a list of nodes (only
// DataNodes) that signify an entire segment.
// The segment name is specified in the first
// DataNode of the list.
type SegmentNode struct {
  NodeType
  List *ListNode
}

func newSegment() *SegmentNode {
  return &SegmentNode{NodeType: NodeSegment, List: newList()}
}

func (s *SegmentNode) String() string {
  return fmt.Sprint(s.List.StringDelim("|"))
}

// DataNode holds a single node. It's basically
// an encapsulation to follow the spec.
type DataNode struct {
  NodeType
  Node Node
}

func newData(node Node) *DataNode {
  return &DataNode{NodeType: NodeData, Node: node}
}

func (c *DataNode) String() string {
  return fmt.Sprint(c.Node)
}

// TextNode holds plain text.
type TextNode struct {
  NodeType
  Text []byte // The text; may span newlines.
}

func newText(text string) *TextNode {
  return &TextNode{NodeType: NodeText, Text: []byte(text)}
}

func (t *TextNode) String() string {
  return fmt.Sprintf("%s", t.Text)
}

// ComponentNode holds a list of nodes that are
// components of one another. Like a person's
// first and last name can be components.
type ComponentNode struct {
  NodeType
  List *ListNode
}

func newComponent() *ComponentNode {
  return &ComponentNode{NodeType: NodeComponent, List: newList()}
}

func (c *ComponentNode) String() string {
  return fmt.Sprint(c.List.StringDelim("~"))
}

// RepetitionNode holds a list of nodes that are
// repetitions. For example, a start and end date
// could be repetitions.
type RepetitionNode struct {
  NodeType
  List *ListNode
}

func newRepetition() *RepetitionNode {
  return &RepetitionNode{NodeType: NodeRepetition, List: newList()}
}

func (r *RepetitionNode) String() string {
  return fmt.Sprint(r.List.StringDelim("^"))
}
