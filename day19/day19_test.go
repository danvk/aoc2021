package main

import (
	"reflect"
	"testing"
)

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
