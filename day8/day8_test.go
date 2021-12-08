package main

import (
	"reflect"
	"testing"
)

func TestIntersect(t *testing.T) {
	a := []Digit{'a', 'b', 'd'}
	b := []Digit{'b', 'd', 'f'}
	out := Intersect(a, b)
	expected := []Digit{'b', 'd'}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Intersect(%#v, %#v)=%#v, want %#v", a, b, out, expected)
	}
}
