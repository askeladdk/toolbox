package sparse

import "golang.org/x/exp/constraints"

// Map represents a sparse map of integer keys
// and arbitrarily typed values.
type Map[E constraints.Integer, T any] struct {
	set    Set[E]
	values []T
}

// NewMap returns a new Map given an initial capacity.
func NewMap[E constraints.Integer, T any](capacity int) *Map[E, T] {
	if capacity < 1 {
		capacity = 1
	}
	return &Map[E, T]{
		set: Set[E]{
			dense:  make([]E, 0, capacity),
			sparse: make([]E, capacity),
		},
		values: make([]T, 0, capacity),
	}
}

// Grow ensures that m has space for at least n members.
func (m *Map[E, T]) Grow(n int) {
	m.set.Grow(n)
}

// Cap reports the capacity of m.
// The complexity is O(1).
func (m *Map[E, T]) Cap() int {
	return m.set.Cap()
}

// Len reports the numbers of members in m.
// The complexity is O(1).
func (m *Map[E, T]) Len() int {
	return m.set.Len()
}

// Empty reports whether the map is empty.
func (m *Map[E, T]) Empty() bool {
	return m.set.Empty()
}

// Reset clears the map.
func (m *Map[E, T]) Reset() {
	var zeroval T
	for i := range m.values {
		m.values[i] = zeroval
	}
	m.values = m.values[:0]
	m.set.Reset()
}

// Has reports whether k is a member of s, where 0 <= x < s.Cap().
// The complexity is O(1).
func (m *Map[E, T]) Has(k E) bool {
	return m.set.Has(k)
}

// Set maps k to v in m, where 0 <= k < s.Cap().
// The complexity is O(1).
func (m *Map[E, T]) Set(k E, v T) {
	if m.set.add(k) {
		m.values = append(m.values, v)
	}
}

// Del removes x from s, where 0 <= x < s.Cap().
// The complexity is O(1).
func (m *Map[E, T]) Del(x E) {
	if i := m.set.del(x); i != -1 {
		var zeroval T
		n := len(m.values) - 1
		m.values[i] = m.values[n]
		m.values[n] = zeroval
		m.values = m.values[:n]
	}
}

// Get returns the value of k if it exists.
// The complexity is O(1).
func (m *Map[E, T]) Get(k E) (v T, exists bool) {
	var i int
	if i, exists = m.Index(k); exists {
		v = m.values[i]
	}
	return
}

// Swap implements sort.Interface.
func (m *Map[E, T]) Swap(i, j int) {
	m.set.Swap(i, j)
	m.values[i], m.values[j] = m.values[j], m.values[i]
}

// Less implements sort.Interface.
func (m *Map[E, T]) Less(i, j int) bool {
	return m.set.Less(i, j)
}

// Keys returns a slice of all keys of m in no particular order.
// It should only be read or the map becomes invalid.
// The complexity is O(1).
func (m *Map[E, T]) Keys() []E {
	return m.set.Keys()
}

// Values returns a slice of all values of m in no particular order.
// The i-th element of Members() maps to the i-th element of Values().
// The complexity is O(1).
func (m *Map[E, T]) Values() []T {
	return m.values[:len(m.values):len(m.values)]
}

// Index returns the index of the value of k in Values() if it exists.
// The complexity is O(1).
func (m *Map[E, T]) Index(k E) (int, bool) {
	a := int(m.set.sparse[k])
	return a, a < len(m.set.dense) && m.set.dense[a] == k
}
