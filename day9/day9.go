package main

import (
	c "aoc/coord"
	"aoc/set"
	"aoc/util"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func FindBasinSize(heights map[c.Coord]int, start c.Coord) int {
	basin := set.SetFrom([]c.Coord{start})
	fringe := basin.Clone()

	for len(fringe) > 0 {
		newFringe := set.SetFrom([]c.Coord{})
		for coord := range fringe {
			for _, n := range coord.Neighbors4() {
				v, ok := heights[n]
				if ok && v < 9 && !basin[n] {
					newFringe.Add(n)
					basin.Add(n)
				}
			}
		}
		fringe = newFringe
	}

	return len(basin)
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	heights := make(map[c.Coord]int)
	for y, line := range linesText {
		for x, digit := range line {
			val, err := strconv.Atoi(string(digit))
			if err != nil {
				panic(err)
			}
			heights[c.Coord{X: x, Y: y}] = val
		}
	}

	sumHeights := 0
	basinSizes := []int{}
	mins := []c.Coord{}
	for pos, height := range heights {
		minNeighbor := height + 1
		for _, p := range pos.Neighbors4() {
			if v, ok := heights[p]; ok && v < minNeighbor {
				minNeighbor = v
			}
		}

		if minNeighbor > height {
			mins = append(mins, pos)
			sumHeights += height + 1

			basinSize := FindBasinSize(heights, pos)
			fmt.Printf("%#v basin size: %d\n", pos, basinSize)
			basinSizes = append(basinSizes, basinSize)
		}
	}

	sort.Ints(basinSizes)
	n := len(basinSizes)
	basinProd := basinSizes[n-1] * basinSizes[n-2] * basinSizes[n-3]

	fmt.Printf("Mins: %#v\n", mins)
	fmt.Printf("Sum heights: %d\n", sumHeights)
	fmt.Printf("Prod basin sizes: %d\n", basinProd)
}
