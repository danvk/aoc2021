package main

import (
	"aoc/graph"
	"aoc/util"
	"fmt"
	"os"
	"strings"
)

// Amber (A),
// Bronze (B),
// Copper (C), and
// Desert (D)

// Amber amphipods require 1 energy per step,
// Bronze amphipods require 10 energy,
// Copper amphipods require 100, and
// Desert ones require 1000

// Amphipods will never stop on the space immediately outside any room.
// Amphipods will never move from the hallway into a room unless that room is their destination room and that room contains no amphipods which do not also have that room as their own destination.
// Once an amphipod stops moving in the hallway, it will stay in that spot until it can move into a room.

//  01234567890
// #############
// #...........# 0
// ###B#C#B#D### 1
//   #A#D#C#A#   2
//   #########
//    A B C D
//    2 4 6 8

type Amphipod struct {
	x, y int
	kind string
}

type State struct {
	amphipods [16]Amphipod
	cost      int
}

func (s State) String() string {
	var rows [5][]string
	rows[0] = strings.Split("#...........#", "")
	rows[1] = strings.Split("###.#.#.#.###", "")
	rows[2] = strings.Split("  #.#.#.#.#", "")
	rows[3] = strings.Split("  #.#.#.#.#", "")
	rows[4] = strings.Split("  #.#.#.#.#", "")
	for _, a := range s.amphipods {
		rows[a.y][a.x+1] = a.kind
	}

	return fmt.Sprintf(
		"#############\n%s\n%s\n%s\n%s\n%s\n  ######### %d",
		strings.Join(rows[0], ""),
		strings.Join(rows[1], ""),
		strings.Join(rows[2], ""),
		strings.Join(rows[3], ""),
		strings.Join(rows[4], ""),
		s.cost,
	)
}

func ParseState(state string) State {
	lines := strings.Split(state, "\n")
	if len(lines) != 7 {
		panic(lines)
	}
	s := State{}
	i := 0
	for y := 0; y <= 4; y++ {
		line := lines[1+y]
		for x, c := range line {
			if c >= 'A' && c <= 'D' {
				s.amphipods[i].kind = string(c)
				s.amphipods[i].x = x - 1
				s.amphipods[i].y = y
				i++
			}
		}
	}
	if i != 16 {
		panic(i)
	}
	return s
}

func (s *State) AmphipodAt(x, y int) (Amphipod, bool) {
	for _, a := range s.amphipods {
		if a.x == x && a.y == y {
			return a, true
		}
	}
	return s.amphipods[0], false
}

func (s *State) IsHallwayOpen(start, stop int, selfIdx int) bool {
	x1, x2 := util.Ordered(start, stop)
	for i, a := range s.amphipods {
		if i != selfIdx && a.y == 0 && a.x >= x1 && a.x <= x2 {
			return false
		}
	}
	return true
}

func (s *State) IsComplete() bool {
	for _, a := range s.amphipods {
		if a.x != destX[a.kind] {
			return false
		}
	}
	return true
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (s State) Move(idx int, x, y int) *State {
	// value receiver, so a copy is made for us
	a := s.amphipods[idx]
	moveD := abs(x-a.x) + abs(y-a.y)
	moveCost := moveD * costs[a.kind]

	s.cost += moveCost
	s.amphipods[idx].x = x
	s.amphipods[idx].y = y
	return &s
}

// Returns a number between 0 and 27
func EncodePos(x, y int) int {
	if y == 0 {
		return x
	}
	return 11 + 4*(y-1) + ((x - 2) / 2)
}

func DecodePos(n int) (x int, y int) {
	if n <= 10 {
		return n, 0
	}
	n -= 11
	y = 1 + n>>2
	x = 2 + 2*(n&0b11)
	return
}

func EncodeKind(k string) int {
	switch k {
	case "A":
		return 1
	case "B":
		return 2
	case "C":
		return 3
	case "D":
		return 4
	}
	panic(k)
}

var kinds = []string{"A", "B", "C", "D"}

func DecodeKind(x int) string {
	return kinds[x-1]
}

func (s State) Encode() uint64 {
	var vals [27]int
	for _, a := range s.amphipods {
		vals[EncodePos(a.x, a.y)] = EncodeKind(a.kind)
	}

	var result uint64 = 0
	for _, v := range vals {
		result *= 5
		result += uint64(v)
	}
	return result
}

// #############
// #...........#
// ###A#B#C#D###
//   #A#B#C#D#
//   #A#B#C#D#
//   #A#B#C#D#
//   #########
// 27 squares
// 16 amphipods
//  2 ** 64 = 18446744073709551616L
// 27 ** 16 = 79766443076872509863361L
//  5 ** 27 = 7450580596923828125L we have a winner!

var decodeKinds = []string{"D", "D", "C", "C", "B", "B", "A", "A"}

func Decode(n uint64) (res State) {
	i := 0
	pos := 26
	for n > 0 {
		v := n % 5
		if v > 0 {
			x, y := DecodePos(pos)
			res.amphipods[i].x = x
			res.amphipods[i].y = y
			res.amphipods[i].kind = DecodeKind(int(v))
			i++
		}
		n = n / 5
		pos--
	}
	return
}

func (s *State) ToGrid() (result [5][11]string) {
	for _, a := range s.amphipods {
		result[a.y][a.x] = a.kind
	}
	return
}

func (s *State) NextStates() []*State {
	var nextStates []*State
	g := s.ToGrid()

	for i, a := range s.amphipods {
		dx := destX[a.kind]
		if a.y == 4 && a.x == dx {
			continue // already in its place, no need to move
		}

		if a.y == 0 {
			dy := -1
			for y := 4; y >= 1; y-- {
				if g[y][dx] == "" {
					dy = y
					break
				}
				if g[y][dx] != a.kind {
					break // wrong amphipod here; can't enter column
				}
			}
			if dy == -1 {
				continue
			}

			if s.IsHallwayOpen(a.x, dx, i) {
				nextStates = append(nextStates, s.Move(i, dx, dy))
			}
			continue
		}

		// Are we stuck?
		isFree := true
		for y := 1; y < a.y; y++ {
			if g[y][a.x] != "" {
				isFree = false
				break
			}
		}
		if !isFree {
			continue
		}

		// Is the stack set from us and below? If so, don't move.
		if a.x == dx {
			isSet := true
			for y := a.y + 1; y <= 4; y++ {
				if g[y][a.x] != a.kind {
					isSet = false
					break
				}
			}
			if isSet {
				continue
			}
		}

		// TODO: move directly to destination?

		// Move into hallway
		for _, x := range hallwaySpots {
			if s.IsHallwayOpen(a.x, x, i) {
				nextStates = append(nextStates, s.Move(i, x, 0))
			}
		}
	}

	return nextStates
}

var costs = map[string]int{
	"A": 1,
	"B": 10,
	"C": 100,
	"D": 1000,
}

var destX = map[string]int{
	"A": 2,
	"B": 4,
	"C": 6,
	"D": 8,
}

// Possible spots to stop in the hallway
var hallwaySpots = []int{
	0, 1, 3, 5, 7, 9, 10,
}

type Graph struct{}

func (g Graph) Neighbors(n uint64) []graph.NodeWithCost[uint64] {
	node := Decode(n)
	nextState := node.NextStates()

	return util.Map(nextState, func(s *State) graph.NodeWithCost[uint64] {
		return graph.NodeWithCost[uint64]{
			Node: s.Encode(),
			Cost: s.cost - node.cost,
		}
	})
}

func (g Graph) String(n uint64) string {
	return Decode(n).String()
}

var final = `#############
#...........#
###A#B#C#D###
  #A#B#C#D#
  #A#B#C#D#
  #A#B#C#D#
  #########`

func main() {
	text := strings.Join(util.ReadLines(os.Args[1]), "\n")
	// text := step3
	state := ParseState(text)
	fmt.Printf("Parsed state:\n%s\n", state)

	/*
		for i, s := range state.NextStates() {
			fmt.Printf("%d %d:\n%s\n\n", i, s.Encode(), s)
		}
	*/

	stop := ParseState(final).Encode()

	/*
		n := state.Encode()
		fmt.Printf("Encoded: %d\n", n)

		back := Decode(n)
		fmt.Printf("\nDecoded:\n%s\n", back)
	*/

	g := Graph{}
	cost, path := graph.Dijkstra[uint64](g, state.Encode(), func(n uint64) bool {
		return n == stop
		// s := Decode(n)
		// return s.IsComplete()
	}, 50000)

	fmt.Printf("\n\nFinal path:\n")
	fmt.Printf("Cost: %d\n\n", cost)
	for i, n := range path {
		s := Decode(n)
		fmt.Printf("%d\n%s\n\n", i, s)
	}

	// 16000 = too high
	// 20000 = too high

	/*
		states := []*State{&state}
		bestSoFar := -1
		for len(states) > 0 {
			nextStates := util.FlatMap(states, func(s *State) []*State { return s.NextStates() })
			sort.Slice(nextStates, func(i, j int) bool {
				return nextStates[i].cost < nextStates[j].cost
			})
			lowestCost := -1
			for i, state := range nextStates {
				if bestSoFar > 0 && state.cost >= bestSoFar {
					nextStates = nextStates[:i-1]
					break
				}
				if state.IsComplete() {
					lowestCost = state.cost
					if i == 0 {
						nextStates = nil
					} else {
						nextStates = nextStates[:i-1]
					}
					bestSoFar = lowestCost
					break
				}
			}

			if lowestCost > 0 {
				fmt.Printf("Best this round: %d, so far: %d\n", lowestCost, bestSoFar)
			}

			fmt.Printf("%d states\n", len(nextStates))
			// if len(nextStates) > 0 {
			// 	fmt.Printf("%s\n\n", nextStates[0])
			// }

			states = nextStates

			// for i, s := range states {
			// 	fmt.Printf("%d\n%s\n\n", i, s)
			// }
			// break
		}
	*/
}
