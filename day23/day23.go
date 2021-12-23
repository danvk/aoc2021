package main

import (
	"aoc/graph"
	"aoc/util"
	"fmt"
	"os"
	"strings"
	"time"
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
	num       int // either 8 or 16
}

func (s State) String() string {
	var rows [5][]string
	rows[0] = strings.Split("#...........#", "")
	rows[1] = strings.Split("###.#.#.#.###", "")
	rows[2] = strings.Split("  #.#.#.#.#", "")
	rows[3] = strings.Split("  #.#.#.#.#", "")
	rows[4] = strings.Split("  #.#.#.#.#", "")
	for _, a := range s.amphipods[:s.num] {
		rows[a.y][a.x+1] = a.kind
	}
	// fmt.Printf("rows: %#v\n", rows)

	result := fmt.Sprintf(
		"#############\n%s\n%s\n%s\n",
		strings.Join(rows[0], ""),
		strings.Join(rows[1], ""),
		strings.Join(rows[2], ""),
	)

	if s.num == 16 {
		result += fmt.Sprintf(
			"%s\n%s\n",
			strings.Join(rows[3], ""),
			strings.Join(rows[4], ""),
		)
	}

	return result + "  #########"
}

func ParseState(state string) State {
	lines := strings.Split(state, "\n")
	var num int
	if len(lines) == 5 {
		num = 8
	} else if len(lines) == 7 {
		num = 16
	} else {
		panic(state)
	}
	s := State{num: num}
	i := 0
	for y := 0; y < len(lines)-2; y++ {
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
	if i != num {
		fmt.Printf("lines: %#v\n", lines)
		panic(i)
	}
	return s
}

func (s *State) IsHallwayOpen(start, stop int, selfIdx int) bool {
	x1, x2 := util.Ordered(start, stop)
	for i, a := range s.amphipods[:s.num] {
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
	for _, a := range s.amphipods[:s.num] {
		result[a.y][a.x] = a.kind
	}
	return
}

func (s *State) NextStates() []*State {
	var nextStates []*State
	g := s.ToGrid()

	maxY := 4
	if s.num == 8 {
		maxY = 2
	}

	for i, a := range s.amphipods[:s.num] {
		dx := destX[a.kind]
		if a.y == maxY && a.x == dx {
			continue // already in its place, no need to move
		}

		if a.y == 0 {
			dy := -1
			for y := maxY; y >= 1; y-- {
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
			for y := a.y + 1; y <= maxY; y++ {
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

	// for i, s := range nextState {
	// 	fmt.Printf("%d:\n%s\n\n", i, s)
	// }

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

var final1 = `#############
#...........#
###A#B#C#D###
  #A#B#C#D#
  #########`

var final2 = `#############
#...........#
###A#B#C#D###
  #A#B#C#D#
  #A#B#C#D#
  #A#B#C#D#
  #########`

func AddStep2Lines(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) != 5 {
		panic(text)
	}

	// Surely there's an easier way to do this!
	var nextLines []string
	nextLines = append(nextLines, lines[:3]...)
	nextLines = append(nextLines, "  #D#C#B#A#", "  #D#B#A#C#")
	nextLines = append(nextLines, lines[3], lines[4])
	return strings.Join(nextLines, "\n")
}

func main() {
	start := strings.Join(util.ReadLines(os.Args[1]), "\n")

	startT := time.Now()
	g := Graph{}
	p1cost, path := graph.Dijkstra[string](g, start, final1)
	p1Time := fmt.Sprintf("%v", time.Since(startT))

	fmt.Printf("Part 1:\n")
	fmt.Printf("Final path:\n")
	for i, s := range path {
		fmt.Printf("%d\n%s  cost: %d\n\n", i, s.Node, s.Cost)
	}

	startT = time.Now()
	start = AddStep2Lines(start)
	p2cost, path := graph.Dijkstra[string](g, start, final2)
	p2Time := fmt.Sprintf("%v", time.Since(startT))

	fmt.Printf("Part 2:\n")
	fmt.Printf("Final path:\n")
	for i, s := range path {
		fmt.Printf("%d\n%s  cost: %d\n\n", i, s.Node, s.Cost)
	}

	fmt.Printf("Total Cost Part 1: %d (%s)\n", p1cost, p1Time)
	fmt.Printf("Total Cost Part 2: %d (%s)\n", p2cost, p2Time)
}
