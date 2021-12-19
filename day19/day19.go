package main

import (
	"aoc/util"
	"fmt"
	"os"
)

type Point struct {
	x, y, z int
}

func (p Point) Add(o Point) Point {
	return Point{p.x + o.x, p.y + o.y, p.z + o.z}
}
func (p Point) Sub(other Point) Point {
	return Point{p.x - other.x, p.y - other.y, p.z - other.z}
}
func (p Point) Scale(k Point) Point {
	return Point{p.x * k.x, p.y * k.y, p.z * k.z}
}
func (p Point) Rot90Z() Point {
	return Point{-p.y, p.x, p.z}
}
func (p Point) Rot90Y() Point {
	return Point{-p.z, p.y, p.x}
}
func (p Point) Rot90X() Point {
	return Point{p.x, -p.z, p.y}
}

type Mat [][]int

func (p Point) Mult(m Mat) Point {
	return Point{
		p.x*m[0][0] + p.y*m[0][1] + p.z*m[0][2],
		p.x*m[1][0] + p.y*m[1][1] + p.z*m[1][2],
		p.x*m[2][0] + p.y*m[2][1] + p.z*m[2][2],
	}
}

var ID = Mat{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
var FLIPS = []Mat{
	{{-1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
	{{1, 0, 0}, {0, -1, 0}, {0, 0, 1}},
	{{1, 0, 0}, {0, 1, 0}, {0, 0, -1}},
}
var ROTS = []Mat{
	{{0, -1, 0}, {1, 0, 0}, {0, 0, 1}}, // rotate 90 z
	{{0, 0, -1}, {0, 1, 0}, {1, 0, 0}}, // rotate 90 y
	{{1, 0, 0}, {0, 0, -1}, {0, 1, 0}}, // rotate 90 y
}

// l[0][0], l[0][1], l[0][2]  r[0][0], r[0][1], r[0][2]
// l[1][0], l[1][1], l[1][2]  r[1][0], r[1][1], r[1][2]
// l[2][0], l[2][1], l[2][2]  r[2][0], r[2][1], r[2][2]

func (l Mat) Mult(r Mat) Mat {
	m := Mat{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m[i][j] = l[i][0]*r[0][j] + l[i][1]*r[1][j] + l[i][2]*r[2][j]
		}
	}
	return m
}
func (m Mat) Clone() Mat {
	return ID.Mult(m)
}
func (m Mat) String() string {
	return fmt.Sprintf("%d,%d,%d;%d,%d,%d;%d,%d,%d",
		m[0][1], m[0][1], m[0][2],
		m[1][1], m[1][1], m[1][2],
		m[2][1], m[2][1], m[2][2],
	)
}

func FindAllOrientations() []Mat {
	orients := map[string]Mat{}
	for nx := 0; nx < 4; nx++ {
		for ny := 0; ny < 4; ny++ {
			for nz := 0; nz < 4; nz++ {
				for fx := 0; fx <= 1; fx++ {
					for fy := 0; fy <= 1; fy++ {
						for fz := 0; fz <= 1; fz++ {
							m := ID
							for n := 0; n < nx; n++ {
								m = m.Mult(ROTS[0])
							}
							for n := 0; n < ny; n++ {
								m = m.Mult(ROTS[1])
							}
							for n := 0; n < nz; n++ {
								m = m.Mult(ROTS[2])
							}
							if fx > 0 {
								m = m.Mult(FLIPS[0])
							}
							if fy > 0 {
								m = m.Mult(FLIPS[1])
							}
							if fz > 0 {
								m = m.Mult(FLIPS[2])
							}
							orients[m.String()] = m
						}
					}
				}
			}
		}
	}
	return util.Values(orients)
}

var ROTATIONS []Mat

func init() {
	ROTATIONS = FindAllOrientations()
}

// facing positive or negative x, y, or z,
// and considering any of four directions "up" from that facing.

func NumOverlapping(as, bs []Point) int {
	aMap := make(map[Point]bool)
	for _, a := range as {
		aMap[a] = true
	}
	n := 0
	for _, b := range bs {
		if _, ok := aMap[b]; ok {
			n++
		}
	}
	return n
}

// Returns shift of bs relative to as, number of overlapping beacons
func FindBestOverlap(as, bs []Point) (Point, int) {
	// Consider that each pair might be the same
	bestOverlap := 0
	var bestShift Point
	for i, a := range as {
		for j, b := range bs {
			if j > i {
				continue
			}

			// Assume a and b are the same
			shift := b.Sub(a)
			newBs := util.Map(bs, func(p Point) Point { return p.Sub(shift) })
			overlap := NumOverlapping(as, newBs)
			if overlap > bestOverlap {
				// Ties are potentially problematic
				// And the greedy strategy might not work.
				bestOverlap = overlap
				bestShift = shift
			}
		}
	}
	return bestShift, bestOverlap
}

func main() {
	chunks := util.ReadChunks(os.Args[1])

	var scanners [][]Point
	for _, chunk := range chunks {
		var beacons []Point
		for _, line := range chunk[1:] {
			var x, y, z int
			_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
			if err != nil {
				_, err = fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
			}
			if err != nil {
				panic(err)
			}
			beacons = append(beacons, Point{x, y, z})
		}
		scanners = append(scanners, beacons)
	}
}
