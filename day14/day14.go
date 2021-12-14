package main

import (
	"aoc/util"
	"fmt"
	"os"
	"sort"
)

func Advance(template string, rules map[string]rune) string {
	chars := []rune(template)
	next := []rune{chars[0]}

	for i := 0; i < len(chars)-1; i++ {
		a := chars[i]
		b := chars[i+1]
		pre := string([]rune{a, b})
		post, ok := rules[pre]
		if !ok {
			panic(pre)
		}
		next = append(next, post, b)
	}
	return string(next)
}

func main() {
	linesText := util.ReadChunks(os.Args[1])
	if len(linesText) != 2 {
		panic(linesText)
	}

	template := linesText[0][0]
	rulesText := linesText[1]

	rules := make(map[string]rune)
	for _, line := range rulesText {
		var pre string
		var post rune
		_, err := fmt.Sscanf(line, "%s -> %c", &pre, &post)
		if err != nil {
			panic(line)
		}
		rules[pre] = post
	}

	fmt.Printf("Template: %s\n", template)
	for step := 1; step <= 10; step++ {
		template = Advance(template, rules)
		if step <= 4 {
			fmt.Printf("%d: %s\n", step, template)
		}
	}

	counts := make(map[rune]int)
	for _, c := range template {
		counts[c] += 1
	}
	syms := util.Keys(counts)
	fmt.Printf("syms: %#v\n", syms)
	sort.Slice(syms, func(i, j int) bool {
		return counts[syms[j]] > counts[syms[i]]
	})

	fmt.Printf("Counts: %d\n", counts)
	fmt.Printf("syms: %#v\n", syms)
	most, least := counts[syms[len(syms)-1]], counts[syms[0]]
	fmt.Printf("Span: %d - %d = %d\n", most, least, most-least)
}
