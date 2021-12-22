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
	x, y, z Interval
	isOn    bool
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

// Closed interval
type Interval struct {
	min, max int
}

func (iv Interval) String() string {
	return fmt.Sprintf("(%d, %d)", iv.min, iv.max)
}

func (iv Interval) IsEmpty() bool {
	return iv.max > iv.min
}

func (a Interval) Intersect(b Interval) Interval {
	return Interval{
		min: max(a.min, b.min),
		max: min(a.max, b.max),
	}
}

func (a Interval) Intersects(b Interval) bool {
	return a.min <= b.max && b.min <= a.max
}

func (a Interval) Union(b Interval) []Interval {
	if a.Intersects(b) {
		return []Interval{{
			min: min(a.min, b.min),
			max: max(a.max, b.max),
		}}
	}
	return []Interval{a, b}
}

// func IntersectRange(aMin, aMax, bMin, bMax int) int, int {
//
// }

func (c Cuboid) Clip() Cuboid {
	return Cuboid{
		x:    c.x.Intersect(Interval{-50, 50}),
		y:    c.y.Intersect(Interval{-50, 50}),
		z:    c.z.Intersect(Interval{-50, 50}),
		isOn: c.isOn,
	}
}

func Set(c Cuboid, grid map[Coord]bool) {
	for x := c.x.min; x <= c.x.max; x++ {
		for y := c.y.min; y <= c.y.max; y++ {
			for z := c.z.min; z <= c.z.max; z++ {
				grid[Coord{x, y, z}] = c.isOn
			}
		}
	}
}

func ParseLine(line string) Cuboid {
	var c Cuboid
	var onOff string
	_, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
		&onOff, &c.x.min, &c.x.max, &c.y.min, &c.y.max, &c.z.min, &c.z.max,
	)
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
