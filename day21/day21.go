package main

import "fmt"

type Die struct {
	val int
	num int
}

func (d *Die) Roll() int {
	val := d.val
	d.val += 1
	d.num += 1
	if d.val == 101 {
		d.val = 1
	}
	return val
}

func Part1(p1, p2 int) {
	die := Die{val: 1}
	p1score, p2score := 0, 0
	for {
		roll := die.Roll() + die.Roll() + die.Roll()
		p1 += roll
		for p1 > 10 {
			p1 -= 10
		}
		p1score += p1
		// fmt.Printf("p1 roll: %d -> %d score: %d\n", roll, p1, p1score)
		if p1score >= 1000 {
			break
		}

		roll = die.Roll() + die.Roll() + die.Roll()
		p2 += roll
		for p2 > 10 {
			p2 -= 10
		}
		p2score += p2
		// fmt.Printf("p2 roll: %d -> %d score: %d\n", roll, p2, p2score)
		if p2score >= 1000 {
			break
		}

		// fmt.Println()
	}

	fmt.Printf("rolls: %d\n", die.num)
	fmt.Printf("p1 score: %d -> %d\n", p1score, p1score*die.num)
	fmt.Printf("p2 score: %d -> %d\n", p2score, p2score*die.num)
}

// 1, 2, 3
// 1, 2, 3
// 1, 2, 3

type State struct {
	p1, p2           int
	p1score, p2score int
}

func Advance(pos int, roll int) int {
	pos += roll
	for pos > 10 {
		pos -= 10
	}
	return pos
}

func Part2(p1, p2 int) {
	states := map[State]int64{{p1, p2, 0, 0}: 1}

	rolls := map[int]int64{}
	for r1 := 1; r1 <= 3; r1++ {
		for r2 := 1; r2 <= 3; r2++ {
			for r3 := 1; r3 <= 3; r3++ {
				rolls[r1+r2+r3] += 1
			}
		}
	}
	fmt.Printf("die: %#v\n", rolls)

	var p1wins, p2wins int64 = 0, 0

	for step := 1; step <= 100000; step++ {
		for turn := 1; turn <= 2; turn++ {
			nextStates := map[State]int64{}
			fmt.Printf("step %d turn %d size: %d\n", step, turn, len(states))
			for state, count := range states {
				for roll, rCount := range rolls {
					p1, p2 := state.p1, state.p2
					p1score, p2score := state.p1score, state.p2score
					n := count * rCount

					if turn == 1 {
						p1 = Advance(p1, roll)
						p1score += p1
						if p1score >= 21 {
							p1wins += n
						} else {
							nextStates[State{p1, p2, p1score, p2score}] += n
						}
					} else {
						p2 = Advance(p2, roll)
						p2score += p2
						if p2score >= 21 {
							p2wins += n
						} else {
							nextStates[State{p1, p2, p1score, p2score}] += n
						}
					}
				}
			}
			states = nextStates
			fmt.Printf("p1wins: %d, p2wins: %d\n", p1wins, p2wins)
			if len(states) == 0 {
				return
			}
		}
	}
}

const mode = "input"

func main() {
	var p1, p2 int
	if mode == "sample" {
		p1 = 4
		p2 = 8
	} else {
		p1 = 7
		p2 = 6
	}

	// Part1(p1, p2)
	Part2(p1, p2)
}
