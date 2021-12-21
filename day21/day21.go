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

	die := Die{val: 1}
	p1score, p2score := 0, 0
	for {
		roll := die.Roll() + die.Roll() + die.Roll()
		p1 += roll
		for p1 > 10 {
			p1 -= 10
		}
		p1score += p1
		fmt.Printf("p1 roll: %d -> %d score: %d\n", roll, p1, p1score)
		if p1score >= 1000 {
			break
		}

		roll = die.Roll() + die.Roll() + die.Roll()
		p2 += roll
		for p2 > 10 {
			p2 -= 10
		}
		p2score += p2
		fmt.Printf("p2 roll: %d -> %d score: %d\n", roll, p2, p2score)
		if p2score >= 1000 {
			break
		}

		fmt.Println()
	}

	fmt.Printf("rolls: %d\n", die.num)
	fmt.Printf("p1 score: %d -> %d\n", p1score, p1score*die.num)
	fmt.Printf("p2 score: %d -> %d\n", p2score, p2score*die.num)
}
