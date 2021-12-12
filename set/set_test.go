package set

import (
	"reflect"
	"testing"
)

func TestUnionWith(t *testing.T) {
	a := SetFromChars("abc")
	b := SetFromChars("bcde")
	expected := SetFromChars("abcde")
	actual := a.UnionWith(b)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%#v.UnionWith(%#v) = %#v want %#v", a, b, actual, expected)
	}
	if !a.Eq(SetFromChars("abc")) {
		t.Errorf("a was mutated by UnionWith, got %v want a,b,c", a)
	}
}

func TestUnion(t *testing.T) {
	a := SetFromChars("abc")
	b := SetFromChars("cde")
	a.Union(b)
	if !a.Eq(SetFromChars("abcde")) {
		t.Errorf("a.Union(%#v) = %s want a,b,c,d,e", a, b)
	}
}

func TestIntersect(t *testing.T) {
	a := SetFromChars("abc")
	b := SetFromChars("cde")
	a.Intersect(b)
	if !a.Eq(SetFromChars("c")) {
		t.Errorf("a.Intersect(%#v) = %s want c", a, b)
	}
}

func TestIntersectWith(t *testing.T) {
	a := SetFromChars("abc")
	b := SetFromChars("cde")
	actual := a.IntersectWith(b)
	if !actual.Eq(SetFromChars("c")) {
		t.Errorf("%#v.Intersect(%#v) = %s want c", a, b, actual)
	}
	if !a.Eq(SetFromChars("abc")) {
		t.Errorf("a was mutated to %#v want a,b,c", a)
	}
}

func TestSubtract(t *testing.T) {
	a := SetFromChars("abc")
	b := SetFromChars("cde")

	a.Subtract(b)
	if !a.Eq(SetFromChars("ab")) {
		t.Errorf("a.Subtract(%s) = %s want a,b", b, a)
	}
}
