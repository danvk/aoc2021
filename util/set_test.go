package util

import (
	"reflect"
	"testing"
)

func TestUnionWith(t *testing.T) {
	a := SetFrom([]string{"a", "b", "c"})
	b := SetFrom([]string{"b", "c", "d", "e"})
	expected := SetFrom([]string{"a", "b", "c", "d", "e"})
	actual := a.UnionWith(b)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%#v.UnionWith(%#v) = %#v want %#v", a, b, actual, expected)
	}
}
