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

func FuelForMove(x int) int {
	return x * (x + 1) / 2
}

func CostOfPosition(crabs []int, pos int) int {
	fuel := 0
	for _, crab := range crabs {
		fuel += FuelForMove(Abs(pos - crab))
	}
	return fuel
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	if len(linesText) != 1 {
		log.Fatalf("Expected just one line of input, got %d", len(linesText))
	}
	crabs := util.ParseLineAsNums(linesText[0], ",", false)

	min, max := util.MinMax(crabs)
	candidates := util.Seq(min, max)
	// fmt.Printf("Candidates: %#v\n", candidates)

	lowestPos, lowestFuel := util.ArgMin(candidates, func(pos int) int {
		return CostOfPosition(crabs, pos)
	})

	fmt.Printf("Lowest fuel %d @ pos %d\n", lowestFuel, lowestPos)
}
