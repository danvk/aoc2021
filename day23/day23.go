package main

import (
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
	amphipods [8]Amphipod
	cost      int
}

func (s State) String() string {
	var rows [3][]string
	rows[0] = strings.Split("#...........#", "")
	rows[1] = strings.Split("###.#.#.#.###", "")
	rows[2] = strings.Split("  #.#.#.#.#", "")
	for _, a := range s.amphipods {
		rows[a.y][a.x+1] = a.kind
	}

	return fmt.Sprintf(
		"#############\n%s\n%s\n%s\n  #########",
		strings.Join(rows[0], ""),
		strings.Join(rows[1], ""),
		strings.Join(rows[2], ""),
	)
}

func ParseState(state string) State {
	lines := strings.Split(state, "\n")
	if len(lines) != 5 {
		panic(lines)
	}
	s := State{}
	i := 0
	for y := 0; y <= 2; y++ {
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

func (s *State) IsHallwayOpen(start, stop int) bool {
	x1, x2 := util.Ordered(start, stop)
	for _, a := range s.amphipods {
		if a.y == 0 && a.x >= x1 && a.x <= x2 {
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

func (s State) Move(idx int, x, y int) State {
	// value receiver, so a copy is made for us
	a := s.amphipods[idx]
	moveD := abs(x-a.x) + abs(y-a.y)
	moveCost := moveD * costs[a.kind]

	s.cost += moveCost
	s.amphipods[idx].x = x
	s.amphipods[idx].y = y
	return s
}

func (s *State) NextStates() []State {
	var nextStates []State
	for i, a := range s.amphipods {
		dx := destX[a.kind]
		if a.y == 2 && a.x == dx {
			continue // already in its place, no need to move
		}
		if a.y == 1 && a.x == dx {
			if below, ok := s.AmphipodAt(a.x, 2); ok && below.kind == a.kind {
				continue // this stack is already set
			}
		}

		if a.y == 0 {
			// Move to final destination if possible
			bottom, okB := s.AmphipodAt(dx, 2)
			_, okT := s.AmphipodAt(dx, 1)
			dy := 2
			if okT && okB {
				continue // this column is fully occupied
			} else if !okB {
				// bottom slot is open, let's take it!
			} else if bottom.kind == a.kind {
				dy = 1 // bottom slot is taken by same kind of amphipod, so we can take the top slot
			} else {
				continue // already another amphipod here
			}

			if s.IsHallwayOpen(a.x, dx) {
				nextStates = append(nextStates, s.Move(i, dx, dy))
			}
			continue
		}

		if a.y == 2 {
			_, ok := s.AmphipodAt(a.x, 1)
			if ok {
				continue // someone is on top of us, we can't move
			}
		} else if a.y != 1 {
			panic(a)
		}

		// Move into hallway
		for _, x := range hallwaySpots {
			if s.IsHallwayOpen(a.x, x) {
				nextStates = append(nextStates, s.Move(i, x, 0))
			}
		}
	}

	return nextStates
}

var costs = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
	"D": 4,
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

func main() {
	text := strings.Join(util.ReadLines(os.Args[1]), "\n")
	state := ParseState(text)
	fmt.Printf("Parsed state:\n%s\n", state)
}
