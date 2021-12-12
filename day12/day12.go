package main

import (
	"aoc/util"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Path struct {
	pos     string
	visited util.Set[string]
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

	paths := []Path{{pos: "start", visited: util.Set[string]{}}}
	completePaths := []Path{}

	for len(paths) > 0 {
		newPaths := []Path{}
		for _, path := range paths {
			nexts, ok := connections[path.pos]
			if !ok {
				continue
			}

			for next := range nexts {
				if !IsLower(next) || !path.visited[next] {
					nextVisited := path.visited.Clone()
					nextVisited.Add(path.pos)
					newPath := Path{
						pos:     next,
						visited: nextVisited,
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
