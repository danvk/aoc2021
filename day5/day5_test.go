package main

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	line := ParseLine("0,9 -> 5,9")
	if !reflect.DeepEqual(line, Line{start: Coord{x: 0, y: 9}, end: Coord{x: 5, y: 9}}) {
		t.Errorf("Expected 0,9 -> 5,9 got %v", line)
		t.Fail()
	}

	line = ParseLine("10,9 -> 15,29")
	if !reflect.DeepEqual(line, Line{start: Coord{x: 10, y: 9}, end: Coord{x: 15, y: 29}}) {
		t.Errorf("Expected 10,9 -> 15,29 got %v", line)
		t.Fail()
	}
}
