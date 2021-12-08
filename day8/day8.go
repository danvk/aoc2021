package main

import (
	"aoc/util"
	"fmt"
	"os"
	"strings"
)

// 0: 6
// 1: 2 *
// 2: 5
// 3: 5
// 4: 4 *
// 5: 5
// 6: 6
// 7: 3 *
// 8: 7 *
// 9: 6

func main() {
	linesText := util.ReadLines(os.Args[1])
	numUniq := 0
	for _, line := range linesText {
		inOut := strings.Split(line, "|")
		// input := strings.TrimSpace(inOut[0])
		output := strings.TrimSpace(inOut[1])

		outParts := strings.Split(output, " ")
		for _, part := range outParts {
			n := len(part)
			if n == 2 || n == 4 || n == 3 || n == 7 {
				numUniq++
			}
		}
	}

	fmt.Printf("Num Unique: %d\n", numUniq)
}
