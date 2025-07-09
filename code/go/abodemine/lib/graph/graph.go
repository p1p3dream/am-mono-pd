package graph

import "gonum.org/v1/gonum/graph"

// Graph represents a generic graph structure.
type Graph struct {
	nodes      []*Node
	nodeLabels map[string]*Node
	edges      map[[2]int64]*Edge
}

// New initializes a new Graph.
func New() *Graph {
	return &Graph{
		nodeLabels: make(map[string]*Node),
	}
}

// Node satisfies graph.Graph.
func (g *Graph) Node(id int64) graph.Node {
	if id >= int64(len(g.nodes)) {
		return nil
	}

	return g.nodes[id]
}

// Nodes satisfies graph.Graph.
func (g *Graph) Nodes() graph.Nodes {
	return &Nodes{items: g.nodes}
}

// From satisfies graph.Graph.
func (g *Graph) From(id int64) graph.Nodes {
	return &Nodes{
		items: g.nodes[id].from,
	}
}

func (g *Graph) edgeBetween(uid, vid int64) *Edge {
	if g.edges == nil {
		g.edges = make(map[[2]int64]*Edge)

		for i := range g.nodes {
			node := g.nodes[i]
			for j := range node.from {
				to := node.from[j]

				g.edges[[2]int64{node.id, to.id}] = &Edge{
					from: node,
					to:   to,
				}
			}
		}
	}

	if edge := g.edges[[2]int64{uid, vid}]; edge != nil {
		return edge
	}

	return g.edges[[2]int64{vid, uid}]
}

// HasEdgeBetween satisfies graph.Graph.
func (g *Graph) HasEdgeBetween(xid, yid int64) bool {
	return g.edgeBetween(xid, yid) != nil || g.edgeBetween(yid, xid) != nil
}

// Edge satisfies graph.Graph.
func (g *Graph) Edge(uid, vid int64) graph.Edge {
	if edge := g.edgeBetween(uid, vid); edge != nil {
		return edge
	}

	return g.edgeBetween(vid, uid)
}

// HasEdgeFromTo satisfies graph.Directed.
func (g *Graph) HasEdgeFromTo(uid, vid int64) bool {
	return g.edgeBetween(uid, vid) != nil
}

// To satisfies graph.Directed.
func (g *Graph) To(id int64) graph.Nodes {
	return &Nodes{
		items: g.nodes[id].to,
	}
}

// GetNode returns the node with label = label, or
// creates a new one with label and returns it.
func (g *Graph) GetNode(label string) *Node {
	n, ok := g.nodeLabels[label]
	if ok {
		return n
	}

	n = &Node{
		id:   int64(len(g.nodes)),
		name: label,
	}

	g.nodes = append(g.nodes, n)
	g.nodeLabels[label] = n

	return n
}
