package main

import (
	c "aoc/coord"
	"aoc/util"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func FoldX(dots map[c.Coord]bool, x int) map[c.Coord]bool {
	res := make(map[c.Coord]bool)
	for pos := range dots {
		if pos.X <= x {
			res[pos] = true
		} else {
			res[c.New(2*x-pos.X, pos.Y)] = true
		}
	}
	return res
}

func FoldY(dots map[c.Coord]bool, y int) map[c.Coord]bool {
	res := make(map[c.Coord]bool)
	for pos := range dots {
		if pos.Y <= y {
			res[pos] = true
		} else {
			res[c.New(pos.X, 2*y-pos.Y)] = true
		}
	}
	return res
}

func PrintDots(dots map[c.Coord]bool) {
	c.PrintGrid(dots, ".", func(v bool) string { return "#" })
}

func main() {
	linesText := util.ReadChunks(os.Args[1])
	if len(linesText) != 2 {
		panic(linesText)
	}

	dotsText := linesText[0]
	foldsText := linesText[1]

	dots := make(map[c.Coord]bool)
	for _, line := range dotsText {
		parts, ok := util.MapErr(strings.Split(line, ","), strconv.Atoi)
		if ok != nil || len(parts) != 2 {
			panic(parts)
		}
		x := parts[0]
		y := parts[1]
		dots[c.Coord{X: x, Y: y}] = true
	}

	for _, fold := range foldsText {
		var dir rune
		var val int
		_, err := fmt.Sscanf(fold, "fold along %c=%d", &dir, &val)
		if err != nil {
			panic(fold)
		}

		if dir == 'x' {
			dots = FoldX(dots, val)
		} else if dir == 'y' {
			dots = FoldY(dots, val)
		} else {
			panic(fold)
		}

		fmt.Printf("After %s num dots=%d\n", fold, len(dots))
	}

	PrintDots(dots)
}
