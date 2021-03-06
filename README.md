# Go Lang Notes

## General thoughts on Go

- Things get copied a lot and it's hard to build a mental model for this.
  - structs are copied in a for loop (e.g. Day 6).
  - Arrays are copied by value but slices, which are like a struct with a pointer, length and cap, don't copy the underlying data.
  - structs are copied when you return them from a function.
  - structs are copied when you invoke a method on them with a non-pointer receiver.

- Generics are a great addition. They'll make Go users really want more type inference and more compact function expressions. It will be really nice when Go libraries start supporting generics -- I would have loved a generic set library, or a generic graph library. Coming from a TypeScript perspective, it's interesting to remember that generics can have runtime implications!

- Go is pretty verbose and the standard library is very intentionally not "batteries included". I found the lack of `min`, `max`, `abs` and `reverse` functions particularly annoying. I also wish its testing tools were less boilerplate-y. Sometimes this seems to be the language authors pushing you towards particular solutions, e.g. "when you reverse a slice, maybe you just want reverse iteration."

- There's significantly less type inference that happens than in TS, e.g. for return types of functions and function parameter types in lambdas. In fairness to Go, inference is more ambiguous when you have nominal types and distinctions like `int32` vs. `int64`. But it does feel a bit incongrous with how obsessive Go is about avoiding certain other types of duplication, e.g. `a int, b int` --> `a, b int`.

Small things:

- You don't need regular expressions as much as you think! I wound up using `fmt.Sscanf` for almost all parsing.

- As in other years, using `map[Coord]T` instead of `[][]T` is almost always a good idea.

- Go's type syntax is almost exactly backwards from TS which was a constant source of confusion for me. For example, `[]T` instead of `T[]`. Or for function parameters, `min(a int, b int)` instead of `min(int a, int b)`. Or even `min(a, b int)`.

- I still find Go's pickiness about where I put my code on disk to be weird, confusing and counterproductive. Why can't I `go build` something in a subdirectory? Or maybe only if it's under `GOROOT`?

- The thing about capitalized names for exports didn't bother me, but it did bother me that you lose the ability to define structs by their order outside a module (check that this is true).

- Go has nominal typing (`type Farenheit int`) but I didn't find it helpful that often (possibly on day 8).

- Using `gotip` to access a pre-release mostly worked great. The one downside was that it somtimes tripped up VS Code, and some of the syntax in `constraints.go` (e.g. `~int`) was consistently marked as a syntax error. Go 1.18 came out on Dec. 14th, but I didn't notice or install it.

- It's nice that Go types all have a meaningful zero value. I started to learn to work with this more and more as the month went on.

- There are some pretty error-prone constructs, e.g. `range` giving you the indices first, or the varying behavior of `:=` with multiple variables (does it introduce new variables or reassign old ones?).

https://go.dev/blog/go1.18beta1

## General thoughts on this year's AoC

The 2021 AoC was harder than 2020 but easier than 2019. As in 2020, none of the solutions built on one another. Each day was independent.

Day 24 was the real standout as a creative puzzle.

I continue to wish I were in a timezone where it was more reasonable for me to do the puzzles as soon as they came out. I continue to think this is a great way to learn a language. AoC problems don't really play to Go's strengths (there's no concurrency) but I did feel much more confident writing Go by the end than I did at the start.

I had more trouble finding a Go expert than I did with Rust last year. I think more of the people on reddit who were doing AoC in Go were trying to learn it than were fluent. Eventually @derat noticed one of my solutions and he became my @AxlLind of 2021.

Days that were hard:

- Day 8 part two
- Day 22, part two
- Day 23, both parts (though in retrospect it was just Dijkstra)
- Day 24, part one

## Advent of Code Day by Day

### Day 25 Sea Cucumber (6508 / 4165)

Took 28 minutes while watching the JWST launch and racing to head off to the airport. Had a few off-by-one bugs but the samples revealed them. Started at 6:40:32 and submitted at 7:06:52.

Final rank for completing the whole AoC was 4165 this year.

### Day 24 Arithmetic Logic Unit (4553 / 4422)

This was hard! The successful sequence for me looked something like:

- Evaluate through a few input digits by hand
- Use search in my editor to find common patterns in the input program.
- Pull out the three things that vary in the processing from input-to-input into a `Step` struct and evaluate that.
- Aggressively simplify this until it runs very fast in a loop. Key insight here is that only `z` carries over from input to input.
- Abandoned attempt to suss out relationships between digits by looking at the table of constants.
- Run over 1B random license plates to find ~5 that evaluate to z=0.
- Find the true relationships by looking at these examples and tinkering with them.
- Realize that with a few constraints on the digits, I should be within range of exhaustive search.

In retrospect, I'm not exactly sure why my formulas for the relationships between digits didn't work. They were _close_! I got some of them. And they were the right sort of formulas. But I don't see how you get the relationship between the first three digits and the fourth.

Abandoned ideas:

- Guess that this is based on prime numbers somehow (they're too rare)
- Calculate a range of possible values for each register after each step.

Full notes: [day24/README.md][].

Nice writeup of the code-free way to do this: <https://github.com/dphilipson/advent-of-code-2021/blob/master/src/days/day24.rs>

### Day 23 Amphipod (7404 / 4877)

Another challenging one. I'm glad I went overboard and wrote a very generic Dijkstra for day 15, it came in handy today. I wish I hadn't had a tricky bug in it, though!

The main trick was to switch from using `*State` as the node in my Dijkstra to a `uint64` encoding. Presumably this eliminated tons of copying and memory pressure, and got me deduping of redundant states.

Part two just barely fits in a `uint64`. There are 27 squares in the room in part two and five possible states for each square (A, B, C, D, empty). And as luck would have it, `5**27 < 2**64`. So I can get away with using `uint64`s for part two as well. Dijkstra is pretty magical for problems like this!

**A few hours later** Upon further reflection, the de-duping aspect of using the `uint64` representation is the only thing that really matters. I switched to using `string` instead with the display value as the encoding and it also works great. Dijkstra was slow with `*State` because there were probably millions of duplicate representations of the same state.

### Day 22 Reactor Reboot (10943 / 7610)

Hardest day by far for me; My first instinct was to implement union and difference by splitting up intersecting cuboids into sub-cuboids. There could be up to 27 if you're really unlucky. It seems that some people on reddit got something like this to work, but mine completely blew up on the sample inputs. _(If you're more careful about how you shape the sub-cuboids, you can get the 27 down to 6.)_

I hoped I could simplify the problem by solving independently for each (x, y) pair, so that you only have to do interval arithmetic. This would have worked and didn't take too much memory, but I clocked it at taking an estimated ~140 years. Ouch.

I hoped that solving for each z value and looking at rectangles in the (x, y) plane would be a nice middle ground, but this also blew up.

While driving in the afternoon, I came up with two new ideas:

1. Store separate "include" and "exclude" cubes. So you can represent A + B as [A, B] - [A ??? B]. It sounds like @igheorghita got something like this to work. I was concerned that operations would blow up as the number of sets in each list increased.
2. There are only so many distinct xs, ys and zs. So if you just make a list of them and index it, you can re-use your solution from part 1 and scale each cell up by its volume.

This is what I wound up doing. It was quite fast on the third sample, but was taking a long but not interminable time on my input (maybe an hour?). I replaced my `map[Coord]bool` structure with a gigantic array (`[1663][1663][1663]bool`) and it ran in ~10s.

@znrk found a similar solution: <https://gist.github.com/znkr/a1c8dfcbedaabb5b97a7886f06a40282>. He made a clever choice to change the closed intervals to open intervals by adding one to the right end of each. This cuts the size of the giant grid by a factor of 8. It seems you can't have dynamically-sized multidimensional arrays in Go. Switching to a big buffer and calculating offsets by hand shaved a bit more time off. All told, I'm at ~1s for part 2.

I'm pretty curious to see how others got this to work using the more geometric approaches.

This `//go:embed` looks interesting: <https://github.com/pemoreau/advent-of-code-2021/blob/main/go/22/day22.go>

@pemoreau seems to have gotten my initial approach to work: <https://github.com/pemoreau/advent-of-code-2021/blob/830619e40b11a0042875b8383978ed8bd1565e5a/go/22/cuboid.go#L75>

This solution is fascinating:
<https://topaz.github.io/paste/#XQAAAQCnAgAAAAAAAAAzHIoib6qqOe07MhJ0XsXE6K08G4Ps1pgTxGMtEZ+kpb0WiMmgclVAwbWGLuEqShaMvaIHbGSZQDr1DzD4YJdTRL4c/0gztLN1zvPRMMuPZI6AjcJ1jQMQwV/eQ4Xx+ZfU1PZoy8ITNoTLg8ND9SSxm0z9oF/VvJvPMcpFJJpTBmm89nO7Kj8zuP1/7GMgKDcDhV3H86rhgHybzKsU1vco+QkpaXh8sDhEfXUe0wM2szbwkYTiIp8UAJyT566Us7JiKvV0S7lxR7dkUDFnAliMQG+BOd+p0kZ8MhuDrFS2Ujxcq1NuqpMlRJVcYAu+2MoMVtqE0OGZQLXqGLMIiVVJuEkZw5OtZ7Pkpo0UcSOcjgq/7o97QucphHs0ZR8Kxnhh/W+3mIz4wqmyXgXqYQy3OqjiI56yy8n688wWq/KFyuzOTLDw+TSopvoOHvFZfB/z4+7ofeGsw0LqJU/tB1d1JQIc9G/WCgt/L18jsahMlvlYcMYrzUD8WeTzOps3J3761hCM>

This is a slightly more legible version that I think is equivalent: <https://github.com/BradonZhang/advent-of-code-2021/blob/main/src/22.py#L30>

The inclusion/exclusion hierarchy is actually a _tree_. For my input it has a depth of 15 and 35,924 nodes. Which is a lot but not excessive!

### Day 21 Dirac Dice (14059 / 8066)

Two very different implementations for part 1 and part 2 today, but neither was very difficult. The state size for part 2 never got larger than ~12,000.

### Day 20 Trench Map (9773 / 9773)

This one doesn't make sense to me -- because the first bit in my decoder string is `#`, doesn't that mean that the entire infinite grid should light up immediately?

YES! But then it turns off again after the second iteration. I like it!

### Day 19 Beacon Scanner (2985 / 3513)

I just guessed the right answer for part 1 without writing any code! Two hours later I had it for real. My idea on how to solve this was correct (try all rotations, find the best shift), it was just tricky to implement without bugs.

I'm getting 48 rotations, not 24.
*It's because I was including flips, whereas the problem just wants rotations. I thought the wording here was pretty ambiguous. See: [reddit comment][48]*

My initial version took ~3 minutes to run. After some light optimizations, it's down to 8 seconds.

@rvilim said he just looked up the list of 24 rotation matrices. In retrospect I'm not sure why that felt like cheating to me, but guessing an answer didn't.

[48]: https://www.reddit.com/r/adventofcode/comments/rjpf7f/2021_day_19_solutions/

### Day 18 Snailfish (8306 / 8048)

Quite a slog today -- I missed that you had to reduce after _each_ addition, rather than after all the additions. I also had a bug where I added the exploded numbers to _all_ the numbers on either side, rather than just the first.

In retrospect, I think a better solution would have been to simultaneously maintain a linear order and tree structure. Storing a pointer to parent also seems very helpful, though this is tricky in Go with its pass by value. Pointers are "contagious" in that, once you start using pointers, you more or less have to use them everywhere if you want to maintain pointer equality.

I reworked the problem with this approach. There is a bit more bookkeeping, but I do think overall it's much simpler.

I also set up some boilerplate for "table-driven" tests following <https://dave.cheney.net/2019/05/07/prefer-table-driven-tests>.

### Day 17 Trick Shot (14797 / 13310)

First wrong answer on part 2! 3092 is too low.

It turned out to be a `<` that should have been a `<=` in my calculation of the minimum possible `vx`. This mattered for my puzzle input but not for the sample input.

### Day 16 Packet Decoder (11533 / 10121)

Using a slice of strings that are `"0"` or `"1"` is inefficient but extremely convenient here. Writing out `PopBits` and `PopAndDecode` helpers early made this pretty straightforward.

My one hiccup was that I used `int64` everywhere, but still had `strconv.ParseInt(bits, 2, 32)`, which meant I `panic`ed on one of the examples.

You _can_ pass methods as function parameters in Go! This:

		vals := util.Map(p.packets, func(p Packet) int64 { return p.Evaluate() })

becomes this:

		vals := util.Map(p.packets, Packet.Evaluate)

Neat! <https://stackoverflow.com/a/38897667/388951>

Why do my `Init()` functions never get called? (@derat points out that it should be `init`. Whoops!)

### Day 15 Chiton (16262 / 12517)

My first attempt at a solution involved keeping track of all the paths through the grid as I did a BFS and pruning once I got the first complete path. This didn't scale at all. Instead I figured I could do Dijkstra: keep track of the minimum distance to each node in the graph, and only add that node to the fringe if you reduce the distance.

Part two was challenging because you really had to read the question carefully: it's not just `risk%10` because it wraps around to `1` not `0`.

Definitely happy that I factored out a `coord` package a few days ago!

I looked at this graph package but I wish it were generic! I think this will really force Go developers to update their packages and create opportunities for new ones <https://pkg.go.dev/github.com/yourbasic/graph>.

I'm trying to write my own Dijkstra using `containers/heap`. One head scratcher: this method makes my heap not match the interface:

		func (h *THeap[T]) Pop() NodeWithCost[T]{}

instead Go wants:

		func (h *THeap[T]) Pop() interface{}

which you'd think would be matched by mine.

Generic Dijkstra works pretty well! And it's faster than what I wrote before because it uses the heap structure to visit nodes in order of increasing risk.

This is subtle:

		node, ok := prev[node]

This is fine if `node` was already declared but `ok` was not. But it makes a new `node`, rather than reassigning it.

### Day 14 Extended Polymerization (23122 / 14966)

I enjoyed today's! My part 1 solution really didn't scale for part 2. Updating it to just keep track of the counts of each pair was simple enough, but I had an "oh crap" moment as I was trying to count the _individual_ molecules to get the final answer. Then I realized that you can just count the second molecule in each pair, and treat the first molecule in the template (which is fixed) specially.

This was vaguely reminiscent of a puzzle in 2019 that I found quite hard.

### Day 13 Transparent Origami (23310 / 22213)

Using a map from `Coord` -> `bool` is a good idea yet again.

With

		func Keys[K comparable](m map[K]interface{}) []K {

this code produces an error:

		func PrintDots(dots map[Coord]bool) {
			keys := util.Keys(dots)
		}
		// type map[Coord]bool of dots does not match inferred type map[Coord]interface{} for map[K]interface{}

whereas with this declaration it works:

		func Keys[K comparable, V any](m map[K]V) []K {

Is `any` not an alias for `interface{}`?

Time to factor out a `coord` package. I'm sure there are reasons, but it's a bit annoying that `Coord{x, y}` only works if `Coord` is defined in your own package. It has to be either `c.Coord{X, Y}` or `c.Coord{X: x, Y: y}` if it's defined in another package.

### Day 12 Passage Pathing (21640 / 19360)

Why

		paths := []Path{{pos: "start", visited: map[string]int{}}}

and not

		paths := []Path{{pos: "start", visited: {}}}

I guess this is even shorter:

		paths := []Path{{pos: "start"}}

### Day 11 Dumbo Octopus (23611 / 23275)

I thought this struct would be the way to go:

```go
type Octopus struct {
	val     int
	flashed bool
}
```

but you cannot assign to a field in a struct in a map in Go (<https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map>) which makes this extremely onerous.

I ran into a very confusing bug where my grid diverged from the sample after 8 steps. But if I used the step 7 grid as the input, it got the first step right. It turned out to be an issue with Go's random order traversal of maps. But is this really an underspecified problem? Or maybe the issue was that I'm mutating a grid, rather than creating a new one within each substep?

### Day 10 Syntax Scoring (30579 / 28392)

I started with:

```go
type Chunk struct {
	openChar    rune
	start, stop int
	children    []Chunk
}
```

but then realized I just needed a stack for part 1. I'm glad I went with the simpler approach -- the remaining elements on the stack were just exactly what I needed for part 2. I embraced "rune"s today. _(In later days I tended to just go with `string`.)_

I was glad my solution today just worked on the first go -- this would have been a real pain to debug!

### Day 9 Smoke Basin (32114 / 22829)

First problem using flood fill / BFS. From doing many whiteboard interviews and previous AoCs, this is definitely one I can implement! My `Set[T]` implementation from yesterday was moderately useful.

### Day 8 Seven Segment Search (33759 / 24434)

Part two was the first puzzle that felt hard. I thought it was going to be _quite_ hard, but then saw that there were some consistent patterns to exploit that made things easier. I think my choice of data structures may have made this a bit harder on me, too.

Nominal typing did come in handy and prevented at least one possibly confusing bug:

```ts
type ScrambledDigit byte
type Digit byte
```

Thinking more about this, a lot of the awkwardness came down to not having a nice set structure. Here's a Gopher who did this using a third-party library:
https://github.com/microhod/adventofcode/blob/main/2021/08/main.go
https://github.com/deckarep/golang-set

I like this but wish it were generic! I started playing around with building my own but got bogged down porting my solution to it. I'm sure this will come in useful again for future problems.

It's interesting that you can do runtime type checks to get around some limitations of Go generics, e.g. to use a `String()` method on `T` to stringify `Set[T]`.

I'm starting to think Go's niche for me personally might be command line tools.

### Day 7 The Treachery of Whales (35665 / 33281)

Not much to this one. The most surprising thing to me was that Go doesn't have an integer `Abs` function, the rationale being that it's easy to implement and you should just write it yourself: <https://stackoverflow.com/questions/57648933/why-doesnt-go-have-a-function-to-calculate-the-absolute-value-of-integers>

I reworked my initial answer to be a bit more functional, with `Seq` and `ArgMin` helpers. Not sure which version I prefer -- the verbosity of Go lambdas doesn't make you want to use them. Performance seems identical for both variants (for loop and `ArgMin`).

### Day 6 Lanternfish (34197 / 27030)

First day where the obvious implementation of part 1 is too slow for part 2.

What tripped me up today was that structs are passed around by value, so there's implicit copying that happens even in a `for` loop. So in this code:

```go
func Advance(school *[]Lanternfish) {
	newFish := []Lanternfish{}
	for _, fish := range *school {
		if fish.timer == 0 {
			newFish = append(newFish, Lanternfish{8})
			fish.timer = 6
		} else {
			fish.timer -= 1
		}
	}
}
```

The struct referenced by `fish` is a _copy_, not the struct in the array. So the changes you make to it get thrown away. You need to mutate it via `fish[i].timer -= 1`.

I need to get a better sense for where implicit copies happen in Go. This also tripped me up with the REPL.

Some notes from skimming the Go Programming Language book:

- Go strings are immutable. I guess JavaScript strings are immutable, too, but it feels more surprising in Go.
- Go arrays are fixed length, which is why you normally work with slices.
- Go has built-in `complex64` and `complex128` types. Might be useful for some puzzles, or as an alternative to a Coordinate structure.
- Every type has a zero value, it's `nil` for most types.
- You can do `type Celsius int32` to get something that's neither comparable nor assignable to a plain `int32`. This explains why the generics proposal talks so much about "underlying types", i.e. `~T`.
- You can use `:=` for reassignment, so long as it declares at least one new var.
- They say Go is fun, so I guess I'm wrong about that.

### Day 5 Hydrothermal Venture (32872 / 29337)

The only really new thing today was using regular expressions to parse input. It looks like `FindStringSubmatch` is going to be my friend <https://pkg.go.dev/regexp#Regexp.FindStringSubmatch>.

_(In retrospect, this was the only time I used a regular expression! I found fmt.Sscanf to be more convenient for parsing.)_

I defined a generic `Min` function and used it along with `FlatMap` to get the max X and Y values for the grid. I think this is cute, but maybe not idiomatic Go. At least not yet? I do think having access to generics greatly improves my happiness working with Go.

### Day 4 Giant Squid (33362 / 29925)

The `%v` formatting just uses spaces to delimit array elements, which is pretty ambiguous about where arrays begin and end!

    len(chunk[1]) = 5, chunk[1]: [22 13 17 11  0  8  2 23  4 24 21  9 14 16  7  6 10  3 18  5  1 12 20 15 19]

(I later found out that `%#v` prints values in a way that looks more like Go syntax.)

Weird that you can return a tuple but can't pass tuples as a parameter. I guess tuple types aren't a thing in Go, just multiple return values.

I got tripped up and misread the solution as the sum of the winning squares, so I wound up implementing some unneeded functions.

The lack of generics is pretty crazy. Today might be a good chance to try out Go 1.18.

I posted my solution on the Megathread. Not a ton of Go solutions:
<https://www.reddit.com/r/adventofcode/comments/r8i1lq/2021_day_4_solutions/>

Generics are nice to have. Very puzzled by this error:

```go
func AllEq[T comparable](vals []T, val T) bool {
	return All(vals, func(x T) { return x == val })
    //               ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    // type func(x T) of (func(x T) literal) does not match inferred type func(T) bool for func(T) boolcompilerErrorCode(138)
}
```

Changing it to `func(x T) bool { return x == val }` fixes the error, but what does Go think the return type is otherwise? (This turns out to be confusing formatting in the error message.)

Here's the proposal: <https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md> It's notably agnostic about the implementation. Unclear to me whether generic functions are implemented via erasure or by instantiating multiple versions of the function.

### Day 3 Binary Diagnostic (56777 / 37323)

The distinction between chars and bytes is a bit annoying. I got tripped up by doing:

    int(str[pos])

which gives you the ASCII value of the byte at that pos, rather than parsing it as a number. So `int('0')` --> `48` rather than `0`.

Factoring out `mostCommonBitAtPos` and `filterByBit` functions made part two easier to think about. I'm starting to get to the point where I'm curious how a Go expert would do this.

### Day 2 Dive! (63117 / 59204)

Introduces some string parsing. In Rust I would have set up an enum to represent the movements, and in TypeScript I'd set up a union type. This construct seems to be absent from Go. I guess I could have set up a `struct` since the moves all have the same shape.

### Day 1 Sonar Sweep (49425 / 43000)

Things to set up here were:

- Iterating over the lines in a file
- Parsing ints
- Manipulating arrays
- Looping over arrays

I got tripped up by the behavior of `range`. If you do:

```go
for x := range xs {}
```

then you loop over the _indices_, not the values. If you want the values, you have to use the two-element form:

```go
for _, x := range xs {}
```

This is pretty hilarious: <https://stackoverflow.com/a/21950283/388951>, <https://github.com/bradfitz/iter/blob/master/iter.go>.

## What I know / remember about Go going in

- I learned a bit of Go back in 2011/2012 and [implemented boggle][gobog].
- I remember running into lots of trouble defining the Trie and Boggler interfaces, and then having an insight about how in Go, the Trie interface belongs with the Boggler (the thing that uses it) rather than the Trie itself (the thing that implements it). This is the essence of Duck Typing.
- But then I couldn't make this work because Go didn't support recursive duck typing. There was an annoying quote about this: "In Go, if it walks like a duck it is a duck, but only if it uses the same legs as a duck." [rec-quote] It's interesting that TypeScript makes the exact opposite choice here given the same hard problem ("it's turtles all the way down, but in practice there are only five turtles").
- I used Go for logs analysis a bit later at Google ("Lingo", aka Logs in Go) but I remember this being a worst-of-both-worlds sort of thing: the Go compiler is famously fast, but that's lost if you have to spend several minutes linking and then kicking off MapReduce jobs.
- I have dabbled in Go a bit since then. I find the reliance on slices to be a bit strange and I don't love how insistent it is about how you organize your code directories at the level above your project.
- Some of Go's quirks, like canonical formatting, have aged well and are standard now.
- I read that Go is getting generics and am excited to give them a try!
- I remember not liking how so many functions return `nil` and errors, leading to lots of `if` statements. But I think now I might not mind that so much? Forcing you to deal with errors seems like a good idea.

Following the example, I'm remembering finding it quite annoying how aggressive `gofmt` is about stripping unused variables and imports when you save. Often it's just that I'm writing code and haven't used them _yet_.

I think the module story has changed since the last time I used Go. I don't remember running `go mod tidy` to download modules or importing from `rsc.io`. These modules get downloaded to `~/go/pkg/mod`.

Following <https://go.dev/doc/tutorial/create-module>, making a `greetings.go` file in the same directory as `hello.go` gives me some confusing errors (commit `3a7924c`). A file with `package greetings` needs to be in a directory called `greetings`. Depending on which directory you open VS Code in, you'll get package errors. (`code .` vs. `cd hello && code .` in the example.) The only real resolution here seemed to be to move my code to `~/go/src/example.com`.

Go 1.18 will be the first release with generics. It's [expected in Feb. 2022][go-118] but maybe I can do Advent of Code in a pre-release? Yes, via `gotip` https://sebastian-holstein.de/post/2021-11-08-go-1.18-features/. A beta may be released sometime in December.

The Go debugger is called??? Delve! <https://github.com/go-delve/delve>

This seems to be the canonical blog post on Go generics at the moment: <https://bitfieldconsulting.com/golang/generics>. Many of its examples break the "golden rule of generics" from my blog.

Q: How do you run a go module from another directory? `go run ??` instead of `cd hello && go run .`

## Working on the REPL

Can you put more than one file in a directory in Go? I keep getting errors about package `main` vs. whatever I want to call the package for the other files. Maybe subdirectories, and your `main` is top level.

I miss union types! <https://making.pusher.com/alternatives-to-sum-types-in-go/>
Go doesn't have a `map` function because no generics: <https://stackoverflow.com/a/50025091/388951>

Is there any downside to pass-by-value in Go? Passing a struct vs. pointer to struct?
-> Looks like it has the same copy mechanics as C++. This can get confusing if you care about pointer equality.

It's pretty shocking to me that there's no built-in funtion to reverse an array? <https://stackoverflow.com/a/19239850/388951>.

I'm not optimistic about Go after doing this exercise. REPL was awkward, it's verbose, and there's lots of missing built-ins. But maybe it gets better from here?

[gobog]: https://github.com/danvk/performance-boggle/tree/master/go/boggle
[rec-quote]: https://github.com/golang/go/issues/1074#issuecomment-66052455
[go-118]: https://go.dev/blog/12years
