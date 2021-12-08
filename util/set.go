package util

type Set[T comparable] map[T]bool

func SetFrom[T comparable](els []T) Set[T] {
	result := make(Set[T])
	for _, el := range els {
		result[el] = true
	}
	return result
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

func (this *Set[T]) UnionWith(other Set[T]) Set[T] {
	result := this.Clone()
	result.Union(other)
	return result
}

func (this *Set[T]) Intersect(other Set[T]) Set[T] {
	if len(*this) < len(other) {
		return other.Intersect(*this)
	}
	result := make(Set[T])
	for v := range other {
		if _, ok := (*this)[v]; ok {
			result[v] = true
		}
	}
	return result
}
