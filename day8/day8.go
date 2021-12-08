package main

import (
	"aoc/util"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// 1:   c  f  2 *
// 7: a c  f  3 *
// 4:  bcd f  4 *
// 2: a cde g 5
// 3: a cd fg 5
// 5: ab d fg 5
// 0: abc efg 6
// 6: ab defg 6
// 9: abcd fg 6
// 8: abcdefg 7 *

// e & g
// c & f
// b & d

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

// TODO: make MapErr work with bools
func GetDigitFromPattern(signal SignalPattern, code map[ScrambledDigit]Digit) (int, error) {
	digits := util.Map(signal.signals, func(d ScrambledDigit) Digit { return code[d] })
	for i, segments := range PATTERNS {
		if EqDigits(digits, segments) {
			return i, nil
		}
	}
	return -1, errors.New("invalid mapping")
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

// a: cf
// b: bd
// c: bd
// d: eg
// e: eg
// f: cf
// g: a

// a -> g
// b -> bc
// d -> bc
// c -> af
// f -> af
// e -> de
// g -> de

// Try all the mappings and find the possible one
func (s *Scramble) TryAllMappings() []int {
	im := InvertMapping(s.mappings)
	aPre := im[Digit('a')]
	bdPre := im[Digit('b')]
	egPre := im[Digit('e')]
	cfPre := im[Digit('c')]
	if len(aPre) != 1 || len(bdPre) != 2 || len(egPre) != 2 || len(cfPre) != 2 {
		panic(im)
	}

	// Try 'em all!
	a := aPre[0]
	for i := 0; i < 2; i++ {
		b := bdPre[i]
		d := bdPre[1-i]
		for j := 0; j < 2; j++ {
			e := egPre[j]
			g := egPre[1-j]
			for k := 0; k < 2; k++ {
				c := cfPre[k]
				f := cfPre[1-k]

				// Try out this mapping
				code := map[ScrambledDigit]Digit{
					a: 'a',
					b: 'b',
					c: 'c',
					d: 'd',
					e: 'e',
					f: 'f',
					g: 'g',
				}

				digits, err := util.MapErr(s.signals, func(s SignalPattern) (int, error) {
					return GetDigitFromPattern(s, code)
				})
				if err == nil {
					fmt.Printf("Decoding success!\n  Digits = %#v\n", digits)
					PrintCode(code)
					return digits
				}
				fmt.Printf(" nope, not that one!\n")
			}
		}
	}
	panic("No valid decoding found")
}

// a->c
// f->d
// g->e
// b->f
// c->g
// d->a
// e->b

// cdfbe -> gadfb (5)

// acedgfb=8
// cdfbe: 5
// gcdfa: 2
// fbcad: 3
// dab: 7 (a, c, f) yes
// cefabd: 9
// cdfgeb: 6
// eafb: 4 (b, c, d, f) yes
// cagedb: 0
// ab: 1 (c, f) yes

func Without[T comparable](xs []T, val T) []T {
	return util.Filter(xs, func(x T) bool { return x != val })
}

func InvertMapping(m map[ScrambledDigit][]Digit) map[Digit][]ScrambledDigit {
	result := make(map[Digit][]ScrambledDigit)
	for scramble, clears := range m {
		for _, clear := range clears {
			result[clear] = append(result[clear], scramble)
		}
	}
	return result
}

// Are the two sets of digits equal, ignoring order?
func EqDigits(a []Digit, b []Digit) bool {
	if len(a) != len(b) {
		return false
	}
	for _, ad := range a {
		found := false
		for _, bd := range b {
			if bd == ad {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
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

func PrintCode(code map[ScrambledDigit]Digit) {
	for pre, post := range code {
		fmt.Printf("%s->%s ", string(pre), string(post))
	}
	fmt.Printf("\n")
}

func main() {
	InitPatterns()
	linesText := util.ReadLines(os.Args[1])
	sum := 0
	for _, line := range linesText {
		inOut := strings.Split(line, "|")
		input := strings.TrimSpace(inOut[0])
		output := strings.TrimSpace(inOut[1])

		parts := strings.Split(input+" "+output, " ")
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

		digits := scramble.TryAllMappings()
		outputs := digits[len(digits)-4:]
		num := 1000*outputs[0] + 100*outputs[1] + 10*outputs[2] + outputs[3]
		fmt.Printf("outputs: %#v = %d\n", outputs, num)
		sum += num
	}

	fmt.Printf("Sum: %d\n", sum)
}
