package main

import (
	"aoc/util"
	"fmt"
	"os"
)

type Coord struct {
	x, y, z int
}

type Cuboid struct {
	min  Coord
	max  Coord
	isOn bool
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

// func IntersectRange(aMin, aMax, bMin, bMax int) int, int {
//
// }

func (c Cuboid) Clip() Cuboid {
	return Cuboid{
		min: Coord{
			x: max(-50, c.min.x),
			y: max(-50, c.min.y),
			z: max(-50, c.min.z),
		},
		max: Coord{
			x: min(50, c.max.x),
			y: min(50, c.max.y),
			z: min(50, c.max.z),
		},
		isOn: c.isOn,
	}
}

func Set(c Cuboid, grid map[Coord]bool) {
	for x := c.min.x; x <= c.max.x; x++ {
		for y := c.min.y; y <= c.max.y; y++ {
			for z := c.min.z; z <= c.max.z; z++ {
				grid[Coord{x, y, z}] = c.isOn
			}
		}
	}
}

func ParseLine(line string) Cuboid {
	var c Cuboid
	var onOff string
	_, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d", &onOff, &c.min.x, &c.max.x, &c.min.y, &c.max.y, &c.min.z, &c.max.z)
	if err != nil {
		panic(err)
	}
	if onOff == "on" {
		c.isOn = true
	}
	return c
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	cuboids := util.Map(linesText, ParseLine)

	grid := map[Coord]bool{}
	for _, c := range cuboids {
		Set(c.Clip(), grid)
	}

	num := 0
	for _, v := range grid {
		if v {
			num++
		}
	}

	fmt.Printf("Num cells on: %d\n", num)
}
