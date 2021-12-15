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

	minRisks := map[c.Coord]int{}
	for k := range risks {
		minRisks[k] = -1
	}
	minRisks[c.Coord{X: 0, Y: 0}] = 0
	fringe := []c.Coord{{X: 0, Y: 0}}
	n := 0
	for len(fringe) > 0 {
		nextFringe := []c.Coord{}
		for _, tip := range fringe {
			tipRisk := minRisks[tip]
			for _, next := range tip.Neighbors4() {
				risk, ok := risks[next]
				if !ok {
					continue // off the grid
				}
				nextRisk := risk + tipRisk
				prevRisk := minRisks[next]
				if nextRisk < prevRisk || prevRisk == -1 {
					minRisks[next] = nextRisk
					if next == end {
						minRisk = nextRisk
					}
					nextFringe = append(nextFringe, next)
				}
			}
		}
		n += 1
		fmt.Printf("iteration %d, pool size=%d\n", n, len(fringe))
		fringe = nextFringe
	}
	fmt.Printf("Min Risk: %d\n", minRisk)
}
