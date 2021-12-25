package main

import (
	c "aoc/coord"
	"aoc/util"
	"fmt"
	"os"
)

// east (>)
// south (v)
// empty (.)

// Every step, the sea cucumbers in the east-facing herd attempt to move forward one location
// then the sea cucumbers in the south-facing herd attempt to move forward one location

// sea cucumbers that move off the right edge of the map appear on the left edge
// Sea cucumbers always check whether their destination location is empty before moving, even if that destination is on the opposite side of the map:

func MoveEast(g map[c.Coord]rune, max c.Coord) bool {
	maxX := max.X
	toMove := map[c.Coord]c.Coord{}
	for p, cuc := range g {
		if cuc == '>' {
			x, y := p.X, p.Y
			next := c.Coord{X: (x + 1) % maxX, Y: y}
			if _, ok := g[next]; !ok {
				toMove[p] = next
			}
		}
	}

	// Make the moves
	for cur, next := range toMove {
		g[next] = g[cur]
		delete(g, cur)
	}

	return len(toMove) > 0
}

func String(g map[c.Coord]rune, max c.Coord) string {
	text := ""
	for y := 0; y <= max.Y; y++ {
		for x := 0; x <= max.X; x++ {
			if cuc, ok := g[c.Coord{X: x, Y: y}]; ok {
				text += fmt.Sprintf("%c", cuc)
			} else {
				text += "."
			}
		}
		text += "\n"
	}
	return text
}

func main() {
	linesText := util.ReadLines(os.Args[1])

	maxX := -1
	maxY := -1
	g := map[c.Coord]rune{}
	for y, line := range linesText {
		for x, ch := range line {
			if ch == '>' || ch == 'v' {
				g[c.Coord{X: x, Y: y}] = ch
			}
			if x > maxX {
				maxX = x
			}
		}
		if y > maxY {
			maxY = y
		}
	}
	max := c.Coord{X: maxX, Y: maxY}
	fmt.Printf("max: %s\n", max)
	fmt.Printf("%s\n", String(g, max))
	MoveEast(g, max)
	fmt.Printf("\n\nEast:\n%s\n", String(g, max))
}
