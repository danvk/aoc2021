package main

import (
	"aoc/util"
	"flag"
	"fmt"
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
		{x + 1, y + 1},
		{x + 1, y - 1},
		{x - 1, y + 1},
		{x - 1, y - 1},
	}
}

func PrintGrid(grid map[Coord]int) {
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			fmt.Printf("%d", grid[Coord{x, y}])
		}
		fmt.Printf("\n")
	}
}

func AdvanceOneStep(grid map[Coord]int) int {
	flashed := make(map[Coord]bool)

	for c := range grid {
		grid[c] += 1
	}

	totalFlashes := 0
	newFlashes := 1
	for newFlashes > 0 {
		newFlashes = 0

		for c, v := range grid {
			if v > 9 && !flashed[c] {
				newFlashes += 1
				flashed[c] = true
				for _, pos := range neighbors(c) {
					grid[pos] += 1
				}
			}
		}

		totalFlashes += newFlashes
	}

	for c, v := range grid {
		if v > 9 {
			grid[c] = 0
		}
	}

	return totalFlashes
}

func main() {
	numSteps := flag.Int("steps", 100, "number of steps to simulate")
	flag.Parse()

	linesText := util.ReadLines(flag.Args()[0])

	octopi := make(map[Coord]int)
	for y, line := range linesText {
		for x, digit := range line {
			val, err := strconv.Atoi(string(digit))
			if err != nil {
				panic(err)
			}
			octopi[Coord{x, y}] = val
		}
	}

	/*
		PrintGrid(octopi)
		AdvanceOneStep(octopi)

		fmt.Printf("\nAfter 1 step:\n")
		PrintGrid(octopi)

		AdvanceOneStep(octopi)
		fmt.Printf("\nAfter 2 steps:\n")
		PrintGrid(octopi)
	*/

	flashes := 0
	for step := 1; step <= *numSteps; step++ {
		flashes += AdvanceOneStep(octopi)
		fmt.Printf("\nAfter %d step(s):\n", step)
		PrintGrid(octopi)
		fmt.Printf("%d flashes\n", flashes)
	}
	fmt.Printf("Flashes after 100 steps: %d\n", flashes)
}
