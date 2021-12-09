package util

import (
	"bufio"
	"constraints"
	"log"
	"os"
	"strconv"
	"strings"
)

// https://github.com/golang/go/discussions/47203#discussioncomment-1005457
func Any[T any](vs []T, p func(T) bool) bool {
	for _, v := range vs {
		if p(v) {
			return true
		}
	}
	return false
}

func All[T any](vs []T, p func(T) bool) bool {
	for _, v := range vs {
		if !p(v) {
			return false
		}
	}
	return true
}

func AllEq[T comparable](vals []T, val T) bool {
	return All(vals, func(x T) bool { return x == val })
}

func Map[T any, U any](vals []T, fn func(T) U) []U {
	us := make([]U, 0, len(vals))
	for _, v := range vals {
		us = append(us, fn(v))
	}
	return us
}

// More convenient variant of Map for functions that return a value and an error.
func MapErr[T any, U any](vals []T, fn func(T) (U, error)) ([]U, error) {
	us := make([]U, len(vals))
	for i, v := range vals {
		u, err := fn(v)
		if err != nil {
			return nil, err
		}
		us[i] = u
	}
	return us, nil
}

func MapBool[T any, U any](vals []T, fn func(T) (U, bool)) ([]U, bool) {
	us := make([]U, len(vals))
	for i, v := range vals {
		u, ok := fn(v)
		if !ok {
			return nil, ok
		}
		us[i] = u
	}
	return us, true
}

func Filter[T any](vals []T, fn func(T) bool) []T {
	result := []T{}
	for _, v := range vals {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func FlatMap[T any, U any](vals []T, fn func(T) []U) []U {
	us := []U{}
	for _, v := range vals {
		us = append(us, fn(v)...)
	}
	return us
}

func ParseLineAsNums(line string, delim string, skipBlanks bool) []int {
	parts := strings.Split(line, delim)

	nums := []int{}
	for _, part := range parts {
		if len(part) == 0 {
			if skipBlanks {
				continue
			}
			log.Fatalf("Blank in %s", part)
		}
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}
	return nums
}

// Read a file to an array of lines
func ReadLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

// Read a file into blank line-delimited "chunks"
func ReadChunks(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	chunks := [][]string{}
	thisChunk := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			chunks = append(chunks, thisChunk)
			thisChunk = []string{}
		} else {
			thisChunk = append(thisChunk, line)
		}
	}
	if len(thisChunk) > 0 {
		chunks = append(chunks, thisChunk)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return chunks
}

func Max[T constraints.Ordered](nums []T) T {
	if len(nums) == 0 {
		panic(nums)
	}
	max := nums[0]
	for _, v := range nums[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[T constraints.Ordered](nums []T) T {
	if len(nums) == 0 {
		panic(nums)
	}
	min := nums[0]
	for _, v := range nums[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Find the min and max simultaneously, which is slightly more efficient
func MinMax[T constraints.Ordered](nums []T) (T, T) {
	if len(nums) == 0 {
		panic(nums)
	}
	min, max := nums[0], nums[0]
	for _, v := range nums[1:] {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	return min, max
}

// Return the index i which maximizes f(xs(i)), and f(xs(i)).
func ArgMin[T any, U constraints.Ordered](xs []T, f func(T) U) (int, U) {
	if len(xs) == 0 {
		panic("Cannot take ArgMax of empty array")
	}

	min := f(xs[0])
	argMin := 0
	for i := 1; i < len(xs); i++ {
		x := xs[i]
		v := f(x)
		if v < min {
			min = v
			argMin = i
		}
	}

	return argMin, min
}

// Returns a slice of numbers from min to max, inclusive on both ends.
func Seq(min int, max int) []int {
	nums := make([]int, 0, max-min+1)
	for v := min; v <= max; v++ {
		nums = append(nums, v)
	}
	return nums
}

// Returns either (a, b) or (b, a) such that the tuple is ordered
func Ordered[T constraints.Ordered](a T, b T) (T, T) {
	if a <= b {
		return a, b
	}
	return b, a
}

// Construct a row-major array of size wxh (access as mat[y][x]).
func Zeros[T any](w int, h int) [][]T {
	xs := make([][]T, h)
	for y := 0; y < h; y++ {
		xs[y] = make([]T, w)
	}
	return xs
}

type Numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

func Sum[T Numeric](xs []T) T {
	var tally T
	for _, x := range xs {
		tally += x
	}
	return tally
}
