package util

import (
	"bufio"
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
	us := []U{}
	for _, v := range vals {
		us = append(us, fn(v))
	}
	return us
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
