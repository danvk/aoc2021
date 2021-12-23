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
		"#############\n%s\n%s\n%s\n%s\n%s\n  #########",
		strings.Join(rows[0], ""),
		strings.Join(rows[1], ""),
		strings.Join(rows[2], ""),
		strings.Join(rows[3], ""),
		strings.Join(rows[4], ""),
		// s.cost,
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

func (s *State) IsHallwayOpen(start, stop int, selfIdx int) bool {
	x1, x2 := util.Ordered(start, stop)
	for i, a := range s.amphipods {
		if i != selfIdx && a.y == 0 && a.x >= x1 && a.x <= x2 {
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

func (g Graph) Neighbors(n string) []graph.NodeWithCost[string] {
	node := ParseState(n)
	nextState := node.NextStates()

	return util.Map(nextState, func(s *State) graph.NodeWithCost[string] {
		return graph.NodeWithCost[string]{
			Node: s.String(),
			Cost: s.cost - node.cost,
		}
	})
}

func (g Graph) String(n string) string {
	return n
}

var final = `#############
#...........#
###A#B#C#D###
  #A#B#C#D#
  #A#B#C#D#
  #A#B#C#D#
  #########`

func main() {
	start := strings.Join(util.ReadLines(os.Args[1]), "\n")

	g := Graph{}
	cost, path := graph.Dijkstra[string](g, start, final)

	fmt.Printf("\n\nFinal path:\n")
	for i, s := range path {
		fmt.Printf("%d\n%s  cost: %d\n\n", i, s.Node, s.Cost)
	}
	fmt.Printf("Total Cost: %d\n\n", cost)
}
