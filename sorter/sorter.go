package sorter

import "sort"

type Interface interface {
	sort.Interface
}
type sorter[T any] struct {
	items    *[]T
	lessFunc func(items []T, i, j int) bool
}

func NewSorter[T any](items *[]T, lessFunc func(items []T, i, j int) bool) Interface {
	return &sorter[T]{
		items:    items,
		lessFunc: lessFunc,
	}
}

// Len is the number of elements in the collection.
func (s *sorter[T]) Len() int {
	return len(*s.items)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//   - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (s *sorter[T]) Less(i int, j int) bool {
	return s.lessFunc(*s.items, i, j)
}

// Swap swaps the elements with indexes i and j.
func (s *sorter[T]) Swap(i int, j int) {
	items := *s.items
	items[i], items[j] = items[j], items[i]
	*s.items = items
}
