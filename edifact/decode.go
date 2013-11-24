package edifact

import (
	"errors"
	"github.com/kdar/health/edifact/parse"
)

// unmarshals passed byte data
func Unmarshal(data []byte) (Values, error) {
	root, err := parse.Parse(string(data))
	if err != nil {
		return Values{}, err
	}

	s := &state{}
	values := s.walk(root)

	return values, s.err
}

// our state for walking the node tree
type state struct {
	err error
}

// walks a parsed node tree and returns values
func (s *state) walk(node parse.Node) Values {
	var ret Values

	switch node := node.(type) {
	case *parse.ListNode:
		for _, node := range node.Nodes {
			ret = append(ret, s.walk(node)...)
		}
	case *parse.HeaderNode:
		ret = append(ret, Header{node.SegmentName.String(), node.Text.String()})
	case *parse.SegmentNode:
		ret = append(ret, s.walk(node.List))
	case *parse.DataNode:
		ret = append(ret, s.walk(node.Node)...)
	case *parse.ComponentNode:
		ret = append(ret, s.walk(node.List))
	case *parse.RepetitionNode:
		ret = append(ret, s.walk(node.List))
	case *parse.TextNode:
		ret = append(ret, string(node.Text))
	default:
		s.err = errors.New("unknown node")
	}

	return ret
}
