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
func (p Point) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.x, p.y, p.z)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (a Point) Manhattan(b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
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
	return fmt.Sprintf("%d,%d,%d; %d,%d,%d; %d,%d,%d",
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2],
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
							/*
								if fx > 0 {
									m = m.Mult(FLIPS[0])
								}
								if fy > 0 {
									m = m.Mult(FLIPS[1])
								}
								if fz > 0 {
									m = m.Mult(FLIPS[2])
								}
							*/
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

func NumOverlapping(aMap map[Point]bool, bs []Point) int {
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
	aMap := make(map[Point]bool)
	for _, a := range as {
		aMap[a] = true
	}

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
			overlap := NumOverlapping(aMap, newBs)
			if overlap > bestOverlap {
				// Ties are potentially problematic
				// And the greedy strategy might not work.
				bestOverlap = overlap
				bestShift = shift
				if bestOverlap >= 12 {
					return bestShift, bestOverlap
				}
			}
		}
	}
	return bestShift, bestOverlap
}

func DeDupe(points []Point) []Point {
	m := map[string]Point{}
	for _, p := range points {
		s := p.String()
		_, ok := m[s]
		if !ok {
			m[s] = p
		}
	}
	return util.Values(m)
}

func FindBestRotatedOverlap(as, bs []Point) (Mat, Point, int) {
	bestOverlap := 0
	var bestShift Point
	bestMat := ID

	for _, m := range ROTATIONS {
		rotBs := util.Map(bs, func(p Point) Point { return p.Mult(m) })
		shift, overlap := FindBestOverlap(as, rotBs)
		if overlap > bestOverlap {
			bestOverlap = overlap
			bestShift = shift
			bestMat = m
			if bestOverlap >= 12 {
				break
			}
		}
	}

	return bestMat, bestShift, bestOverlap
}

func main() {
	chunks := util.ReadChunks(os.Args[1])

	var scanners [][]Point
	for _, chunk := range chunks {
		var beacons []Point
		for _, line := range chunk[1:] {
			var x, y, z int
			_, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
			if err != nil {
				panic(err)
			}
			beacons = append(beacons, Point{x, y, z})
		}
		scanners = append(scanners, beacons)
	}

	aligned := scanners[0]
	remaining := scanners[1:]

	shifts := []Point{{0, 0, 0}}

	for len(remaining) > 0 {
		gotOne := false
		for i, next := range remaining {
			mat, shift, overlap := FindBestRotatedOverlap(aligned, next)
			if overlap >= 12 {
				shifted := util.Map(next, func(p Point) Point { return p.Mult(mat).Sub(shift) })
				// XXX more efficient to de-dupe here
				aligned = DeDupe(append(aligned, shifted...))
				remaining = append(remaining[:i], remaining[i+1:]...)
				fmt.Printf("Got one! %d/%d\n", len(aligned), len(remaining))
				shifts = append(shifts, shift)
				gotOne = true
				break
			}
		}

		if !gotOne {
			panic("Unable to align!")
		}
	}

	allPoints := map[string]Point{}
	for _, point := range aligned {
		// for _, point := range a {
		allPoints[point.String()] = point
		// }
	}

	fmt.Printf("Total beacons: %d\n", len(allPoints))

	longest := 0
	for i, a := range shifts {
		for j, b := range shifts {
			if j >= i {
				break
			}
			d := a.Manhattan(b)
			if d > longest {
				longest = d
			}
		}
	}
	fmt.Printf("Longest distance: %d\n", longest)

	/*
		for i, a := range scanners {
			for j, b := range scanners {
				if i >= j {
					continue
				}

				mat, shift, overlap := FindBestRotatedOverlap(a, b)
				if overlap >= 12 {
					fmt.Printf("%d -> %d: %d %s %s\n", i, j, overlap, shift, mat)
					nextBs := util.Map(b, func(p Point) Point { return p.Mult(mat).Sub(shift) })
					fmt.Printf("  %d\n", NumOverlapping(a, nextBs))
				}
			}
		}
	*/
}
