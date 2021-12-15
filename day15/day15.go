package main

import (
	c "aoc/coord"
	"aoc/util"
	"fmt"
	"os"
	"strconv"
)

func PrintGrid(risks map[c.Coord]int) {
	c.PrintGrid(risks, ".", func(risk int) string {
		return strconv.Itoa(risk)
	})
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
	PrintGrid(risks)

	end := c.MaxXY(risks)
	maxX, maxY := end.X+1, end.Y+1
	fmt.Printf("maxX, maxY = %d, %d\n", maxX, maxY)

	nr := make(map[c.Coord]int)
	for rx := 0; rx <= 4; rx++ {
		for ry := 0; ry <= 4; ry++ {
			for pos, risk := range risks {
				x, y := pos.X, pos.Y
				nextRisk := (risk + rx + ry)
				if nextRisk > 9 {
					nextRisk -= 9
				}
				nr[c.Coord{X: x + rx*maxX, Y: y + ry*maxY}] = nextRisk
			}
		}
	}
	risks = nr
	PrintGrid(risks)

	end = c.MaxXY(risks)
	maxX, maxY = end.X+1, end.Y+1
	fmt.Printf("maxX, maxY = %d, %d\n", maxX, maxY)

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
