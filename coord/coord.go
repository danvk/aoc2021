package coord

import (
	"aoc/util"
	"fmt"
)

type Coord struct {
	X, Y int
}

func New(X, Y int) Coord {
	return Coord{X, Y}
}

func CoordX(c Coord) int {
	return c.X
}

func CoordY(c Coord) int {
	return c.Y
}

func MaxXY[V any](m map[Coord]V) (int, int) {
	maxX := util.Max(util.Map(util.Keys(m), CoordX))
	maxY := util.Max(util.Map(util.Keys(m), CoordY))

	return maxX, maxY
}

func PrintGrid[V any](m map[Coord]V, blank string, printer func(v V) string) {
	maxX, maxY := MaxXY(m)

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			v, ok := m[Coord{x, y}]
			if ok {
				fmt.Print(printer(v))
			} else {
				fmt.Print(blank)
			}
		}
		fmt.Printf("\n")
	}
}
