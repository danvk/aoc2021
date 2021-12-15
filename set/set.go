package set

import (
	"fmt"
	"strings"
)

type Set[T comparable] map[T]bool

// Would be nice if this method required Stringer | string
func (this Set[T]) String() string {
	keys := make([]string, 0, len(this))
	for k := range this {
		fmt.Printf("String %#v\n", k)
		var ki interface{} = k
		s, ok := ki.(string)
		if !ok {
			ss, ok := ki.(fmt.Stringer)
			if !ok {
				panic(k)
			}
			s = ss.String()
		}

		keys = append(keys, s)
	}
	return strings.Join(keys, ",")
}

func SetFrom[T comparable](els []T) Set[T] {
	result := make(Set[T])
	for _, el := range els {
		result[el] = true
	}
	return result
}

func SetFromChars(s string) Set[string] {
	result := make(Set[string])
	for _, c := range s {
		result[string(c)] = true
	}
	return result
}

func (this Set[T]) Eq(other Set[T]) bool {
	if len(this) != len(other) {
		return false
	}
	for el := range this {
		if _, ok := other[el]; !ok {
			return false
		}
	}
	return true
}

func (this *Set[T]) Union(other Set[T]) {
	for v := range other {
		(*this)[v] = true
	}
}

func (this *Set[T]) Clone() Set[T] {
	result := make(Set[T])
	result.Union(*this)
	return result
}

func (this Set[T]) UnionWith(other Set[T]) Set[T] {
	result := this.Clone()
	result.Union(other)
	return result
}

func (this Set[T]) Intersect(other Set[T]) {
	for k := range this {
		if _, ok := other[k]; !ok {
			delete(this, k)
		}
	}
}

func (this Set[T]) IntersectWith(other Set[T]) Set[T] {
	if len(this) < len(other) {
		return other.IntersectWith(this)
	}
	result := make(Set[T])
	for v := range other {
		if _, ok := this[v]; ok {
			result[v] = true
		}
	}
	return result
}

func (this Set[T]) Subtract(other Set[T]) {
	for v := range other {
		delete(this, v)
	}
}

func (this Set[T]) LoneElement() T {
	if len(this) == 1 {
		for k := range this {
			return k
		}
	}
	panic(this)
}

func (this Set[T]) Add(el T) {
	this[el] = true
}

func (this Set[T]) CloneAndAdd(el T) Set[T] {
	n := this.Clone()
	n.Add(el)
	return n
}

func (this Set[T]) Remove(el T) {
	delete(this, el)
}

func (this Set[T]) AsSlice() []T {
	keys := make([]T, 0, len(this))
	for k := range this {
		keys = append(keys, k)
	}
	return keys
}

func (this Set[T]) Contains(el T) bool {
	_, ok := this[el]
	return ok
}
