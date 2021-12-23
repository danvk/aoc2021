package main

import (
	"testing"
)

func TestEncodeDecodePos(t *testing.T) {
	tests := map[string]struct {
		x int
		y int
	}{
		"0": {x: 0, y: 0},
		"a": {x: 10, y: 0},
		"b": {x: 2, y: 1},
		"c": {x: 4, y: 2},
		"d": {x: 6, y: 3},
		"e": {x: 8, y: 4},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			enc := EncodePos(tc.x, tc.y)
			dX, dY := DecodePos(enc)
			if dX != tc.x || dY != tc.y || enc >= 27 {
				t.Errorf("%v -> %d -> (%d, %d)", tc, enc, dX, dY)
			}
		})
	}
}

func TestEncodeFinal(t *testing.T) {
	final := `#############
#...........#
###A#B#C#D###
###A#B#C#D###
###A#B#C#D###
  #A#B#C#D#
  #########`
	state := ParseState(final)

	n1 := state.Encode()

	state.amphipods[1], state.amphipods[0] = state.amphipods[0], state.amphipods[1]
	n2 := state.Encode()

	if n1 != n2 {
		t.Errorf("Multiple final state encodings: %d, %d", n1, n2)
	}
}

func TestEncodeAndBack(t *testing.T) {
	input := `#############
#D........A.#
###.#C#.#B###
###B#C#B#B###
###B#C#B#A###
  #A#D#C#A#
  #########`
	state := ParseState(input)

	n := state.Encode()
	back := Decode(n)

	if state.String() != back.String() {
		t.Errorf("Mismatch:\n%s\nvs:\n%s", state, back)
	}
}
