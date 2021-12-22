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

func (iv Interval) Length() int {
	if iv.min <= iv.max {
		return iv.max - iv.min + 1
	}
	return 0
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

// Return the parts of a that are not in b
func (a Interval) Subtract(b Interval) []Interval {
	if b.max < a.min || b.min > a.max {
		return []Interval{a}
	}
	if a.min < b.min && a.max > b.max {
		return []Interval{
			{a.min, b.min - 1},
			{b.max + 1, a.max},
		}
	}
	if a.min < b.min {
		return []Interval{
			{a.min, b.min - 1},
		}
	}
	if b.max < a.max {
		return []Interval{
			{b.max + 1, a.max},
		}
	}
	return []Interval{}
}

// Returns a-b, a&b, b-a
func (a Interval) Overlaps(b Interval) []Interval {
	if !a.Intersects(b) {
		return []Interval{a, b}
	}
	ranges := append(b.Subtract(a), a.Intersect(b))
	ranges = append(ranges, a.Subtract(b)...)
	return ranges
}

func (c Cuboid) Clip() Cuboid {
	return Cuboid{
		x: c.x.Intersect(Interval{-50, 50}),
		y: c.y.Intersect(Interval{-50, 50}),
		z: c.z.Intersect(Interval{-50, 50}),
	}
}

func (c Cuboid) Intersects(other Cuboid) bool {
	return c.x.Intersects(other.x) && c.y.Intersects(other.y) && c.z.Intersects(other.z)
}

func (c Cuboid) Volume() int {
	return c.x.Length() * c.y.Length() * c.z.Length()
}

func (c Cuboid) Union(other Cuboid) []Cuboid {
	var result []Cuboid
	for _, x := range c.x.Overlaps(other.x) {
		for _, y := range c.y.Overlaps(other.y) {
			for _, z := range c.z.Overlaps(other.z) {
				part := Cuboid{x, y, z}
				if part.Intersects(c) || part.Intersects(other) {
					result = append(result, c)
				}
			}
		}
	}
	return result
}

func (c Cuboid) Subtract(other Cuboid) []Cuboid {
	if !c.Intersects(other) {
		return []Cuboid{c}
	}

	var result []Cuboid
	for _, x := range c.x.Overlaps(other.x) {
		for _, y := range c.y.Overlaps(other.y) {
			for _, z := range c.z.Overlaps(other.z) {
				part := Cuboid{x, y, z}
				if part.Volume() > 0 && part.Intersects(c) && !part.Intersects(other) {
					result = append(result, c)
				}
			}
		}
	}

	return result
}

// cuboids are disjoint
func SubtractCuboid(cuboids []Cuboid, sub Cuboid) []Cuboid {
	var result []Cuboid

	for _, cuboid := range cuboids {
		if !cuboid.Intersects(sub) {
			result = append(result, cuboid)
		} else {
			result = append(result, cuboid.Subtract(sub)...)
		}
	}
	return result
}

// cuboids are disjoint
func AddCuboid(cuboids []Cuboid, add Cuboid) []Cuboid {
	remains := SubtractCuboid(cuboids, add)
	remains = append(remains, add)
	return remains
}

func Set(c Cuboid, grid map[Coord]bool) {
	for x := c.x.min; x <= c.x.max; x++ {
		for y := c.y.min; y <= c.y.max; y++ {
			for z := c.z.min; z <= c.z.max; z++ {
				grid[Coord{x, y, z}] = true
			}
		}
	}
}

func ParseLine(line string) (bool, Cuboid) {
	var c Cuboid
	var onOff string
	_, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
		&onOff, &c.x.min, &c.x.max, &c.y.min, &c.y.max, &c.z.min, &c.z.max,
	)
	if err != nil {
		panic(err)
	}
	return onOff == "on", c
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	var disjointCuboids []Cuboid

	for i, line := range linesText {
		isOn, cuboid := ParseLine(line)
		if isOn {
			disjointCuboids = AddCuboid(disjointCuboids, cuboid)
		} else {
			disjointCuboids = SubtractCuboid(disjointCuboids, cuboid)
		}
		disjointCuboids = util.Filter(disjointCuboids, func(c Cuboid) bool { return c.Volume() > 0 })

		fmt.Printf("%d %s -> %d disjoint cuboids\n", i, line, len(disjointCuboids))
	}

	num := 0
	for _, cuboid := range disjointCuboids {
		num += cuboid.Volume()
	}

	fmt.Printf("Num cells on: %d\n", num)
}
