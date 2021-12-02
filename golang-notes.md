# Go Lang Notes

## Advent of Code

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
