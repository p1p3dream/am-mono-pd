package graph

// Node represents a node in the graph.
// https://godoc.org/gonum.org/v1/gonum/graph#Node
type Node struct {
	id   int64
	name string

	// to represents connections TO this node: (other node) -> (this node)
	// from represents connections FROM this node: (this node) -> (other node)
	//
	// The connection should always be considered as unidirectional, and
	// the direction of the connection represents the dependency relation:
	// the element on the right depends on the element on the left.
	to   []*Node
	from []*Node

	value any
}

// ID satisfies graph.Node.
func (n *Node) ID() int64 {
	return n.id
}

// Name returns the node name.
func (n *Node) Name() string {
	return n.name
}

// AddFrom adds an outgoing edge.
func (n *Node) AddFrom(node *Node) {
	n.from = append(n.from, node)
}

// LenFrom return the len of from slice.
func (n *Node) LenFrom() int {
	return len(n.from)
}

// AddTo adds an incoming edge.
func (n *Node) AddTo(node *Node) {
	n.to = append(n.to, node)
}

// LenTo return the len of to slice.
func (n *Node) LenTo() int {
	return len(n.to)
}

// SetValue returns the element associated with the node.
func (n *Node) SetValue(v any) {
	n.value = v
}

// GetValue returns the value associated with the node.
func (n *Node) GetValue() any {
	return n.value
}
