package main

import (
	"aoc/util"
	"fmt"
	"os"
)

type Move struct {
	dir  string
	dist int
}

func parseLine(line string) (Move, error) {
	var move Move
	_, err := fmt.Sscanf(line, "%s %d", &move.dir, &move.dist)
	return move, err
}

func main() {
	moves, err := util.MapErr(util.ReadLines(os.Args[1]), parseLine)
	if err != nil {
		panic(err)
	}

	pos := 0
	depth := 0
	aim := 0
	for _, move := range moves {
		dist := move.dist
		switch move.dir {
		case "forward":
			pos += dist
			depth += aim * dist
		case "down":
			aim += dist
		case "up":
			aim -= dist
		}
	}

	fmt.Printf("pos: %d, depth: %d, product: %d\n", pos, depth, pos*depth)
}
