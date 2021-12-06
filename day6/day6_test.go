package main

import (
	"reflect"
	"testing"
)

func TestAdvance(t *testing.T) {
	school := []Lanternfish{{timer: 6}}

	Advance(&school)
	expect := []Lanternfish{{timer: 5}}
	if !reflect.DeepEqual(expect, school) {
		t.Errorf("Expected %v, got %v", expect, school)
		t.Fail()
	}
}
