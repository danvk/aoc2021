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
	counts := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(counts) == 0 {
			for range line {
				counts = append(counts, []int{0, 0})
			}
		}

		for i, char := range line {
			d, err := strconv.Atoi(string(char))
			if err != nil {
				log.Fatal(err)
			}
			counts[i][d] += 1
		}
	}

	gamma := ""
	epsilon := ""
	for i, digitCounts := range counts {
		zeros := digitCounts[0]
		ones := digitCounts[1]

		if zeros > ones {
			gamma += "0"
			epsilon += "1"
		} else if ones > zeros {
			gamma += "1"
			epsilon += "0"
		} else {
			log.Fatalf("Tie at position %d: %v", i, digitCounts)
		}
	}

	fmt.Printf("gamma: %s\n", gamma)
	fmt.Printf("espilon: %s\n", epsilon)

	gammaD, err := strconv.ParseInt(gamma, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	epsilonD, err := strconv.ParseInt(epsilon, 2, 64)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d * %d = %d\n", gammaD, epsilonD, gammaD*epsilonD)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
