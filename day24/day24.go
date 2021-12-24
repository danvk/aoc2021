package main

import (
	"aoc/util"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// 89999999999999 = too high
// 900000000000000 = too high
// 13579246899999
// 012345678901234

type State struct {
	reg   [4]uint64
	input []int
}

var registers = map[string]int{
	"x": 0,
	"y": 1,
	"z": 2,
	"w": 3,
}

func (s *State) GetValue(str string) uint64 {
	n, ok := registers[str]
	if ok {
		return s.reg[n]
	}
	// must be a number
	n, err := strconv.Atoi(str)
	if err != nil {
		panic(str)
	}
	return uint64(n)
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
		s.reg[n] = uint64(s.input[0])
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

func GetZForState(state State, lines []string) uint64 {
	for _, line := range lines {
		state = RunInstruction(state, line)
		// fmt.Printf("%8s  # %s\n", line, state)
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

	for i := 0; i < 1000000; i++ {
		var s State
		inp := RandomPlate()
		s.input = inp
		z := GetZForState(s, linesText)
		// fmt.Printf("%v --> z=%d\n", s.input, z)
		if z == 0 {
			fmt.Printf("\n\n   WE HAVE A WINNER! %v\n\n", inp)
		}
	}

	// for n := num - 100; n <= num+10; n++ {
	// 	fmt.Printf("%d --> z=%d\n", n, GetZ(n, linesText))
	// }

	// digits, err := util.MapErr(strings.Split(os.Args[2], ""), strconv.Atoi)
	// if err != nil {
	// 	panic(err)
	// }

	// state := State{}
	// state.input = digits
	// fmt.Printf("%8s  # %s\n", "(init)", state)

	// for _, line := range linesText {
	// 	state = RunInstruction(state, line)
	// 	fmt.Printf("%8s  # %s\n", line, state)
	// }
}
