package main

import (
	"aoc/util"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Coord struct {
	X, Y int
}

func FoldX(dots map[Coord]bool, x int) map[Coord]bool {
	res := make(map[Coord]bool)
	for pos := range dots {
		if pos.X <= x {
			res[pos] = true
		} else {
			res[Coord{2*x - pos.X, pos.Y}] = true
		}
	}
	return res
}

func FoldY(dots map[Coord]bool, y int) map[Coord]bool {
	res := make(map[Coord]bool)
	for pos := range dots {
		if pos.Y <= y {
			res[pos] = true
		} else {
			res[Coord{pos.X, 2*y - pos.Y}] = true
		}
	}
	return res
}

func PrintDots(dots map[Coord]bool) {
	maxX := util.Max(util.Map(util.Keys(dots), func(pos Coord) int { return pos.X }))
	maxY := util.Max(util.Map(util.Keys(dots), func(pos Coord) int { return pos.Y }))

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			_, ok := dots[Coord{x, y}]
			if ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	linesText := util.ReadChunks(os.Args[1])
	if len(linesText) != 2 {
		panic(linesText)
	}

	dotsText := linesText[0]
	foldsText := linesText[1]

	dots := make(map[Coord]bool)
	for _, line := range dotsText {
		parts, ok := util.MapErr(strings.Split(line, ","), strconv.Atoi)
		if ok != nil || len(parts) != 2 {
			panic(parts)
		}
		x := parts[0]
		y := parts[1]
		dots[Coord{x, y}] = true
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
