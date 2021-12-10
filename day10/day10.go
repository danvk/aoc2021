package main

import (
	"aoc/util"
	"fmt"
	"os"
	"sort"
)

var openToClose = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}

var scores = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

// Parse the line and return true, 0, completion score;
// otherwise return false, first illegal char, 0
func parseLine(line string) (bool, rune, int) {
	stack := []rune{}
	for _, c := range line {
		if close, ok := openToClose[c]; ok {
			stack = append(stack, close)
		} else if c == stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		} else {
			return false, c, 0
		}
	}

	score := 0
	for i := len(stack) - 1; i >= 0; i-- {
		score *= 5
		points, ok := scores[stack[i]]
		if !ok {
			panic(stack[i])
		}
		score += points
	}

	return true, 0, score
}

var charToPoints = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	pointsTally := 0
	scores := []int{}
	for _, line := range linesText {
		ok, badchar, score := parseLine(line)
		if ok {
			fmt.Printf("%s: ok! score: %d\n", line, score)
			scores = append(scores, score)
		} else {
			fmt.Printf("%s: illegal %c\n", line, badchar)
			pointsTally += charToPoints[badchar]
		}
	}

	fmt.Printf("Invalid lines score: %d\n", pointsTally)

	sort.Ints(scores)
	score := scores[(len(scores)-1)/2]
	fmt.Printf("Completion score: %d\n", score)
}
