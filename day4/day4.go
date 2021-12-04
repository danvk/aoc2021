package main

import (
	"aoc/day1/util"
	"fmt"
	"log"
	"os"
)

type Board struct {
	cells  [][]int
	called [][]bool
}

func ParseBoard(chunk []string) Board {
	cells := make([][]int, len(chunk))
	called := make([][]bool, len(chunk))
	for i, row := range chunk {
		cells[i] = util.ParseLineAsNums(row, " ", true)
		called[i] = make([]bool, len(cells[i]))
		if i > 1 && len(cells[i]) != len(cells[i-1]) {
			log.Fatalf("Length mismatch: %v vs %v", cells[i-1], cells[i])
		}
	}
	return Board{
		cells:  cells,
		called: called,
	}
}

func (b *Board) CallNumber(num int) {
	for i, row := range b.cells {
		for j, cell := range row {
			if cell == num {
				b.called[i][j] = true
			}
		}
	}
}

func readRowColB(called *[][]bool, x0 int, y0 int, dx int, dy int) []bool {
	result := []bool{}
	x := x0
	y := y0
	for {
		if x < 0 || x >= len(*called) || y < 0 || y >= len((*called)[x]) {
			break
		}
		result = append(result, (*called)[x][y])
		x += dx
		y += dy
	}
	return result
}

func readRowColI(cells *[][]int, x0 int, y0 int, dx int, dy int) []int {
	result := []int{}
	x := x0
	y := y0
	for {
		if x < 0 || x >= len(*cells) || y < 0 || y >= len((*cells)[x]) {
			break
		}
		result = append(result, (*cells)[x][y])
		x += dx
		y += dy
	}
	return result
}

func (b *Board) IsWinner() (bool, []int) {
	// Cols
	for x := range b.cells {
		if util.AllTrue(readRowColB(&b.called, x, 0, 0, 1)) {
			return true, readRowColI(&b.cells, x, 0, 0, 1)
		}
	}

	// Rows
	for y := range b.cells {
		if util.AllTrue(readRowColB(&b.called, 0, y, 1, 0)) {
			return true, readRowColI(&b.cells, 0, y, 1, 0)
		}
	}

	return false, nil
}

func (b *Board) SumUnmarked() int {
	count := 0
	for i, row := range b.cells {
		for j, val := range row {
			if !b.called[i][j] {
				count += val
			}
		}
	}
	return count
}

func main() {
	chunks := util.ReadChunks(os.Args[1])

	if len(chunks[0]) != 1 {
		log.Fatalf("Chunks[0] = %v", chunks[0])
	}
	calledNums := util.ParseLineAsNums(chunks[0][0], ",", false)
	fmt.Printf("Called nums: %v\n", calledNums)

	boards := make([]Board, len(chunks)-1)
	for i, chunk := range chunks[1:] {
		boards[i] = ParseBoard(chunk)
		fmt.Printf("Board %d: %v\n", i, boards[i])
	}

	hasWon := make([]bool, len(boards))
out:
	for _, num := range calledNums {
		for i, b := range boards {
			b.CallNumber(num)
			isWinner, cells := b.IsWinner()
			if isWinner {
				hasWon[i] = true
				if util.AllTrue(hasWon) {
					sumUnmarked := b.SumUnmarked()
					fmt.Printf("Winning board %d, number: %d, numbers: %v, sum: %d, answer: %d\n", i, num, cells, sumUnmarked, sumUnmarked*num)
					break out
				}
			}
		}
	}
}
