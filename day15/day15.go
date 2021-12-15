package main

import (
	c "aoc/coord"
	"aoc/set"
	"aoc/util"
	"fmt"
	"os"
	"strconv"
)

type Path struct {
	visited set.Set[c.Coord]
	tip     c.Coord
	risk    int
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	risks := make(map[c.Coord]int)
	for y, line := range linesText {
		for x, digit := range line {
			val, err := strconv.Atoi(string(digit))
			if err != nil {
				panic(err)
			}
			risks[c.Coord{X: x, Y: y}] = val
		}
	}

	end := c.MaxXY(risks)

	minRisk := 0
	// Come up with a good upper bound for minimum risk
	for x := 1; x <= end.X; x++ {
		v, ok := risks[c.Coord{X: x, Y: 0}]
		if !ok {
			panic(x)
		}
		minRisk += v
	}
	for y := 1; y <= end.X; y++ {
		v, ok := risks[c.Coord{X: end.X, Y: y}]
		if !ok {
			panic(y)
		}
		minRisk += v
	}
	fmt.Printf("Min Risk: %d\n", minRisk)

	init := Path{tip: c.Coord{X: 0, Y: 0}}
	pool := []Path{init}
	n := 0
	for len(pool) > 0 {
		fringe := []Path{}
		path := pool[0]
		pool = pool[1:]
		for _, next := range path.tip.Neighbors4() {
			if path.visited.Contains(next) {
				continue
			}
			risk, ok := risks[next]
			if !ok {
				continue // off the grid
			}
			nextRisk := risk + path.risk
			if nextRisk < minRisk {
				if next == end {
					minRisk = nextRisk
				}
				fringe = append(fringe, Path{tip: next, risk: nextRisk, visited: path.visited.CloneAndAdd(next)})
			}
		}

		n += 1
		fmt.Printf("iteration %d, pool size=%d\n", n, len(fringe))
		pool = fringe
	}
	fmt.Printf("Min Risk: %d\n", minRisk)
}
