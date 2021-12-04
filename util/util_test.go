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

func TestMap(t *testing.T) {
	nums := Map([]int{1, 2, 3}, func(x int) int { return x * x })
	if !reflect.DeepEqual(nums, []int{1, 4, 9}) {
		t.Fatalf("Expected 1, 4, 9 got %v", nums)
	}

	nums = Map([]string{"hi", "bye"}, func(x string) int { return len(x) })
	if !reflect.DeepEqual(nums, []int{2, 3}) {
		t.Fatalf("Expected 2, 3 got %v", nums)
	}
}

func TestFilter(t *testing.T) {
	nums := Filter([]int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 == 1 })
	if !reflect.DeepEqual(nums, []int{1, 3, 5}) {
		t.Fatalf("Expected 1, 3, 5 got %v", nums)
	}
}

func TestFlatMap(t *testing.T) {
	nums := FlatMap([]int{1, 2, 3}, func(x int) []int {
		if x%2 == 0 {
			return []int{}
		}
		return []int{x, 2 * x}
	})
	if !reflect.DeepEqual(nums, []int{1, 2, 3, 6}) {
		t.Fatalf("Expected 1, 2, 3, 6 got %v", nums)
	}
}
