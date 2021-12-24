package graph

import (
	"aoc/util"
	"reflect"
	"testing"
)

type UnweightedGraph map[string][]string

func (g UnweightedGraph) String(n string) string {
	return n
}

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

	actual, path := Dijkstra[string](g, "a", "d")
	want := []NodeWithCost[string]{
		{"a", 0},
		{"b", 1},
		{"d", 2},
	}
	if actual != 2 || !reflect.DeepEqual(path, want) {
		t.Errorf("Dijkstra(g, a, b) = %d, %#v want 2, %#v", actual, path, want)
	}
}
