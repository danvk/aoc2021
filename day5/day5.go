package main

import (
	"aoc/day1/util"
	"constraints"
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

func ordered[T constraints.Ordered](a T, b T) (T, T) {
	if a <= b {
		return a, b
	}
	return b, a
}

func (line *Line) Stroke(mat [][]int) {
	if line.start.x == line.end.x {
		// It's a column
		x := line.start.x
		y0, y1 := ordered(line.start.y, line.end.y)
		for y := y0; y <= y1; y++ {
			mat[x][y] += 1
		}
	} else if line.start.y == line.end.y {
		// It's a row
		y := line.start.y
		x0, x1 := ordered(line.start.x, line.end.x)
		for x := x0; x <= x1; x++ {
			mat[x][y] += 1
		}
	} else {
		log.Printf("Ignoring line %v", *line)
	}
}

func Xs(line Line) []int {
	return []int{line.start.x, line.end.x}
}

func Ys(line Line) []int {
	return []int{line.start.y, line.end.y}
}

func PrintMat(mat [][]int) {
	for _, row := range mat {
		for _, count := range row {
			fmt.Printf("%d", count)
		}
		fmt.Printf("\n")
	}
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	lines := util.Map(linesText, ParseLine)
	maxX := util.Max(util.FlatMap(lines, Xs))
	maxY := util.Max(util.FlatMap(lines, Ys))
	fmt.Printf("Size: %d x %d\n", maxX, maxY)

	counts := make([][]int, maxX+1)
	for x := 0; x <= maxX; x++ {
		counts[x] = make([]int, maxY+1)
	}

	for _, line := range lines {
		line.Stroke(counts)
	}

	PrintMat(counts)

	numMultiple := 0
	for _, row := range counts {
		for _, count := range row {
			if count >= 2 {
				numMultiple++
			}
		}
	}
	fmt.Printf("Num >=2: %d\n", numMultiple)
}
