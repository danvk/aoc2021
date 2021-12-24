package main

import (
	"aoc/util"
	"constraints"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Step struct {
	a   int
	b   int
	div bool
}

func (s Step) String() string {
	return fmt.Sprintf("a: %-2d b: %-2d %v", s.a, s.b, s.div)
}

var steps = []Step{
	{a: 14, b: 12},             // d0
	{a: 11, b: 8},              // d1
	{a: 11, b: 7},              // d2
	{a: 14, b: 4},              // d3
	{a: -11, b: 4, div: true},  // d4
	{a: 12, b: 1},              // d5
	{a: -1, b: 10, div: true},  // d6
	{a: 10, b: 8},              // d7
	{a: -3, b: 12, div: true},  // d8
	{a: -4, b: 10, div: true},  // d9
	{a: -13, b: 15, div: true}, // d10
	{a: -8, b: 4, div: true},   // d11
	{a: 13, b: 10},             // d12
	{a: -11, b: 9, div: true},  // d13
}

// 89999999999999 = too high
// 900000000000000 = too high
// 13579246899999
// 01234567890123

type State struct {
	reg   [4]int64
	input []int
}

type Interval struct {
	min, max int64
}

func (iv Interval) String() string {
	if iv.min == iv.max {
		return fmt.Sprintf("%d", iv.min)
	}
	return fmt.Sprintf("%d..%d", iv.min, iv.max)
}

func (iv *Interval) Add(other Interval) {
	iv.min += other.min
	iv.max += other.max
}

func min[T constraints.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func (iv *Interval) Mul(other Interval) {
	a, b, c, d := iv.min*other.min, iv.min*other.max, iv.max*other.min, iv.max*other.max
	iv.min = min(a, min(b, min(c, d)))
	iv.max = max(a, max(b, max(c, d)))
}

type RangeState struct {
	reg [4]Interval
}

var registers = map[string]int{
	"x": 0,
	"y": 1,
	"z": 2,
	"w": 3,
}

func (s *RangeState) GetValue(str string) Interval {
	n, ok := registers[str]
	if ok {
		return s.reg[n]
	}
	return Interval{1, 9}
}

func (s *State) GetValue(str string) int64 {
	n, ok := registers[str]
	if ok {
		return s.reg[n]
	}
	// must be a number
	n, err := strconv.Atoi(str)
	if err != nil {
		panic(str)
	}
	return int64(n)
}

func RunInstruction(s State, instr string) State {
	if len(instr) == 0 {
		return s // blank lines are OK
	}

	parts := strings.Split(instr, " ")
	if len(parts) < 2 {
		panic(instr)
	}
	cmd := parts[0]

	n, ok := registers[parts[1]]
	if !ok {
		panic(instr)
	}

	if cmd == "inp" {
		s.reg[n] = int64(s.input[0])
		s.input = s.input[1:]
		return s
	}
	if len(parts) < 3 {
		panic(instr)
	}

	v := s.GetValue(parts[2])
	switch cmd {
	case "add":
		s.reg[n] += v
	case "mul":
		s.reg[n] *= v
	case "mod":
		if v < 0 {
			panic(instr)
		}
		s.reg[n] = s.reg[n] % v
	case "div":
		if v == 0 {
			panic(instr)
		}
		s.reg[n] /= v
	case "eql":
		if s.reg[n] == v {
			s.reg[n] = 1
		} else {
			s.reg[n] = 0
		}
	}
	return s
}

func (s State) String() string {
	return fmt.Sprintf(
		"x=%d y=%d z=%d w=%d digits=%v",
		s.reg[0], s.reg[1], s.reg[2], s.reg[3],
		s.input,
	)
}

func (s State) RunStep(step Step) State {
	x, y, z := s.reg[0], s.reg[1], s.reg[2]
	w := int64(s.input[0])
	s.input = s.input[1:]

	x = z % 26
	if step.div {
		z = z / 26
	}
	x += int64(step.a)
	if x == w {
		x = 0
	} else {
		z *= 26
		x = 1
	}
	y = w + int64(step.b)
	y = y * x
	z += y

	s.reg[0] = x
	s.reg[1] = y
	s.reg[2] = z
	s.reg[3] = w
	return s
}

func (s RangeState) RunInstruction(instr string) RangeState {
	if len(instr) == 0 {
		return s // blank lines are OK
	}

	parts := strings.Split(instr, " ")
	if len(parts) < 2 {
		panic(instr)
	}
	cmd := parts[0]

	n, ok := registers[parts[1]]
	if !ok {
		panic(instr)
	}

	if cmd == "inp" {
		s.reg[n] = Interval{1, 9}
		// (Consume input)
		return s
	}
	if len(parts) < 3 {
		panic(instr)
	}

	v := s.GetValue(parts[2])
	switch cmd {
	case "add":
		s.reg[n].Add(v)
	case "mul":
		s.reg[n].Mul(v)
	case "mod":
		if v.min != 26 || v.max != 26 {
			panic(v)
		}
		if s.reg[n].min >= 0 && s.reg[n].max <= 25 {
			// no op
		} else {
			s.reg[n].min = 0
			s.reg[n].max = 25
		}
	case "div":
		// s.reg[n] /= v
	case "eql":
		// if s.reg[n] == v {
		// 	s.reg[n] = 1
		// } else {
		// 	s.reg[n] = 0
		// }
	}
	return s
}

func GetZ(n uint64, lines []string) int64 {
	digits := strconv.FormatUint(n, 10)
	if len(digits) != 14 {
		panic(digits)
	}
	input, err := util.MapErr(strings.Split(digits, ""), strconv.Atoi)
	if err != nil {
		panic(err)
	}
	for _, i := range input {
		if i == 0 {
			return -1
		}
	}
	state := State{}
	state.input = input
	for _, line := range lines {
		state = RunInstruction(state, line)
		// fmt.Printf("%8s  # %s\n", line, state)
	}
	return int64(state.reg[2])
}

func GetZForState(state State, lines []string) int64 {
	for _, line := range lines {
		state = RunInstruction(state, line)
		// fmt.Printf("%8s  # %s\n", line, state)
	}
	return state.reg[2]
}

func GetZBySteps(state State) int64 {
	for _, step := range steps {
		state = state.RunStep(step)
	}
	return state.reg[2]
}

func RandomPlate() []int {
	ns := make([]int, 14)
	for i := 0; i < len(ns); i++ {
		ns[i] = 1 + rand.Intn(9)
	}
	return ns
}

func main() {
	linesText := util.ReadLines(os.Args[1])
	// num, err := strconv.ParseUint(os.Args[2], 10, 64)

	// if err != nil {
	// 	panic(err)
	// }
	for i := 0; i < 100000; i++ {
		var s1 State
		inp := RandomPlate()
		s1.input = inp
		z1 := GetZForState(s1, linesText)
		// fmt.Printf("%v --> z=%d\n", s.input, z)
		if z1 == 0 {
			fmt.Printf("\n\n   WE HAVE A WINNER! %v\n\n", inp)
		}

		var s2 State
		s2.input = inp
		z2 := GetZBySteps(s2)
		if z1 != z2 {
			fmt.Printf("Mismatch for input: %v %d != %d\n", inp, z1, z2)
			return
		} else {
			// fmt.Printf("Match!\n")
		}
	}

	// for n := num - 100; n <= num+10; n++ {
	// 	fmt.Printf("%d --> z=%d\n", n, GetZ(n, linesText))
	// }

	// run program z=5387967764
	// run steps   z=3958084764
	/*
		digits, err := util.MapErr(strings.Split(os.Args[2], ""), strconv.Atoi)
		if err != nil {
			panic(err)
		}
		if len(digits) != 14 {
			panic(len(digits))
		}

		fmt.Printf("Line by line:\n")
		state := State{}
		state.input = digits
		fmt.Printf("%8s  # %s\n", "(init)", state)
		for _, line := range linesText {
			if line[0:3] == "inp" {
				fmt.Println()
			}
			state = RunInstruction(state, line)
			fmt.Printf("%8s  # %s\n", line, state)
		}
		// fmt.Printf("run program z=%d\n", GetZForState(state, linesText))
		// fmt.Printf("%8s  # %s\n", "(init)", state)

		fmt.Printf("\n\nStep by step:\n")

		state = State{}
		state.input = digits
		fmt.Printf("%8s  # %s\n", "(init)", state)
		for i, step := range steps {
			state = state.RunStep(step)
			fmt.Printf("%2d: %s\n", i, state)
		}
		fmt.Printf("run steps z=%d\n", state.reg[2])
	*/
}
