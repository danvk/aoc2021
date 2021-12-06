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

	counts := map[int]int{}
	for _, timer := range school {
		counts[timer] += 1
	}

	for day := 1; day <= 256; day++ {
		// school = util.FlatMap(school, AdvanceOne)
		nextCounts := map[int]int{}
		for timer, count := range counts {
			if timer == 0 {
				nextCounts[6] += count
				nextCounts[8] += count
			} else {
				nextCounts[timer-1] += count
			}
		}
		counts = nextCounts

		schoolSize := 0
		for _, count := range counts {
			schoolSize += count
		}

		fmt.Printf("Day %d, size: %d\n", day, schoolSize)
		// PrintSchool(school)
	}
}
