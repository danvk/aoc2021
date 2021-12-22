package main

import (
	"reflect"
	"testing"
)

func add(a, b int) int {
	return a + b
}

func TestClip(t *testing.T) {
	tests := map[string]struct {
		a    string
		want Cuboid
	}{
		"outside": {
			a:    "on x=-54112..-39298,y=-85059..-49293,z=-27449..7877",
			want: Cuboid{},
		},
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
