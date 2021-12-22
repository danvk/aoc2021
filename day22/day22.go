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
	return fmt.Sprintf("[%d, %d)", iv.min, iv.max)
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

func Set(c Cuboid, grid *[1000][1000][1000]bool) {
	for x := c.x.min; x < c.x.max; x++ {
		for y := c.y.min; y < c.y.max; y++ {
			for z := c.z.min; z < c.z.max; z++ {
				grid[x][y][z] = c.isOn
			}
		}
	}
}

func (c Cuboid) String() string {
	onOff := "off"
	if c.isOn {
		onOff = " on"
	}
	return fmt.Sprintf("%s x:%s y:%s z:%s", onOff, c.x, c.y, c.z)
}

func ParseLine(line string) Cuboid {
	var c Cuboid
	var onOff string
	_, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
		&onOff, &c.x.min, &c.x.max, &c.y.min, &c.y.max, &c.z.min, &c.z.max,
	)
	// Make the intervals half-open
	c.x.max += 1
	c.y.max += 1
	c.z.max += 1
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
	fmt.Printf(" %#v\n", xs)

	m := map[int]int{}
	lens := []int{}

	m[xs[0]] = 0

	// [-10, 10] x=-10..9
	// [-5, 15] x=-5..14
	// -10, -5, 10, 15
	// [-10, -5] x=-10..-6
	// [-5, 10]  x=-5..9
	// [10, 15]  x=10..14

	for i := 0; i < len(xs)-1; i++ {
		a := xs[i]
		b := xs[i+1]
		lens = append(lens, b-a)
		m[b] = len(lens)
	}

	return m, lens
}

func IndexCuboids(lines []Cuboid) ([]Cuboid, []int, []int, []int) {
	xM, xL := MakeDistinctIntervals(util.Map(lines, func(line Cuboid) Interval { return line.x }))
	fmt.Printf("xs: (%v) %v\n\n", xM, xL)

	yM, yL := MakeDistinctIntervals(util.Map(lines, func(line Cuboid) Interval { return line.y }))
	fmt.Printf("ys: (%v) %v\n\n", yM, yL)

	zM, zL := MakeDistinctIntervals(util.Map(lines, func(line Cuboid) Interval { return line.z }))
	fmt.Printf("ys: (%v) %v\n\n", zM, zL)

	fmt.Printf("Distinct xs: %d, ys: %d, zs: %d\n", len(xL), len(yL), len(zL))

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
	fmt.Printf("cuboids: %v\n", lines)

	cubes, xL, yL, zL := IndexCuboids(lines)
	fmt.Printf("cubes: %v\n", cubes)
	// fmt.Printf("Cubes: %#v\n", cubes)
	start := time.Now()
	// grid := map[Coord]bool{}
	if len(xL) > 1000 {
		panic(len(xL))
	}
	if len(yL) > 1000 {
		panic(len(yL))
	}
	if len(zL) > 1000 {
		panic(len(zL))
	}
	var grid [1000][1000][1000]bool
	for i, c := range cubes {
		Set(c, &grid)
		fmt.Printf("%3d elapsed: %v\n", i, time.Since(start))
	}

	var num int64 = 0
	for x := 0; x < len(xL); x++ {
		for y := 0; y < len(yL); y++ {
			for z := 0; z < len(zL); z++ {
				if grid[x][y][z] {
					num += int64(xL[x]) * int64(yL[y]) * int64(zL[z])
				}
			}
		}
	}

	fmt.Printf("Num cells on: %d\n", num)
}
