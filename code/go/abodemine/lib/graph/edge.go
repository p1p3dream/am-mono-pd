package graph

import "gonum.org/v1/gonum/graph"

// Edge represents a connection between two nodes.
// https://godoc.org/gonum.org/v1/gonum/graph#Edge
type Edge struct {
	from *Node
	to   *Node
}

// From satisfies graph.Edge.
func (e *Edge) From() graph.Node {
	return e.from
}

// To satisfies graph.Edge.
func (e *Edge) To() graph.Node {
	return e.to
}

// ReversedEdge satisfies graph.Edge.
func (e *Edge) ReversedEdge() graph.Edge {
	return &Edge{
		from: e.to,
		to:   e.from,
	}
}
