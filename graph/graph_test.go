package graph

import (
	"testing"
)

// Tests method Findconnected of type Graph
func TestFindConnected(t *testing.T) {

	g := NewGraph()

	g.AddEdge(
		Edge{
			Vertex{"1"},
			Vertex{"2"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"3"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"3"},
			Vertex{"4"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"5"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"6"},
			1000,
		})

	g.RemoveEdge(
		Edge{
			Vertex{"1"},
			Vertex{"2"},
			1000, // same timestamp, it should not "remove" from the final graph in LWW-Element-Set model
		})

	g.RemoveEdge(
		Edge{
			Vertex{"2"},
			Vertex{"6"},
			1001, // bigger timestamp, it should "remove"
		})

	// expected connected vertices [{1} {3} {5}]
	result := len(g.FindConnected(Vertex{"2"}))
	expected := 3

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

}

// Tests methods AddEdge and RemoveEdge of type Graph
func TestAddRemove(t *testing.T) {

	g := NewGraph()

	g.AddEdge(
		Edge{
			Vertex{"1"},
			Vertex{"2"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"3"},
			1000,
		})

	result := len(g.FindConnected(Vertex{"2"}))
	expected := 2

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

	g.RemoveEdge(
		Edge{
			Vertex{"3"},
			Vertex{"2"},
			1001,
		})

	result = len(g.FindConnected(Vertex{"2"}))
	expected = 1

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

}

// Tests method CheckInGraph of type Graph
func TestCheckInGraph(t *testing.T) {

	g := NewGraph()

	g.AddEdge(
		Edge{
			Vertex{"1"},
			Vertex{"2"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"3"},
			1000,
		})

	g.RemoveEdge(
		Edge{
			Vertex{"3"},
			Vertex{"2"},
			1001,
		})

	e := Edge{
		Vertex{"1"},
		Vertex{"2"},
		0,
	}

	result := g.CheckInGraph(e)
	expected := true

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

	e = Edge{
		Vertex{"2"},
		Vertex{"1"},
		0,
	}

	result = g.CheckInGraph(e)
	expected = true

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

}

// Tests method Merge of type Graph
func TestMerge(t *testing.T) {

	g1 := NewGraph()

	g1.AddEdge(
		Edge{
			Vertex{"1"},
			Vertex{"2"},
			1000,
		})

	g1.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"3"},
			1000,
		})

	g2 := NewGraph()

	g1.RemoveEdge(
		Edge{
			Vertex{"2"},
			Vertex{"3"},
			1001,
		})

	g1.AddEdge(
		Edge{
			Vertex{"3"},
			Vertex{"4"},
			1001,
		})

	g1.Merge(g2)

	result := len(g1.FindConnected(Vertex{"2"}))
	expected := 1

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

}

// Tests method FindPath of type Graph
func TestFindPath(t *testing.T) {

	g := NewGraph()

	g.AddEdge(
		Edge{
			Vertex{"1"},
			Vertex{"2"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"2"},
			Vertex{"3"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"3"},
			Vertex{"4"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"4"},
			Vertex{"5"},
			1000,
		})

	g.AddEdge(
		Edge{
			Vertex{"1"},
			Vertex{"5"},
			1000,
		})

	g.RemoveEdge(
		Edge{
			Vertex{"4"},
			Vertex{"5"},
			1001,
		})

	path := g.FindPath(Vertex{"2"}, Vertex{"5"})

	// expected path is only [{2} {1} {5}]
	result := len(path)
	expected := 3

	t.Log(path)

	if result != expected {
		t.Fatalf("Expected \"%v\", returned \"%v\"", expected, result)
	}

}
