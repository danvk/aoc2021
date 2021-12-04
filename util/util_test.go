package util

import (
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
