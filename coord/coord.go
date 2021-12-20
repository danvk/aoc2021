package coord

import (
	"aoc/util"
	"fmt"
)

type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func New(X, Y int) Coord {
	return Coord{X, Y}
}

// TODO: delete these
func CoordX(c Coord) int {
	return c.X
}

func CoordY(c Coord) int {
	return c.Y
}

func (c *Coord) Neighbors4() []Coord {
	x, y := c.X, c.Y
	return []Coord{
		{x + 1, y},
		{x - 1, y},
		{x, y - 1},
		{x, y + 1},
	}
}

func (c *Coord) Neighbors8() []Coord {
	x, y := c.X, c.Y
	return []Coord{
		{x + 1, y},
		{x - 1, y},
		{x, y - 1},
		{x, y + 1},
		{x + 1, y + 1},
		{x + 1, y - 1},
		{x - 1, y + 1},
		{x - 1, y - 1},
	}
}

func MaxXY[V any](m map[Coord]V) Coord {
	maxX := util.Max(util.Map(util.Keys(m), CoordX))
	maxY := util.Max(util.Map(util.Keys(m), CoordY))

	return Coord{maxX, maxY}
}

func MinXY[V any](m map[Coord]V) Coord {
	maxX := util.Min(util.Map(util.Keys(m), CoordX))
	maxY := util.Min(util.Map(util.Keys(m), CoordY))

	return Coord{maxX, maxY}
}

func PrintGrid[V any](m map[Coord]V, blank string, printer func(v V) string) {
	c := MinXY(m)
	minX, minY := c.X, c.Y
	c = MaxXY(m)
	maxX, maxY := c.X, c.Y

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
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
