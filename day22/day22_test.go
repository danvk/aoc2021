package main

import (
	"reflect"
	"testing"
)

func TestSubtractInterval(t *testing.T) {
	tests := map[string]struct {
		a    Interval
		b    Interval
		want []Interval
	}{
		"overlap": {
			a:    Interval{10, 20},
			b:    Interval{15, 25},
			want: []Interval{{10, 14}},
		},
		"disjoint": {
			a:    Interval{10, 20},
			b:    Interval{25, 35},
			want: []Interval{{10, 20}},
		},
		"overlapB": {
			a:    Interval{15, 25},
			b:    Interval{10, 20},
			want: []Interval{{21, 25}},
		},
		"inside": {
			a:    Interval{5, 25},
			b:    Interval{10, 20},
			want: []Interval{{5, 9}, {21, 25}},
		},
		"outside": {
			a:    Interval{10, 20},
			b:    Interval{5, 25},
			want: []Interval{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Subtract(tc.b)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("%v.Subtract(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := map[string]struct {
		a    Interval
		b    Interval
		want Interval
	}{
		"overlap": {
			a:    Interval{10, 100},
			b:    Interval{50, 120},
			want: Interval{50, 100},
		},
		"empty": {
			a:    Interval{10, 100},
			b:    Interval{-100, 0},
			want: Interval{10, 0},
		},
		"touching": {
			a:    Interval{10, 100},
			b:    Interval{100, 110},
			want: Interval{100, 100},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Intersect(tc.b)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("%v.Intersect(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestIntersects(t *testing.T) {
	tests := map[string]struct {
		a    Interval
		b    Interval
		want bool
	}{
		"overlap": {
			a:    Interval{10, 100},
			b:    Interval{50, 120},
			want: true,
		},
		"empty": {
			a:    Interval{10, 100},
			b:    Interval{-100, 0},
			want: false,
		},
		"touching": {
			a:    Interval{10, 100},
			b:    Interval{100, 110},
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Intersects(tc.b)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("%v.Intersect(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
