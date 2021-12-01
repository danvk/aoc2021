package main

import (
	"bufio"
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
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	nums := make([]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		this, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, this)
	}

	sums := calcSums(nums, 3)
	for _, sum := range sums {
		fmt.Printf("%d\n", sum)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
