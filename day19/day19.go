package main

import (
	"aoc/util"
	"fmt"
	"os"
)

type Point struct {
	x, y, z int
}

func Add(a, b Point) Point {
	return Point{a.x + b.x, a.y + b.y, a.z + b.z}
}
func (p Point) Sub(other Point) Point {
	return Point{p.x - other.x, p.y - other.y, p.z - other.z}
}

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
