package main

import (
	"aoc/util"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Coord struct {
	x int
	y int
}

type Line struct {
	start Coord
	end   Coord
}

var LinePat, _ = regexp.Compile("(\\d+),(\\d+) -> (\\d+),(\\d+)")

func ParseLine(line string) Line {
	matches := LinePat.FindStringSubmatch(line)
	if matches == nil {
		log.Fatalf("No match for %s", line)
	}
	x0, _ := strconv.Atoi(matches[1])
	y0, _ := strconv.Atoi(matches[2])
	x1, _ := strconv.Atoi(matches[3])
	y1, _ := strconv.Atoi(matches[4])
	return Line{
		start: Coord{
			x: x0,
			y: y0,
		},
		end: Coord{
			x: x1,
			y: y1,
		},
	}
}

func (line *Line) Stroke(mat map[Coord]int) {
	if line.start.x == line.end.x {
		// It's a column
		x := line.start.x
		y0, y1 := util.Ordered(line.start.y, line.end.y)
		for y := y0; y <= y1; y++ {
			mat[Coord{x: x, y: y}] += 1
		}
	} else if line.start.y == line.end.y {
		// It's a row
		y := line.start.y
		x0, x1 := util.Ordered(line.start.x, line.end.x)
		for x := x0; x <= x1; x++ {
			mat[Coord{x: x, y: y}] += 1
		}
	} else {
		// Normalize so it's going L->R
		x0, y0 := line.start.x, line.start.y
		x1, y1 := line.end.x, line.end.y
		if x1 < x0 {
			x1, x0 = x0, x1
			y1, y0 = y0, y1
		}
		var dy int
		if y0 < y1 {
			dy = 1
		} else {
			dy = -1
		}

		// fmt.Printf("(%d, %d) - (%d, %d) dy: %d\n", x0, y0, x1, y1, dy)
		for x, y := x0, y0; x <= x1; x, y = x+1, y+dy {
			mat[Coord{x: x, y: y}] += 1
		}
	}
}

func PrintMat(mat [][]int) {
	for _, col := range mat {
		for _, count := range col {
			fmt.Printf("%d", count)
		}
		fmt.Printf("\n")
	}
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	lines := util.Map(linesText, ParseLine)
	maxX := util.Max(util.FlatMap(lines, func(x Line) []int { return []int{x.start.x, x.end.x} }))
	maxY := util.Max(util.FlatMap(lines, func(x Line) []int { return []int{x.start.y, x.end.y} }))
	fmt.Printf("Size: %d x %d\n", maxY, maxX)

	counts := map[Coord]int{}
	for _, line := range lines {
		line.Stroke(counts)
	}

	// PrintMat(counts)

	numMultiple := 0
	for _, count := range counts {
		if count >= 2 {
			numMultiple++
		}
	}
	fmt.Printf("Num >=2: %d\n", numMultiple)
}
