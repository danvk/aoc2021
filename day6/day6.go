package main

import (
	"aoc/util"
	"fmt"
	"os"
)

type Lanternfish struct {
	timer int
}

func Advance(school *[]Lanternfish) {
	newFish := []Lanternfish{}
	for _, fish := range *school {
		if fish.timer == 0 {
			newFish = append(newFish, Lanternfish{8})
			fish.timer = 6
		} else {
			fish.timer -= 1
		}
	}
	if len(newFish) > 0 {
		*school = append(*school, newFish...)
	}
}

func PrintSchool(school []Lanternfish) {
	fmt.Printf("%v\n", school)
}

func AdvanceOne(timer int) []int {
	if timer == 0 {
		return []int{6, 8}
	}
	return []int{timer - 1}
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	school := util.ParseLineAsNums(linesText[0], ",", false)
	// school := util.Map(initCounts, func(n int) Lanternfish { return Lanternfish{timer: n} })

	for day := 1; day <= 80; day++ {
		school = util.FlatMap(school, AdvanceOne)

		fmt.Printf("Day %d, size: %d\n", day, len(school))
		// PrintSchool(school)
	}
}
