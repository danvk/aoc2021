package util

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func AllEq[T comparable](vals []T, val T) bool {
	for _, v := range vals {
		if v != val {
			return false
		}
	}
	return true
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
