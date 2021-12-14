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

func AdvanceMap(pairCounts map[string]int64, rules map[string]rune) map[string]int64 {
	newCounts := make(map[string]int64)
	for pair, count := range pairCounts {
		post, ok := rules[pair]
		if !ok {
			panic(pair)
		}
		pre := []rune(pair)
		a := pre[0]
		b := pre[1]

		newCounts[string([]rune{a, post})] += count
		newCounts[string([]rune{post, b})] += count
	}
	return newCounts
}

func TemplateToPairs(template string) map[string]int64 {
	chars := []rune(template)
	res := make(map[string]int64)
	for i := 1; i < len(chars); i++ {
		two := string([]rune{chars[i-1], chars[i]})
		res[two] += 1
	}
	return res
}

func GetCharCounts(initTemplate string, counts map[string]int64) map[rune]int64 {
	out := make(map[rune]int64)
	for pair, count := range counts {
		chars := []rune(pair)
		if len(chars) != 2 {
			panic(pair)
		}
		b := chars[1]
		out[b] += count
	}

	initRunes := []rune(initTemplate)
	out[initRunes[0]] += 1
	return out
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
	pairs := TemplateToPairs(template)
	for step := 1; step <= 40; step++ {
		pairs = AdvanceMap(pairs, rules)
		if step <= 4 {
			fmt.Printf("%d: %#v\n", step, pairs)
		}
	}

	counts := GetCharCounts(template, pairs)
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
