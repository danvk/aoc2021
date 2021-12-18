package main

import (
	"aoc/util"
	"fmt"
	"os"
)

func main() {
	linesText := util.ReadLines(os.Args[1])
	for _, line := range linesText {
		fmt.Printf("%s\n", line)
	}
}
