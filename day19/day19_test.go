package main

import (
	"reflect"
	"testing"
)

func TestFindAllOrientations(t *testing.T) {
	rots := FindAllOrientations()
	if len(rots) != 24 {
		t.Errorf("Got %d rotations, want %d", len(rots), 24)
	}
}

func TestMult(t *testing.T) {
	x := Point{5, 0, 0}
	want := Point{0, 5, 0}
	got := x.Mult(ROTS[0])
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%#v.Mult(ROTS[0]) = %#v want %#v", x, got, want)
	}

	want = Point{-5, 0, 0}
	got = x.Mult(ROTS[0]).Mult(ROTS[0])
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%#v.Mult(ROTS[0]).Mult(ROTS[0]) = %#v want %#v", x, got, want)
	}

	m := ROTS[0].Mult(ROTS[0])
	got = x.Mult(m)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%#v.Mult(ROTS[0]).Mult(ROTS[0]) = %#v want %#v", x, got, want)
	}
}

func TestRot90Z(t *testing.T) {
	tests := map[string]struct {
		p    Point
		want Point
	}{
		"x": {p: Point{1, 0, 0}, want: Point{0, 1, 0}},
		"y": {p: Point{0, 1, 0}, want: Point{-1, 0, 0}},
		"z": {p: Point{-1, 0, 10}, want: Point{0, -1, 10}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.p.Rot90Z()
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("%v = %#v, want %#v", tc, got, tc.want)
			}
		})
	}
}

func TestFindBestOverlap(t *testing.T) {
	tests := map[string]struct {
		a     []Point
		b     []Point
		wantN int
		wantP Point
	}{
		"sample": {
			a: []Point{
				{0, 2, 0},
				{4, 1, 0},
				{3, 3, 0},
			},
			b: []Point{
				{-1, -1, 0},
				{-5, 0, 0},
				{-2, 1, 0},
			},
			wantN: 3,
			wantP: Point{-5, -2, 0},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotP, gotN := FindBestOverlap(tc.a, tc.b)
			if !reflect.DeepEqual(tc.wantP, gotP) || tc.wantN != gotN {
				t.Errorf("%v = %#v, %#v, want %#v, %#v", tc, gotP, gotN, tc.wantP, tc.wantN)
			}
		})
	}
}

func TestNumOverlapping(t *testing.T) {
	tests := map[string]struct {
		a    []Point
		b    []Point
		want int
	}{
		"some": {a: []Point{{0, 0, 0}, {1, 0, 0}}, b: []Point{{1, 0, 0}, {1, 1, 0}}, want: 1},
		"zero": {a: []Point{}, b: []Point{{1, 0, 0}, {1, 1, 0}}, want: 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NumOverlapping(tc.a, tc.b)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("%v = %#v, want %#v", tc, got, tc.want)
			}
		})
	}
}

func add(a, b int) int {
	return a + b
}

func Test(t *testing.T) {
	tests := map[string]struct {
		a    int
		b    int
		want int
	}{
		"positive": {a: 10, b: 20, want: 30},
		"zero":     {a: 0, b: 1, want: 1},
		"negative": {a: -10, b: 10, want: 1},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := add(tc.a, tc.b)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("%v = %#v, want %#v", tc, got, tc.want)
			}
		})
	}
}
