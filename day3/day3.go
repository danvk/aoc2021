package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func mostCommonBitAtPos(nums *[]string, pos int) int {
	counts := []int{0, 0}
	for _, num := range *nums {
		d := num[pos] - '0'
		counts[d] += 1
	}

	zeros := counts[0]
	ones := counts[1]
	if zeros > ones {
		return 0
	} else if ones > zeros {
		return 1
	}
	return -1 // tie
}

func filterByBit(nums *[]string, pos int, bit int) []string {
	result := make([]string, 0)
	var c byte
	if bit == 0 {
		c = '0'
	} else {
		c = '1'
	}
	for _, num := range *nums {
		if num[pos] == c {
			result = append(result, num)
		}
	}
	return result
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nums := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		nums = append(nums, line)
	}

	oxygenRatingSet := nums
	for pos := range nums[0] {
		// fmt.Printf("nums[0]: %s, pos: %d", nums[0], pos)
		mostCommon := mostCommonBitAtPos(&oxygenRatingSet, pos)
		if mostCommon == -1 {
			mostCommon = 1
		}
		oxygenRatingSet = filterByBit(&oxygenRatingSet, pos, mostCommon)
		if len(oxygenRatingSet) == 1 {
			break
		}
	}
	oxygenSetRating := oxygenRatingSet[0]

	co2ScrubberSet := nums
	for pos := range nums[0] {
		mostCommon := mostCommonBitAtPos(&co2ScrubberSet, pos)
		leastCommon := 1 - mostCommon
		if mostCommon == -1 {
			leastCommon = 0
		}
		co2ScrubberSet = filterByBit(&co2ScrubberSet, pos, leastCommon)
		if len(co2ScrubberSet) == 1 {
			break
		}
	}
	co2ScrubberRating := co2ScrubberSet[0]

	fmt.Printf("Oxygen Set Rating: %s\n", oxygenSetRating)
	fmt.Printf("CO2 Scrubber Rating: %s\n", co2ScrubberRating)

	od, err := strconv.ParseInt(oxygenSetRating, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	cd, err := strconv.ParseInt(co2ScrubberRating, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d * %d = %d\n", od, cd, od*cd)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
