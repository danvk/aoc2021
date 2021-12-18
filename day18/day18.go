package main

import (
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

func (p Pair) Add(other Pair) Pair {
	return Pair{
		left:  &p,
		right: &other,
	}
}

func (p Pair) IsValue() bool {
	return p.left == nil
}

// Returns new pair, exploding left, exploding right, was there an explosion?
func (p Pair) Explode(depth int) (Pair, *int, *int, bool) {
	if p.IsValue() {
		return p, nil, nil, false
	}
	if p.left.IsValue() && p.right.IsValue() && depth >= 4 {
		// this is a pair of two values; explode it!
		fmt.Printf("Explode! %s, depth=%d\n", p, depth)
		return Pair{value: 0}, &p.left.value, &p.right.value, true
	}

	// Try to explode the left
	left, expL, expR, exploded := p.left.Explode(depth + 1)
	if expR != nil && p.right.IsValue() {
		// There's an explosion! And we can handle the right part.
		return Pair{left: &left, right: &Pair{value: p.right.value + *expR}}, expL, nil, exploded
	}
	if exploded {
		right := p.right
		if expR != nil {
			newRight := right.ExplodeDownLeft(*expR)
			right = &newRight
		}
		return Pair{left: &left, right: right}, expL, expR, exploded
	}

	right, expL, expR, exploded := p.right.Explode(depth + 1)
	if expL != nil && p.left.IsValue() {
		// There's an explosion and we can handle the left part.
		return Pair{left: &Pair{value: p.left.value + *expL}, right: &right}, nil, expR, exploded
	}
	if exploded {
		left := p.left
		if expL != nil {
			newLeft := left.ExplodeDownRight(*expL)
			left = &newLeft
		}
		return Pair{left: left, right: &right}, expL, expR, exploded
	}

	// No explosions! We're unchanged.
	return p, nil, nil, false
}

func (p Pair) ExplodeDownLeft(val int) Pair {
	if p.IsValue() {
		return Pair{value: p.value + val}
	}
	newLeft := p.left.ExplodeDownLeft(val)
	return Pair{left: &newLeft, right: p.right}
}

func (p Pair) ExplodeDownRight(val int) Pair {
	if p.IsValue() {
		return Pair{value: p.value + val}
	}
	newRight := p.right.ExplodeDownRight(val)
	return Pair{left: p.left, right: &newRight}
}

func (p Pair) ReduceOnce() (Pair, bool) {
	if p.IsValue() {
		return p, false
	}
	np, _, _, exploded := p.Explode(0)
	if exploded {
		return np, true
	}
	return p, false
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
	// linesText := util.ReadLines(os.Args[1])

	// for _, line := range linesText {
	for _, line := range []string{os.Args[1]} {
		pair, rest := ParsePair(line)
		if len(rest) != 0 {
			panic(line)
		}

		np, reduced := pair.ReduceOnce()
		fmt.Printf("%s -> %s, %v\n", pair, np, reduced)
	}
}
