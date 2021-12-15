package graph

import (
	"aoc/util"
	"testing"
)

type UnweightedGraph map[string][]string

func (g UnweightedGraph) Neighbors(n string) []NodeWithCost[string] {
	nodes, ok := g[n]
	if !ok {
		return []NodeWithCost[string]{}
	}
	edges := util.Map(nodes, func(n string) NodeWithCost[string] {
		return NodeWithCost[string]{Node: n, Cost: 1}
	})
	return edges
}

func TestDijkstra(t *testing.T) {
	g := UnweightedGraph{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
	}

	actual := Dijkstra[string](g, "a", "d")
	if actual != 2 {
		t.Errorf("Dijkstra(g, a, b) = %d want 2", actual)
	}
}
