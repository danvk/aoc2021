package main

import (
	"aoc/util"
	"fmt"
	"os"
)

type Chunk struct {
	openChar    rune
	start, stop int
	children    []Chunk
}

// Parse the line and return true, 0; otherwise return false, first illegal char
func parseLine(line string) (bool, rune) {
	openToClose := map[rune]rune{
		'(': ')',
		'{': '}',
		'<': '>',
		'[': ']',
	}

	stack := []rune{}
	for _, c := range line {
		if close, ok := openToClose[c]; ok {
			stack = append(stack, close)
		} else if c == stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		} else {
			return false, c
		}
	}
	return true, 0
}

func main() {
	charToPoints := map[rune]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}

	linesText := util.ReadLines(os.Args[1])
	pointsTally := 0
	for _, line := range linesText {
		ok, badchar := parseLine(line)
		if ok {
			fmt.Printf("%s: ok!\n", line)
		} else {
			fmt.Printf("%s: illegal %c\n", line, badchar)
			pointsTally += charToPoints[badchar]
		}
	}
	fmt.Printf("Score: %d\n", pointsTally)
}
