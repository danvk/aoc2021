package main

import (
	"aoc/util"
	"fmt"
	"log"
	"os"
	"strconv"
)

func calcSums(nums []int, n int) []int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += nums[i]
	}

	sums := []int{sum}
	for i := n; i < len(nums); i++ {
		sums = append(sums, sums[len(sums)-1]+nums[i]-nums[i-n])
	}

	return sums
}

func main() {
	nums, err := util.MapErr(util.ReadLines(os.Args[1]), strconv.Atoi)
	if err != nil {
		log.Fatal(err)
	}

	sums := calcSums(nums, 3)
	for _, sum := range sums {
		fmt.Printf("%d\n", sum)
	}

	last := -1
	numInc := 0
	for _, n := range sums {
		if n > last && last != -1 {
			numInc++
		}
		last = n
	}

	fmt.Printf("Num increased: %d\n", numInc)
}
