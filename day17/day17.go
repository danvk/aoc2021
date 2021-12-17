package main

import (
	"aoc/util"
	"fmt"
	"os"
)

// Notes:
// - There's some min x velocity below which you won't make it to the target area
// - There's some max x velocity above which you will overshoot

// What's the final x you get to with a given initial velocity?
func finalX(vx int) int {
	x := 0
	for vx > 0 {
		x += vx
		vx -= 1
	}
	return x
}

func maxHeight(vx, vy int, xMin, xMax, yMin, yMax int) int {
	maxY := 0
	x, y := 0, 0
	for x <= xMax {
		// fmt.Printf("pos: %d,%d v: %d,%d\n", x, y, vx, vy)

		x += vx
		y += vy
		if y > maxY {
			maxY = y
		}
		if x >= xMin && x <= xMax && y >= yMin && y <= yMax {
			return maxY
		}
		if (vx < 0 && x < xMin) || y < yMin {
			// fmt.Printf("38: x, y=%d,%d, vx=%d, xMin,yMin=%d,%d\n", x, y, vx, xMin, yMin)
			return -1
		}
		if vx < 0 {
			vx += 1
		} else if vx > 0 {
			vx -= 1
		}
		vy -= 1
	}
	return -1
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	if len(linesText) != 1 {
		panic(linesText)
	}

	// target area: x=153..199, y=-114..-75
	var xMin, xMax, yMin, yMax int
	_, err := fmt.Sscanf(linesText[0], "target area: x=%d..%d, y=%d..%d", &xMin, &xMax, &yMin, &yMax)
	if err != nil {
		panic(err)
	}

	fmt.Println(linesText[0])
	var minVx int
	for vx := 1; vx <= xMin; vx++ {
		fmt.Printf("vx: %3d final x: %4d\n", vx, finalX(vx))
		if finalX(vx) > xMin {
			minVx = vx
			break
		}
	}

	// m := maxHeight(6, 3, xMin, xMax, yMin, yMax)
	// fmt.Printf("%d\n", m)

	best := 0
	for vx := minVx; vx < xMin; vx++ {
		// what's the right upper bound here?
		for vy := 1; vy <= 1000; vy++ {
			// fmt.Printf("%d, %d\n", vx, vy)
			m := maxHeight(vx, vy, xMin, xMax, yMin, yMax)
			if m >= 0 {
				fmt.Printf("v: %3d, %3d max height: %4d\n", vx, vy, m)
				if m > best {
					best = m
					fmt.Printf("  best so far!\n")
				}
			}
		}
	}

	fmt.Printf("Max height: %d\n", best)
}
