package main

import (
	"aoc/util"
	"fmt"
	"os"
	"sort"
	"time"
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
	c.isOn = onOff == "on"
	return c
}

func (a Interval) ExtentWith(b Interval) Interval {
	return Interval{
		min: min(a.min, b.min),
		max: max(a.max, b.max),
	}
}

// Get all the distinct numbers in a list of intervals and sort them
func GetDistinctSorted(ivs []Interval) []int {
	xs := map[int]bool{}
	for _, iv := range ivs {
		xs[iv.min] = true
		xs[iv.max] = true
	}
	distinct := util.Keys(xs)
	sort.Ints(distinct)
	return distinct
}

// Returns a map from value --> index and a list of lengths
func MakeDistinctIntervals(ivs []Interval) (map[int]int, []int) {
	xs := GetDistinctSorted(ivs)

	m := map[int]int{}
	lens := []int{1}

	m[xs[0]] = 0

	for i := 1; i < len(xs); i++ {
		a := xs[i-1]
		b := xs[i]
		lens = append(lens, b-a-1, 1) // exclusive on both ends
		m[b] = len(lens) - 1
	}

	return m, lens
}

func IndexCuboids(lines []Cuboid) ([]Cuboid, []int, []int, []int) {
	xM, xL := MakeDistinctIntervals(util.Map(lines, func(line Cuboid) Interval { return line.x }))
	// fmt.Printf("xs: (%v) %v\n\n", xL, xM)

	yM, yL := MakeDistinctIntervals(util.Map(lines, func(line Cuboid) Interval { return line.y }))
	// fmt.Printf("ys: (%v) %v\n\n", yM, yL)

	zM, zL := MakeDistinctIntervals(util.Map(lines, func(line Cuboid) Interval { return line.z }))
	// fmt.Printf("ys: (%v) %v\n\n", zM, zL)

	return util.Map(lines, func(c Cuboid) Cuboid {
		return Cuboid{
			x:    Interval{min: xM[c.x.min], max: xM[c.x.max]},
			y:    Interval{min: yM[c.y.min], max: yM[c.y.max]},
			z:    Interval{min: zM[c.z.min], max: zM[c.z.max]},
			isOn: c.isOn,
		}
	}), xL, yL, zL
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	lines := util.Map(linesText, ParseLine)

	cubes, xL, yL, zL := IndexCuboids(lines)
	// fmt.Printf("Cubes: %#v\n", cubes)
	start := time.Now()
	grid := map[Coord]bool{}
	for i, c := range cubes {
		Set(c, grid)
		fmt.Printf("%3d elapsed: %v\n", i, time.Since(start))
	}

	var num int64 = 0
	for c, v := range grid {
		if v {
			num += int64(xL[c.x]) * int64(yL[c.y]) * int64(zL[c.z])
		}
	}

	fmt.Printf("Num cells on: %d\n", num)
}
