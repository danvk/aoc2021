package main

import "testing"

func TestFuelForMove(t *testing.T) {
	ins := [...]int{0, 1, 2, 3, 4}
	outs := [...]int{0, 1, 3, 6, 10}

	for i, in := range ins {
		expected := outs[i]
		f := FuelForMove(in)
		if f != expected {
			t.Fatalf("Expected fuel(%d) = %d, got %d", in, expected, f)
		}
	}
}
