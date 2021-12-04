package util

import (
	"reflect"
	"testing"
)

func TestAllEq(t *testing.T) {
	if !AllEq([]bool{true, true, true}, true) {
		t.Log("Should be true")
		t.Fail()
	}

	if AllEq([]bool{true, false, true}, true) {
		t.Log("Should be false")
		t.Fail()
	}

	if !AllEq([]bool{}, true) {
		t.Log("Empty array is all true")
		t.Fail()
	}

	if !AllEq([]int{42, 42, 42}, 42) {
		t.Log("Should be all 42")
		t.Fail()
	}
}

func TestParseLineAsNums(t *testing.T) {
	nums := ParseLineAsNums("1, 2, 3, 4, 5", ",", false)
	if !reflect.DeepEqual(nums, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("Expected [1, 2, 3, 4, 5], got %v", nums)
	}

	nums = ParseLineAsNums("1  2    3 4", " ", true)
	if !reflect.DeepEqual(nums, []int{1, 2, 3, 4}) {
		t.Fatalf("Expected [1, 2, 3, 4], got %v", nums)
	}
}
