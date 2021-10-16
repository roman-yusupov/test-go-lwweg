// Package graph implements type Graph in LWW-Element-Set model.
package graph

import ()

// Vertex type represents a vertex of the graph.
type Vertex struct {
	id string // represents unique ID of the vertex globally, over all computer systems making updates
}

// Edge type represents an edge of the graph.
type Edge struct {
	v1        Vertex // one of two vertices of the edge
	v2        Vertex // one of two vertices of the edge
	timestamp int64  // represents general form of timestamp for the model
}

// Graph type represents a graph in LWW-Element model.
type Graph struct {
	aSet []Edge // the added set
	rSet []Edge // the removed set
}

// equals checks if two vertices are qual using unique vertex id.
func (v1 Vertex) equals(v2 Vertex) bool {
	return v1.id == v2.id
}

// equals checks if two edges are qual using based on equivalence of their vertices.
// Assumed that the edges are not directional.
func (e1 Edge) equals(e2 Edge) bool {
	return e1.v1.equals(e2.v1) && e1.v2.equals(e2.v2) || e1.v1.equals(e2.v2) && e1.v2.equals(e2.v1)
}

// NewGraph returns new empty graph with zero number of edges.
func NewGraph() *Graph {
	return &Graph{
		aSet: make([]Edge, 0),
		rSet: make([]Edge, 0),
	}
}

// findInASet returns true if the edge included into added set.
func (g *Graph) findInASet(e Edge) (*Edge, bool) {
	for _, el := range g.aSet {
		if el.equals(e) {
			return &el, true
		}
	}

	return nil, false
}

// findInASet returns true if the edge is included into removed set.
func (g *Graph) findInRSet(e Edge) (*Edge, bool) {
	for _, el := range g.rSet {
		if el.equals(e) {
			return &el, true
		}
	}
	return nil, false
}

// AddEdge adds edge e to the graph g.
func (g *Graph) AddEdge(e Edge) {

	ea, ok := g.findInASet(e)
	if !ok {
		g.aSet = append(g.aSet, e)
	} else if ea.timestamp < e.timestamp {
		ea.timestamp = e.timestamp // keeping the lastest add timestamp
	}

}

// RemoveEdge removes edge e from the graph g.
func (g *Graph) RemoveEdge(e Edge) {

	er, ok := g.findInRSet(e)
	if !ok {
		g.rSet = append(g.rSet, e)
	} else if er.timestamp < e.timestamp {
		er.timestamp = e.timestamp // keeping the lastest remove timestamp
	}

}

// Merge combines graph g with another graph g2 using LWW-Element-Set model.
// Graph g is being updated and returned as result.
// Merging is done be merging added sets and deleted sets of two graphs.
func (g *Graph) Merge(g2 *Graph) *Graph {

	for _, e := range g2.aSet {
		g.AddEdge(e)
	}

	for _, e := range g2.rSet {
		g.RemoveEdge(e)
	}

	return g
}

// CheckInGraph checks if the edge e belongs to the graph g using both added set and deleted set in .
// The algorithm is biased towards added for equivalent timestamps in added and deleted sets.
func (g *Graph) CheckInGraph(e Edge) bool {

	ea, ok := g.findInASet(e)
	if !ok {
		return false
	}

	er, ok := g.findInRSet(e)
	if !ok {
		return true
	}

	if ea.timestamp < er.timestamp { // biased towards adds for equivalent timestamps
		return false
	}

	return true
}

// currentGraph returns the current state of the graph according to LWW-Element-Set model.
// The algorithm is biased towards added for equivalent timestamps in added and deleted sets.
func (g *Graph) currentGraph() []Edge {
	result := make([]Edge, 0)

	for _, ea := range g.aSet {

		er, ok := g.findInRSet(ea)
		if ok && ea.timestamp < er.timestamp { // biased towards adds for equivalent timestamps
			continue
		}

		result = append(result, ea)
	}

	return result
}

// findConnectedVertices returns all connected vertices to vertex v for given array of edges represendng a graph.
func findConnectedVertices(edges []Edge, v Vertex) []Vertex {
	result := make([]Vertex, 0)

	for _, e := range edges {
		if v.equals(e.v1) {
			result = append(result, e.v2)
		} else if v.equals(e.v2) {
			result = append(result, e.v1)
		}
	}

	return result
}

// FindConnected returns all connected vertices to vertex v for a graph g.
func (g *Graph) FindConnected(v Vertex) []Vertex {
	return findConnectedVertices(g.currentGraph(), v)
}

type Path []Vertex

// FindPath returns a path from vertex v1 to vertex v2.
func (g *Graph) FindPath(v1 Vertex, v2 Vertex) []Vertex {

	edges := g.currentGraph()

	branches := make(map[Vertex]Path, 0)
	branches[v1] = Path{v1}

	visited := make([]Vertex, 0)
	visited = append(visited, v1)

	for {

		// path not found
		if len(branches) == 0 {
			return Path{}
		}

		// v2 reached
		for vertex, path := range branches {
			if vertex.equals(v2) {
				return path
			}
		}

		// next step of iterations over connected vertices

		newbranches := make(map[Vertex]Path, 0)

		for vertex, path := range branches {

			allconnected := findConnectedVertices(edges, vertex)

			for _, connected := range allconnected {

				// excluding already for visited
				skip := false

				for _, v := range visited {
					if connected.equals(v) {
						skip = true
						break
					}
				}

				if skip {
					continue
				}

				// adding to visited
				visited = append(visited, connected)

				// extending path
				newbranches[connected] = append(path, connected)

			}

		}

		// step finished
		branches = newbranches

	}

}
