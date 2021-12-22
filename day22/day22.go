package main

import (
	"aoc/util"
	"fmt"
	"os"
	"time"
)

type Coord struct {
	x, y, z int
}

type Cuboid struct {
	x, y, z Interval
}

type Rect struct {
	x, y Interval
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

func (a Rect) Intersects(b Rect) bool {
	return a.x.Intersects(b.x) && a.y.Intersects(b.y)
}

func (a Rect) Intersection(b Rect) Rect {
	return Rect{
		x: a.x.Intersect(b.x),
		y: a.y.Intersect(b.y),
	}
}

func (a Rect) Area() int {
	return a.x.Length() * a.y.Length()
}

func (r Rect) Union(other Rect) []Rect {
	var result []Rect
	for _, x := range r.x.Overlaps(other.x) {
		for _, y := range r.y.Overlaps(other.y) {
			part := Rect{x, y}
			if part.Intersects(r) || part.Intersects(other) {
				result = append(result, r)
			}
		}
	}
	return result
}

func (r Rect) Subtract(other Rect) []Rect {
	if !r.Intersects(other) {
		return []Rect{r}
	}

	var result []Rect
	for _, x := range r.x.Overlaps(other.x) {
		for _, y := range r.y.Overlaps(other.y) {
			part := Rect{x, y}
			if part.Area() > 0 && part.Intersects(r) && !part.Intersects(other) {
				result = append(result, r)
			}
		}
	}

	return result
}

func SubtractRects(rects []Rect, sub Rect) []Rect {
	var result []Rect

	for _, rect := range rects {
		if !rect.Intersects(sub) {
			result = append(result, rect)
		} else {
			result = append(result, rect.Subtract(sub)...)
		}
	}
	return result
}

func AddRects(rects []Rect, add Rect) []Rect {
	remains := SubtractRects(rects, add)
	remains = append(remains, add)
	return remains
}

func (c Cuboid) Clip() Cuboid {
	return Cuboid{
		x: c.x.Intersect(Interval{-50, 50}),
		y: c.y.Intersect(Interval{-50, 50}),
		z: c.z.Intersect(Interval{-50, 50}),
	}
}

func (c Cuboid) InClipArea() bool {
	iv := Interval{-50, 50}
	return c.x.Intersects(iv) && c.y.Intersects(iv) && c.z.Intersects(iv)
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

func ParseLine(line string) Line {
	var c Cuboid
	var onOff string
	_, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
		&onOff, &c.x.min, &c.x.max, &c.y.min, &c.y.max, &c.z.min, &c.z.max,
	)
	if err != nil {
		panic(err)
	}
	return Line{
		c:    c,
		isOn: onOff == "on",
	}
}

type Line struct {
	c    Cuboid
	isOn bool
}

func (a Interval) ExtentWith(b Interval) Interval {
	return Interval{
		min: min(a.min, b.min),
		max: max(a.max, b.max),
	}
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	lines := util.Map(linesText, ParseLine)

	extent := Cuboid{}
	for _, line := range lines {
		if line.isOn {
			cuboid := line.c
			extent.x = extent.x.ExtentWith(cuboid.x)
			extent.y = extent.y.ExtentWith(cuboid.y)
			extent.z = extent.z.ExtentWith(cuboid.z)
		}
	}

	var disjointRects []Rect

	start := time.Now()

	// bigOnes := util.Filter(lines, func(line Line) bool { return !line.c.InClipArea() })

	num := 0
	for z := extent.z.min; z <= extent.z.max; z++ {
		for _, line := range lines {
			c := line.c
			if c.z.min <= z && z <= c.z.max {
				r := Rect{
					x: line.c.x,
					y: line.c.y,
				}
				if line.isOn {
					disjointRects = AddRects(disjointRects, r)
				} else {
					disjointRects = SubtractRects(disjointRects, r)
				}
				disjointRects = util.Filter(disjointRects, func(r Rect) bool { return r.Area() > 0 })

				// fmt.Printf("%d %s -> %d disjoint rects\n",
				// 	i, linesText[i], len(disjointRects),
				// )
			}
		}

		n := 0
		for _, rect := range disjointRects {
			n += rect.Area()
		}

		num += n
		nDone := z - extent.z.min + 1
		if nDone%1000 == 0 {
			fmt.Printf("y=%d done %d after %v, %d rects area=%d\n", z, nDone, time.Since(start), len(disjointRects), n)
		}
	}

	fmt.Printf("Num cells on: %d\n", num)
}
