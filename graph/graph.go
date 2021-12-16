package graph

import (
	"container/heap"
)

type NodeWithCost[T any] struct {
	Node T
	Cost int
}

type THeap[T comparable] []NodeWithCost[T]

func (h THeap[T]) Len() int           { return len(h) }
func (h THeap[T]) Less(i, j int) bool { return h[i].Cost < h[j].Cost }
func (h THeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *THeap[T]) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(NodeWithCost[T]))
}

func (h *THeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Traversable[T comparable] interface {
	Neighbors(node T) []NodeWithCost[T]
}

func Dijkstra[T comparable](graph Traversable[T], start T, end T) (int, []T) {
	costs := make(map[T]int)
	prev := make(map[T]T)

	fringe := THeap[T]{NodeWithCost[T]{start, 0}}
	heap.Init(&fringe)

	for len(fringe) > 0 {
		cur := heap.Pop(&fringe).(NodeWithCost[T])
		node, cost := cur.Node, cur.Cost
		if node == end {
			path := []T{node}
			for path[len(path)-1] != start {
				prevNode, ok := prev[node]
				if !ok {
					panic(node)
				}
				node = prevNode
				path = append(path, node)
			}
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return cost, path
		}

		for _, n := range graph.Neighbors(node) {
			next, edgeCost := n.Node, n.Cost
			if _, ok := prev[next]; !ok {
				prev[next] = node
				nextCost := cost + edgeCost
				costs[next] = nextCost
				nextNode := NodeWithCost[T]{next, nextCost}
				heap.Push(&fringe, nextNode)
			}
		}
	}

	return -1, nil
}
