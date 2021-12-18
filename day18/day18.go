package main

import (
	"aoc/util"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// ADT would be nice for this!
type Pair struct {
	value int
	left  *Pair
	right *Pair
}

func (p Pair) String() string {
	if p.left != nil {
		return fmt.Sprintf("[%s,%s]", p.left.String(), p.right.String())
	}
	return strconv.Itoa(p.value)
}

var NumPat = regexp.MustCompile("\\d+")

// Returns the Pair and remaining, unparsed string
func ParsePair(text string) (Pair, string) {
	if text[0] == '[' {
		p1, text := ParsePair(text[1:])
		if text[0] != ',' {
			panic(text)
		}
		p2, text := ParsePair(text[1:])
		if text[0] != ']' {
			panic(text)
		}
		return Pair{left: &p1, right: &p2}, text[1:]
	}
	// XXX can I get fmt.Sscanf to return the number of bytes it consumes?
	// n, err := fmt.Sscanf(text, "%d", &value)
	pos := NumPat.FindStringIndex(text)
	if pos == nil {
		panic(text)
	}
	value, err := strconv.Atoi(text[pos[0]:pos[1]])
	if err != nil {
		panic(text)
	}
	return Pair{value: value}, text[pos[1]:]
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	for _, line := range linesText {
		pair, rest := ParsePair(line)
		if len(rest) != 0 {
			panic(line)
		}
	}
}
