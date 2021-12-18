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
	value  int
	index  int // only for values
	left   *Pair
	right  *Pair
	parent *Pair
}

func (p Pair) String() string {
	if p.left != nil {
		return fmt.Sprintf("[%s,%s]", p.left.String(), p.right.String())
	}
	return strconv.Itoa(p.value)
}

func (p *Pair) ReIndex(start int) int {
	if p.IsValue() {
		p.index = start
		return start + 1
	}
	start = p.left.ReIndex(start)
	return p.right.ReIndex(start)
}

func (p *Pair) SetParents() {
	if !p.IsValue() {
		p.left.parent = p
		p.right.parent = p
		p.left.SetParents()
		p.right.SetParents()
	}
}

func (p *Pair) Add(other *Pair) *Pair {
	sum := &Pair{
		left:  p,
		right: other,
	}
	sum.left.parent = sum
	sum.right.parent = sum
	sum.ReIndex(0)
	// fmt.Printf("  sum: %s\n", sum)
	sum.Reduce()
	return sum
}

func (p Pair) IsValue() bool {
	return p.left == nil
}

func (p *Pair) FindFirstRegular(depth int) *Pair {
	if p.IsValue() {
		if depth >= 5 {
			return p.parent
		}
		return nil
	}
	if r := p.left.FindFirstRegular(depth + 1); r != nil {
		return r
	}
	return p.right.FindFirstRegular(depth + 1)
}

// Invalidates index & parent
func (p *Pair) Split() bool {
	if p.IsValue() {
		v := p.value
		if v < 10 {
			return false
		}
		m := v % 2
		left := (v - m) / 2
		right := (v + m) / 2
		p.left = &Pair{value: left}
		p.right = &Pair{value: right}
		return true
	}
	return p.left.Split() || p.right.Split()
}

func (p *Pair) AddByIndex(updates map[int]int) {
	if p.IsValue() {
		update, ok := updates[p.index]
		if ok {
			p.value += update
		}
	} else {
		p.left.AddByIndex(updates)
		p.right.AddByIndex(updates)
	}
}

func (p *Pair) ReduceOnce() bool {
	rp := p.FindFirstRegular(0)
	if rp != nil {
		// fmt.Printf("  explode! First regular: %s\n", *rp)
		updates := map[int]int{rp.left.index - 1: rp.left.value, rp.right.index + 1: rp.right.value}
		// fmt.Printf("    updates: %#v\n", updates)
		p.AddByIndex(updates)
		// fmt.Printf("   --> %s\n", *p)
		rp.value = 0
		rp.left = nil
		rp.right = nil
		p.ReIndex(0)
		// fmt.Printf("  --> %s\n", *p)
		return true
	}

	if p.Split() {
		// fmt.Printf("  split! %s\n", *p)
		p.SetParents()
		p.ReIndex(0)
		return true
	}
	return false
}

func (p *Pair) Reduce() {
	for p.ReduceOnce() {
	}
}

func (p Pair) Magnitude() int {
	if p.IsValue() {
		return p.value
	}
	return 3*p.left.Magnitude() + 2*p.right.Magnitude()
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

func Parse(text string) *Pair {
	pair, rest := ParsePair(text)
	if len(rest) != 0 {
		panic(text)
	}
	pair.ReIndex(0)
	pair.SetParents()
	return &pair
}

func Part1(linesText []string) {
	var pair *Pair
	for _, line := range linesText {
		thisPair := Parse(line)

		if pair == nil {
			pair = thisPair
		} else {
			pair = pair.Add(thisPair)
		}
		// fmt.Printf("Reduced: %s\n", pair)
	}
	fmt.Printf("Sum: %s\n", pair)

	fmt.Printf("Final pair: %s\n", pair)
	fmt.Printf("Magnitude: %d\n", pair.Magnitude())
}

func Part2(linesText []string) {
	topMag := -1
	for i := 0; i < len(linesText); i++ {
		for j := 0; j < len(linesText); j++ {
			if i == j {
				continue
			}
			// Could implement Pair.Clone(), but this is very convenient!
			pair1 := Parse(linesText[i])
			pair2 := Parse(linesText[j])
			sum := pair1.Add(pair2)
			mag := sum.Magnitude()
			if mag > topMag {
				topMag = mag
			}
		}
	}
	fmt.Printf("Top magnitude: %d\n", topMag)
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	Part1(linesText)
	Part2(linesText)
}
