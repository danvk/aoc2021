package main

import (
	"aoc/util"
	"fmt"
	"log"
	"os"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CostOfPosition(crabs []int, pos int) int {
	fuel := 0
	for _, crab := range crabs {
		fuel += Abs(pos - crab)
	}
	return fuel
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	if len(linesText) != 1 {
		log.Fatalf("Expected just one line of input, got %d", len(linesText))
	}
	crabs := util.ParseLineAsNums(linesText[0], ",", false)

	min, max := util.Min(crabs), util.Max(crabs)

	lowestFuel := -1
	lowestPos := -1
	for pos := min; pos <= max; pos++ {
		fuel := CostOfPosition(crabs, pos)
		if lowestFuel == -1 || fuel < lowestFuel {
			lowestFuel = fuel
			lowestPos = pos
		}
		fmt.Printf("Pos %d, fuel: %d\n", pos, fuel)
	}

	fmt.Printf("Lowest fuel %d @ pos %d\n", lowestFuel, lowestPos)
}
