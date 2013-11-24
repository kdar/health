// Simple stack package for aid the parser
// in shift-reducing parsing.

package parse

type NodeStack []Node

func (n *NodeStack) clear() {
	*n = []Node{}
}

func (n *NodeStack) push(node Node) {
	*n = append(*n, node)
}

func (n *NodeStack) last() Node {
	c := len(*n)
	if c == 0 {
		return nil
	}

	return (*n)[c-1]
}

func (n *NodeStack) setLast(node Node) {
	(*n)[n.len()-1] = node
}

func (n *NodeStack) len() int {
	return len(*n)
}
