package graph

import "gonum.org/v1/gonum/graph"

// Nodes represents a collection of Node elements.
// https://godoc.org/gonum.org/v1/gonum/graph#Nodes
type Nodes struct {
	index int
	items []*Node
}

// Next satisfies graph.Iterator.
func (n *Nodes) Next() bool {
	if len(n.items) == 0 || n.index >= len(n.items) {
		return false
	}

	n.index++

	return true
}

// Len satisfies graph.Iterator.
func (n *Nodes) Len() int {
	return len(n.items) - n.index
}

// Reset satisfies graph.Iterator.
func (n *Nodes) Reset() {
	n.index = 0
}

// Node satisfies graph.Nodes.
func (n *Nodes) Node() graph.Node {
	return n.items[n.index-1]
}
