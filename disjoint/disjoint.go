// Package disjoint implements operations on disjoint sets, also known as union find.
package disjoint

// Set is a compact representation of a disjoint set containing nodes
// numbering from 0 to Len()-1.
// Initially each node belongs to a unique group containing only that node.
//
// The nodes are encoded such that for each value x:
//
//   - x > 0 is the root of a group of size x;
//   - x < 0 is a link to node -x;
//   - x = 0 is a link to node 0.
type Set []int

// New creates a disjoint set containing n singleton groups.
func New(n int) Set {
	s := make(Set, n)
	s.Reset()
	return s
}

// Reset clears the set to have Len() singleton groups.
// The complexity is O(n).
func (s Set) Reset() {
	for i := range s {
		s[i] = 1
	}
}

// Find returns the root of i.
func (s Set) Find(i int) int {
	val := s[i]
	switch {
	case val > 0:
		return i
	case val == 0:
		if i == 0 {
			panic("disjoint: node 0 points to itself")
		}
		fallthrough
	default:
		root := s.Find(-val)
		if val != root { // path compression
			s[i] = -root
		}
		return root
	}
}

// Union merges the groups of i and j together
// whereby the smaller group is merged into the larger.
// Returns true if the groups were merged or false if
// i and j are already in the same group.
func (s Set) Union(i, j int) bool {
	p := s.Find(i)
	q := s.Find(j)

	if p == q {
		return false
	}

	// union by size
	if s[q] >= s[p] {
		p, q = q, p
	}

	s[p] += s[q]
	s[q] = -p
	return true
}

// Same reports whether nodes i and j are in the same group.
func (s Set) Same(i, j int) bool {
	return s.Find(i) == s.Find(j)
}

// Len reports the number of nodes in s.
func (s Set) Len() int {
	return len(s)
}

// Size reports the size of the group of i.
func (s Set) Size(i int) int {
	return s[s.Find(i)]
}

// CountGroups reports the number of groups in s.
func (s Set) CountGroups() int {
	n := 0
	for _, v := range s {
		if v > 0 {
			n++
		}
	}
	return n
}
