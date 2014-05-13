package node

type NodeType int

const (
	NodeSegment NodeType = iota
	NodeElement
	NodeComponent
)

// type Node struct {
// 	Type NodeType
// 	Tag  string
// 	Data string

// 	Nodes []Node
// 	Sub   Node
// }

type Segment struct {
	Tag      string
	Data     []interface{}
	Children []Segment
}

type Component []string
