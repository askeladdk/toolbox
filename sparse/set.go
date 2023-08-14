// Package sparse provides sparse set and map implementations
// that efficiently pack their data in contiguous memory.
// The space complexity is O(n) and the time complexity of most operations is O(1).
// Due to the particulars of the representation Set and Map can only store integer keys.
//
// The implementation is based on the paper An Efficient Representation for Sparse Sets:
// https://citeseerx.ist.psu.edu/doc/10.1.1.30.7319
package sparse

import "slices"

// Integer is a constraint that permits any integer type.
// Only integer types can be used as keys in sparse sets and maps.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Set represents a sparse set of integers.
type Set[E Integer] struct {
	dense  []E
	sparse []E
}

// NewSet returns a new Set given an initial capacity.
func NewSet[E Integer](capacity int) *Set[E] {
	if capacity < 1 {
		capacity = 1
	}
	sparse := make([]E, capacity)
	for i := range sparse {
		sparse[i] = ^E(0)
	}
	return &Set[E]{
		dense:  make([]E, 0, capacity),
		sparse: sparse,
	}
}

// Grow ensures that s has space for at least n members.
func (s *Set[E]) Grow(n int) {
	if z := n - len(s.sparse); z > 0 {
		oldn := len(s.sparse)
		s.sparse = slices.Grow(s.sparse, z)
		s.sparse = s.sparse[:cap(s.sparse)]
		for i := oldn; i < len(s.sparse); i++ {
			s.sparse[i] = ^E(0)
		}
	}
}

// Cap reports the capacity of s.
// The complexity is O(1).
func (s *Set[E]) Cap() int {
	return len(s.sparse)
}

// Len reports the numbers of members in s.
// The complexity is O(1).
func (s *Set[E]) Len() int {
	return len(s.dense)
}

// Empty reports whether the set is empty.
func (s *Set[E]) Empty() bool {
	return len(s.dense) == 0
}

// Reset clears the set.
func (s *Set[E]) Reset() {
	for i := range s.sparse {
		s.sparse[i] = ^E(0)
	}
	s.dense = s.dense[:0]
}

// Has reports whether k is a member of s, where 0 <= x < s.Cap().
// The complexity is O(1).
func (s *Set[E]) Has(k E) bool {
	a := s.sparse[k]
	return a < E(len(s.dense)) && s.dense[a] == k
}

// Add includes k in s, where 0 <= k < s.Cap().
// The complexity is O(1).
func (s *Set[E]) Add(k E) {
	s.add(k)
}

// Del removes k from s, where 0 <= k < s.Cap().
// The complexity is O(1).
func (s *Set[E]) Del(k E) {
	s.del(k)
}

// Keys returns a slice of all members of s in no particular order.
// It should only be read or the set becomes invalid.
// The complexity is O(1).
func (s *Set[E]) Keys() []E {
	return s.dense[:len(s.dense):len(s.dense)]
}

// Swap implements sort.Interface.
func (s *Set[E]) Swap(i, j int) {
	p := s.dense[i]
	q := s.dense[j]
	s.sparse[p], s.sparse[q] = s.sparse[q], s.sparse[p]
	s.dense[i], s.dense[j] = s.dense[j], s.dense[i]
}

// Less implements sort.Interface.
func (s *Set[E]) Less(i, j int) bool {
	return s.dense[i] < s.dense[j]
}

func (s *Set[E]) add(k E) bool {
	s.Grow(int(k) + 1)
	a := s.sparse[k]
	n := len(s.dense)
	if a >= E(n) || s.dense[a] != k {
		s.sparse[k] = E(n)
		s.dense = append(s.dense, k)
		return true
	}
	return false
}

func (s *Set[E]) del(k E) int {
	a := s.sparse[k]
	n := len(s.dense)
	if a < E(n) && s.dense[a] == k {
		e := s.dense[n-1]
		s.sparse[e] = a
		s.sparse[k] = ^E(0)
		s.dense[a] = e
		s.dense = s.dense[:n-1]
		return int(a)
	}
	return -1
}
