package main

import (
	"aoc/util"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Path struct {
	pos         string
	visited     map[string]int
	doubleVisit bool
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func main() {
	// numSteps := flag.Int("steps", 100, "number of steps to simulate")
	// flag.Parse()

	linesText := util.ReadLines(os.Args[1])

	connections := make(map[string]util.Set[string])
	for _, line := range linesText {
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic(parts)
		}
		a := parts[0]
		b := parts[1]
		if _, ok := connections[a]; !ok {
			connections[a] = util.Set[string]{}
		}
		if _, ok := connections[b]; !ok {
			connections[b] = util.Set[string]{}
		}
		connections[a].Add(b)
		connections[b].Add(a)
	}

	paths := []Path{{pos: "start"}}
	completePaths := []Path{}

	for len(paths) > 0 {
		newPaths := []Path{}
		for _, path := range paths {
			pos := path.pos
			nexts, ok := connections[pos]
			if !ok {
				continue
			}

			for next := range nexts {
				if next == "start" {
					continue
				}
				count := path.visited[next]
				isSmallCave := IsLower(next)
				if !isSmallCave || count == 0 || (count == 1 && !path.doubleVisit) {
					nextVisited := util.CopyMap(path.visited)
					nextVisited[pos] += 1
					newPath := Path{
						pos:         next,
						visited:     nextVisited,
						doubleVisit: path.doubleVisit || (isSmallCave && count == 1),
					}
					if next == "end" {
						completePaths = append(completePaths, newPath)
					} else {
						newPaths = append(newPaths, newPath)
					}
				}
			}
		}
		paths = newPaths
	}

	fmt.Printf("Completed paths: %d\n", len(completePaths))
}
