package util

import (
	"reflect"
	"testing"
)

func TestAllTrue(t *testing.T) {
	if !AllTrue([]bool{true, true, true}) {
		t.Log("Should be true")
		t.Fail()
	}

	if AllTrue([]bool{true, false, true}) {
		t.Log("Should be false")
		t.Fail()
	}

	if !AllTrue([]bool{}) {
		t.Log("Empty array is all true")
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
