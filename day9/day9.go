package main

import (
	"aoc/util"
	"fmt"
	"os"
	"strconv"
)

type Coord struct {
	X, Y int
}

func neighbors(pos Coord) []Coord {
	x, y := pos.X, pos.Y
	return []Coord{
		{x + 1, y},
		{x - 1, y},
		{x, y - 1},
		{x, y + 1},
	}
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	heights := make(map[Coord]int)
	for y, line := range linesText {
		for x, digit := range line {
			val, err := strconv.Atoi(string(digit))
			if err != nil {
				panic(err)
			}
			heights[Coord{x, y}] = val
		}
	}

	sumHeights := 0
	mins := []Coord{}
	for pos, height := range heights {
		minNeighbor := height + 1
		for _, p := range neighbors(pos) {
			if v, ok := heights[p]; ok && v < minNeighbor {
				minNeighbor = v
			}
		}

		if minNeighbor > height {
			mins = append(mins, pos)
			sumHeights += height + 1
		}
	}

	fmt.Printf("Mins: %#v\n", mins)
	fmt.Printf("Sum heights: %d\n", sumHeights)
}
