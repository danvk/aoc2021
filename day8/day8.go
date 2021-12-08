package main

import (
	"aoc/util"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// 0: abc efg 6
// 1:   c  f  2 *
// 2: a cde g 5
// 3: a cd fg 5
// 4:  bcd f  4 *
// 5: ab d fg 5
// 6: ab defg 6
// 7: a c  f  3 *
// 8: abcdefg 7 *
// 9: abcd fg 6

var PATTERNS_STR = [...]string{
	0: "abcefg",
	1: "cf",
	2: "acdeg",
	3: "acdfg",
	4: "bcdf",
	5: "abdfg",
	6: "abdefg",
	7: "acf",
	8: "abcdefg",
	9: "abcdfg",
}

var PATTERNS [][]Digit

// Why can't I call this Init()?
func InitPatterns() {
	PATTERNS = make([][]Digit, 10)
	for num, s := range PATTERNS_STR {
		digits := make([]Digit, len(s))
		for i, b := range s {
			digits[i] = Digit(b)
		}
		PATTERNS[num] = digits
	}
}

// acedgfb cdfbe gcdfa fbcad dab cefabd cdfgeb eafb cagedb ab
// 8                         7                 4           1
//                           acf               bcdf        cf
// | cdfeb fcadb cdfeb cdbaf

// 1
// a -> {c, f}
// b -> {c, f}
//
// 4
// a -> {b, c, d, f}
// e -> {b, c, d, f}
// f -> {b, c, d, f}
// b -> {b, c, d, f}
//
// 7
// d -> {a, c, f}  d -> a
// a -> {a, c, f}
// b -> {a, c, f}

type ScrambledDigit byte
type Digit byte

type SignalPattern struct {
	signals       []ScrambledDigit
	candidateNums []int
}

func MakeSignalPattern(text string) (result SignalPattern) {
	for _, code := range text {
		result.signals = append(result.signals, ScrambledDigit(code))
	}
	n := len(result.signals)
	if n == 2 {
		result.candidateNums = []int{1}
	} else if n == 3 {
		result.candidateNums = []int{7}
	} else if n == 4 {
		result.candidateNums = []int{4}
	} else if n == 5 {
		result.candidateNums = []int{2, 3, 5}
	} else if n == 6 {
		result.candidateNums = []int{0, 6, 9}
	} else if n == 7 {
		result.candidateNums = []int{8}
	}
	return result
}

func Intersect(a []Digit, b []Digit) []Digit {
	var ok [10]int
	for _, d := range a {
		c := d - 'a'
		ok[c] += 1
	}
	for _, d := range b {
		c := d - 'a'
		ok[c] += 1
	}

	out := []Digit{}
	for i, v := range ok {
		if v == 2 {
			out = append(out, Digit('a'+i))
		}
	}
	return out
}

type Scramble struct {
	signals  []SignalPattern
	mappings map[ScrambledDigit][]Digit
}

func NewScramble() (result Scramble) {
	result.mappings = make(map[ScrambledDigit][]Digit)
	for _, char := range "abcdefg" {
		result.mappings[ScrambledDigit(char)] = []Digit{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	}
	return result
}

func (s *Scramble) PrintCandidates() {
	for _, sd := range []ScrambledDigit{'a', 'b', 'c', 'd', 'e', 'f', 'g'} {
		fmt.Printf("%s: ", string(sd))
		for _, d := range s.mappings[sd] {
			fmt.Printf("%s", string(d))
		}
		fmt.Printf("\n")
	}
}

// Narrow down the possible mappings based on how many are displayed
func (s *Scramble) NarrowByLength() {
	for _, signal := range s.signals {
		if len(signal.candidateNums) != 1 {
			continue
		}

		num := signal.candidateNums[0]
		clear_segs := PATTERNS[num]
		for _, scramble := range signal.signals {
			s.mappings[scramble] = Intersect(s.mappings[scramble], clear_segs)
		}
	}
}

// Narrow down the possible mappings based on the signals
func (s *Scramble) NarrowMappings() {

}

func Without[T comparable](xs []T, val T) []T {
	return util.Filter(xs, func(x T) bool { return x != val })
}

func KeyForValue[K comparable, V comparable](m map[K]V, val V) (K, bool) {
	var first K
	for k, v := range m {
		first = k
		if v == val {
			return k, true
		}
	}
	return first, false
}

// The 1 uses only C&F and the 7 uses A, C and F.
// So with 1 & 7, we know what A is.
func (s *Scramble) FindTheA() {
	var ones []ScrambledDigit
	for _, signal := range s.signals {
		if !reflect.DeepEqual(signal.candidateNums, []int{1}) {
			continue
		}
		if len(signal.signals) != 2 {
			panic(signal)
		}
		// Must be the 1
		ones = signal.signals
	}

	fmt.Printf("Ones: %#v\n", ones)
	var sevens []ScrambledDigit
	for _, signal := range s.signals {
		if !reflect.DeepEqual(signal.candidateNums, []int{7}) {
			continue
		}
		if len(signal.signals) != 3 {
			panic(signal)
		}
		// Must be the 7
		sevens = signal.signals
	}
	fmt.Printf("Sevens: %#v\n", sevens)

	for scramble, clear := range s.mappings {
		if scramble != ones[0] && scramble != ones[1] {
			s.mappings[scramble] = Without(Without(clear, 'c'), 'f')
		}
	}

	// Figure out which one (uniquely) maps to A and remove it from the others
	var a ScrambledDigit
	for scramble, clear := range s.mappings {
		if len(clear) == 1 && clear[0] == Digit('a') {
			a = scramble
			break
		}
	}
	fmt.Printf("A: %v\n", a)
	for scramble, clear := range s.mappings {
		if scramble != a {
			s.mappings[scramble] = Without(clear, 'a')
		}
	}

	// Now B&D should be uniquely determined
	var bdPre []ScrambledDigit
	for scramble, clear := range s.mappings {
		if len(clear) == 2 && clear[0] == Digit('b') && clear[1] == Digit('d') {
			bdPre = append(bdPre, scramble)
		}
	}
	if len(bdPre) != 2 {
		panic(bdPre)
	}
	fmt.Printf("BD: %#v\n", bdPre)
	for scramble, clear := range s.mappings {
		if scramble != bdPre[0] && scramble != bdPre[1] {
			s.mappings[scramble] = Without(Without(clear, 'b'), 'd')
		}
	}
}

func main() {
	InitPatterns()
	linesText := util.ReadLines(os.Args[1])
	for _, line := range linesText {
		inOut := strings.Split(line, "|")
		input := strings.TrimSpace(inOut[0])
		output := strings.TrimSpace(inOut[1])

		parts := strings.Split(output+" "+input, " ")
		signals := util.Map(parts, MakeSignalPattern)

		scramble := NewScramble()
		scramble.signals = signals

		fmt.Printf("Line %s\n", line)
		fmt.Printf("Init:\n")
		scramble.PrintCandidates()
		scramble.NarrowByLength()
		fmt.Printf("Narrowed by length:\n")
		scramble.PrintCandidates()
		fmt.Printf("Find the seven/A:\n")
		scramble.FindTheA()
		scramble.PrintCandidates()
	}

	//
}
