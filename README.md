# Go Lang Notes

## Advent of Code

### Day 11

I thought this struct would be the way to go:

```go
type Octopus struct {
	val     int
	flashed bool
}
```

but you cannot assign to a field in a struct in a map in Go (<https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map>) which makes this extremely onerous.

I ran into a very confusing bug where my grid diverged from the sample after 8 steps. But if I used the step 7 grid as the input, it got the first step right. It turned out to be an issue with Go's random order traversal of maps. But is this really an underspecified problem? Or maybe the issue was that I'm mutating a grid, rather than creating a new one within each substep?

### Day 10

I started with:

```go
type Chunk struct {
	openChar    rune
	start, stop int
	children    []Chunk
}
```

but then realized I just needed a stack for part 1. I'm glad I went with the simpler approach -- the remaining elements on the stack were just exactly what I needed for part 2. I embraced "rune"s today.

### Day 9

First problem using flood fill / BFS. From doing many whiteboard interviews and previous AoCs, this is definitely one I can implement! My `Set[T]` implementation from yesterday was moderately useful.

### Day 8

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

### Day 7

Not much to this one. The most surprising thing to me was that Go doesn't have an integer `Abs` function, the rationale being that it's easy to implement and you should just write it yourself: <https://stackoverflow.com/questions/57648933/why-doesnt-go-have-a-function-to-calculate-the-absolute-value-of-integers>

I reworked my initial answer to be a bit more functional, with `Seq` and `ArgMin` helpers. Not sure which version I prefer -- the verbosity of Go lambdas doesn't make you want to use them. Performance seems identical for both variants (for loop and `ArgMin`).

### Day 6

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

The struct referenced by `fish` is a copy, not the struct in the array. So the changes you make to it get thrown away. You need to mutate it via `fish[i].timer -= 1`.

I need to get a better sense for where implicit copies happen in Go. This also tripped me up with the REPL.

Some notes from skimming the Go Programming Language book:

- Go strings are immutable. I guess JavaScript strings are immutable, too, but it feels more surprising in Go.
- Go arrays are fixed length, which is why you normally work with slices.
- Go has built-in `complex64` and `complex128` types. Might be useful for some puzzles, or as an alternative to a Coordinate structure.
- Every type has a zero value, it's `nil` for most types.
- You can do `type Celsius int32` to get something that's neither comparable nor assignable to a plain `int32`. This explains why the generics proposal talks so much about "underlying types", i.e. `~T`.
- You can use `:=` for reassignment, so long as it declares at least one new var.
- They say Go is fun, so I guess I'm wrong about that.

### Day 5

The only really new thing today was using regular expressions to parse input. It looks like `FindStringSubmatch` is going to be my friend <https://pkg.go.dev/regexp#Regexp.FindStringSubmatch>.

I defined a generic `Min` function and used it along with `FlatMap` to get the max X and Y values for the grid. I think this is cute, but maybe not idiomatic Go. At least not yet? I do think having access to generics greatly improves my happiness working with Go.

### Day 4

The `%v` formatting just uses spaces to delimit array elements, which is pretty ambiguous about where arrays begin and end!

    len(chunk[1]) = 5, chunk[1]: [22 13 17 11  0  8  2 23  4 24 21  9 14 16  7  6 10  3 18  5  1 12 20 15 19]

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

Changing it to `func(x T) bool { return x == val }` fixes the error, but what does Go think the return type is otherwise?

Here's the proposal: <https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md> It's notably agnostic about the implementation. Unclear to me whether generic functions are implemented via erasure or by instantiating multiple versions of the function.

### Day 3

The distinction between chars and bytes is a bit annoying. I got tripped up by doing:

    int(str[pos])

which gives you the ASCII value of the byte at that pos, rather than parsing it as a number. So `int('0')` --> `48` rather than `0`.

Factoring out `mostCommonBitAtPos` and `filterByBit` functions made part two easier to think about. I'm starting to get to the point where I'm curious how a Go expert would do this.

### Day 2

Introduces some string parsing. In Rust I would have set up an enum to represent the movements, and in TypeScript I'd set up a union type. I guess I could have set up a `struct` in Go since the moves all have the same shape.

### Day 1

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

The Go debugger is called… Delve! <https://github.com/go-delve/delve>

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
