package main

import (
	"aoc/util"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Step struct {
	a   int64
	b   int64
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

// 28692221693676 = too low
// 28692994993698 = incorrect
// 57592994982998 = incorrect
// 89999999999999 = too high
// 900000000000000 = too high
// 13579246899999
// 01234567890123

type State struct {
	reg   [4]int64
	input []int
}

var registers = map[string]int{
	"x": 0,
	"y": 1,
	"z": 2,
	"w": 3,
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

func GetZLoop(digits []int) int64 {
	var z int64 = 0

	// var zs []int64
	for i := 0; i < 14; i++ {
		digit := digits[i]
		step := steps[i]

		x := z % 26
		if step.div {
			z = z / 26
		}
		x += step.a
		w := int64(digit)
		if x != w {
			z = 26*z + w + step.b
		} else {
			// fmt.Printf("%2d match %d = %d a=%d b=%d \n", i, x, w, step.a, step.b)
		}
		// zs = append(zs, z)
	}

	// fmt.Printf("zs: %v\n", zs)

	return z
}

func (s State) RunStep(step Step) State {
	z := s.reg[2]
	x := z % 26
	if step.div {
		z = z / 26
	}
	x += step.a
	w := int64(s.input[0])
	s.input = s.input[1:]
	if x != w {
		z = 26*z + w + step.b
	}
	s.reg[2] = z
	return s

	/*
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
	*/
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
	/*
		digits, err := util.MapErr(os.Args[1:], strconv.Atoi)
		if err != nil {
			panic(err)
		}
		if len(digits) != 14 {
			panic(len(digits))
		}
		z := GetZLoop(digits)
		fmt.Printf("%v --> %d\n", digits, z)
		return
	*/

	var best []int
	for d0 := 1; d0 <= 9; d0++ {
		for d1 := 1; d1 <= 9; d1++ {
			for d2 := 1; d2 <= 9; d2++ {
				for d3 := 8; d3 <= 9; d3++ {
					d4 := d3 - 7
					d5 := 9
					d6 := 9
					d7 := 4
					d8 := 9
					d12 := 9
					d13 := 8
					for d9 := 1; d9 <= 9; d9++ {
						for d10 := 1; d10 <= 9; d10++ {
							for d11 := 1; d11 <= 9; d11++ {
								input := []int{
									d0, d1, d2, d3, d4, d5, d6,
									d7, d8, d9, d10, d11, d12, d13,
								}
								if GetZLoop(input) == 0 {
									fmt.Printf("got one: %v\n", best)
									best = input
								}
							}
						}
					}
				}
			}
			fmt.Printf("Completed d0,d1=%d,%d\n", d0, d1)
		}
	}
	fmt.Printf("Best: %v\n", best)

	/*
		// lowestZ := int64(-1)
		for i := 0; i < 1000000000; i++ {
			inp := RandomPlate()
			z := GetZLoop(inp)
			if z == 0 {
				fmt.Printf("%v -> %d\n", inp, z)
			}
			// if lowestZ == -1 || z < lowestZ {
			// 	lowestZ = z
			// 	fmt.Printf("Lowest z: %v -> %d\n", inp, z)
			// }
		}
	*/

	// linesText := util.ReadLines(os.Args[1])
	// num, err := strconv.ParseUint(os.Args[2], 10, 64)

	// if err != nil {
	// 	panic(err)
	// }
	/*
		for i := 0; i < 100000; i++ {
			var s1 State
			inp := RandomPlate()
			s1.input = inp
			z1 := GetZForState(s1, linesText)
			// fmt.Printf("%v --> z=%d\n", s.input, z)
			if z1 == 0 {
				fmt.Printf("\n\n   WE HAVE A WINNER! %v\n\n", inp)
			}

			// var s2 State
			// s2.input = inp
			z2 := GetZLoop(inp)
			if z1 != z2 {
				fmt.Printf("Mismatch for input: %v %d != %d\n", inp, z1, z2)
				return
			} else {
				// fmt.Printf("Match!\n")
			}
		}
	*/
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
