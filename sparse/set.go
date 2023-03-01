// Package sparse provides sparse set and map implementations
// that are efficiently pack their data in contiguous memory.
// The space complexity is O(n) and the time complexity of most operations is O(1).
package sparse

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// An Efficient Representation For Sparse Sets
// https://citeseerx.ist.psu.edu/doc/10.1.1.30.7319

// Set represents a sparse set of integers.
type Set[E constraints.Integer] struct {
	dense  []E
	sparse []E
}

// NewSet returns a new Set with a maximum capacity of cap.
func NewSet[E constraints.Integer](cap int) *Set[E] {
	if cap < 1 {
		cap = 1
	}
	sparse := make([]E, cap)
	for i := range sparse {
		sparse[i] = ^E(0)
	}
	return &Set[E]{
		dense:  make([]E, 0, cap),
		sparse: sparse,
	}
}

// Grow ensures that s has space for at least n members.
func (m *Set[E]) Grow(n int) {
	if z := n - len(m.sparse) + 1; z > 0 {
		oldn := len(m.sparse)
		m.sparse = slices.Grow(m.sparse, z)
		m.sparse = m.sparse[:cap(m.sparse)]
		for i := oldn; i < len(m.sparse); i++ {
			m.sparse[i] = ^E(0)
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

// Add includes k in s, where 0 <= x < s.Cap().
// The complexity is O(1).
func (s *Set[E]) Add(k E) {
	s.add(k)
}

// Del removes k from s, where 0 <= x < s.Cap().
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
func (m *Set[E]) Swap(i, j int) {
	p := m.dense[i]
	q := m.dense[j]
	m.sparse[p], m.sparse[q] = m.sparse[q], m.sparse[p]
	m.dense[i], m.dense[j] = m.dense[j], m.dense[i]
}

// Less implements sort.Interface.
func (s *Set[E]) Less(i, j int) bool {
	return s.dense[i] < s.dense[j]
}

func (s *Set[E]) add(k E) bool {
	s.Grow(int(k))
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
