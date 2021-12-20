package main

import (
	c "aoc/coord"
	"aoc/util"
	"fmt"
	"os"
	"strings"
)

func ReadCell(cell string) int {
	if cell == "#" {
		return 1
	} else if cell == "." {
		return 0
	}
	panic(cell)
}

func PrintCell(cell int) string {
	if cell == 0 {
		return "."
	} else if cell == 1 {
		return "#"
	}
	panic(cell)
}

func Advance(grid map[c.Coord]int, bg int, decoder []int) (map[c.Coord]int, int) {
	min := c.MinXY(grid)
	max := c.MaxXY(grid)
	fmt.Printf("min: %s, max: %s\n", min, max)

	next := map[c.Coord]int{}
	for x := min.X - 1; x <= max.X+1; x++ {
		for y := min.Y - 1; y <= max.Y+1; y++ {
			idx := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					idx *= 2
					v, ok := grid[c.Coord{X: x + dx, Y: y + dy}]
					if !ok {
						v = bg
					}
					if v != 0 {
						idx++
					}
				}
			}
			// fmt.Printf("%d,%d -> %d\n", x, y, idx)
			v := decoder[idx]
			next[c.Coord{X: x, Y: y}] = v
		}
	}
	if bg == 0 {
		bg = decoder[0]
	} else if bg == 1 {
		bg = decoder[511]
	} else {
		panic(bg)
	}
	return next, bg
}

func main() {
	chunks := util.ReadChunks(os.Args[1])
	if len(chunks) != 2 {
		panic(len(chunks))
	}

	decoder := util.Map(strings.Split(chunks[0][0], ""), ReadCell)
	if len(decoder) != 512 {
		panic(chunks[0][0])
	}

	grid := map[c.Coord]int{}
	for y, line := range chunks[1] {
		for x, cell := range strings.Split(line, "") {
			if cell == "#" {
				grid[c.Coord{X: x, Y: y}] = 1
			}
		}
	}
	bg := 0

	for step := 1; step <= 50; step++ {
		grid, bg = Advance(grid, bg, decoder)
	}

	// c.PrintGrid(grid, ".", PrintCell)
	// fmt.Printf("Background: %d\n", bg)
	// grid, bg = Advance(grid, bg, decoder)
	// c.PrintGrid(grid, ".", PrintCell)
	// fmt.Printf("Background: %d\n", bg)

	if bg != 0 {
		panic(bg)
	}
	num := 0
	for _, v := range grid {
		if v != 0 {
			num++
		}
	}
	fmt.Printf("Num Set: %d\n", num)
}

// 10,000 = too high
//  5,991 = too high
//  5,469 = too high
//  5,372 = wrong
