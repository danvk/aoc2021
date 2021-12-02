package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	pos := 0
	depth := 0
	aim := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			log.Fatal(line)
		}

		dir := parts[0]
		dist, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(line)
		}

		if dir == "forward" {
			pos += dist
			depth += aim * dist
		} else if dir == "down" {
			aim += dist
		} else if dir == "up" {
			aim -= dist
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pos: %d, depth: %d, product: %d\n", pos, depth, pos*depth)
}
