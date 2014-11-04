package main

import (
	"github.com/jteeuwen/go-pkg-xmlx"
)

type Node struct {
	*xmlx.Node
}

func NewNode(n *xmlx.Node) *Node {
	return &Node{Node: n}
}

func (node *Node) FindNode(args ...string) *Node {
	n := node.Node
	for _, a := range args {
		n = n.SelectNode("*", a)
	}

	return NewNode(n)
}
