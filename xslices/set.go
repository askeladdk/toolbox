package xslices

import "cmp"

// Difference copies all elements in s1 excluding those that are also in s2 to d.
// It is assumed that s1 and s2 are sorted.
// At most len(s1) elements are copied.
//
// Difference does not allocate memory. Rather, the result is limited by len(d),
// and the return value is the number of elements copied to d.
// The time complexity is O(n).
//
// A slice may be updated in-place by passing it as d:
//
//	s1 = s1[:xslices.Difference(s1, s1, s2)]
func Difference[S ~[]E, E cmp.Ordered](d, s1, s2 S) int {
	i, j, k := 0, 0, 0
	for i < len(s1) && j < len(s2) && k < len(d) {
		v1, v2 := s1[i], s2[j]
		switch {
		case v1 < v2:
			d[k] = v1
			i++
			k++
		case v1 > v2:
			j++
		default:
			i++
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	return k
}

// DifferenceFunc is like [Difference] but uses a comparison function.
func DifferenceFunc[S ~[]E, E any](d, s1, s2 S, cmp func(E, E) int) int {
	i, j, k := 0, 0, 0
	for i < len(s1) && j < len(s2) && k < len(d) {
		v1, v2 := s1[i], s2[j]
		c := cmp(v1, v2)
		switch {
		case c < 0:
			d[k] = v1
			i++
			k++
		case c > 0:
			j++
		default:
			i++
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	return k
}

// Includes determines whether all elements in the subset slice
// are also elements of the superset slice.
// It is assumed that superset and subset are sorted.
// The time complexity is O(n).
func Includes[S ~[]E, E cmp.Ordered](superset, subset S) bool {
	j := 0
	for i := range subset {
		for j < len(superset) && subset[i] > superset[j] {
			j++
		}
		if j == len(superset) || subset[i] != superset[j] {
			return false
		}
	}
	return true
}

// IncludesFunc is like [Includes] but uses a comparison function.
func IncludesFunc[S ~[]E, E any](superset, subset S, cmp func(E, E) int) bool {
	j := 0
	for i := range subset {
		for j < len(superset) && cmp(subset[i], superset[j]) > 0 {
			j++
		}
		if j == len(superset) || cmp(subset[i], superset[j]) != 0 {
			return false
		}
	}
	return true
}

// Intersect copies all elements that appear in both s1 and s2 to d.
// It is assumed that s1 and s2 are sorted.
// At most min(len(s1), len(s2)) elements are copied.
//
// Intersect does not allocate memory. Rather, the result is limited by len(d),
// and the return value is the number of elements copied to d.
// The time complexity is O(n).
//
// A slice may be updated in-place by passing it as d:
//
//	s1 = s1[:Intersect(s1, s1, s2)]
func Intersect[S ~[]E, E cmp.Ordered](d, s1, s2 S) int {
	i, j, k := 0, 0, 0
	for i < len(s1) && j < len(s2) && k < len(d) {
		v1, v2 := s1[i], s2[j]
		switch {
		case v1 < v2:
			i++
		case v1 > v2:
			j++
		default:
			d[k] = v1
			i++
			j++
			k++
		}
	}
	return k
}

// IntersectFunc is like [Intersect] but uses a comparison function.
func IntersectFunc[S ~[]E, E any](d, s1, s2 S, cmp func(E, E) int) int {
	i, j, k := 0, 0, 0
	for i < len(s1) && j < len(s2) && k < len(d) {
		v1, v2 := s1[i], s2[j]
		c := cmp(v1, v2)
		switch {
		case c < 0:
			i++
		case c > 0:
			j++
		default:
			d[k] = v1
			i++
			j++
			k++
		}
	}
	return k
}

// Merge copies all elements in s1 and s2 in sequential order to d.
// It is assumed that s1 and s2 are sorted.
// The difference with Union is that elements that appear in both slices
// are copied twice.
// Exactly len(s1)+len(s2) elements are copied.
//
// Merge does not allocate memory. Rather, the result is limited by len(d),
// and the return value is the number of elements copied to d.
// The time complexity is O(n).
func Merge[S ~[]E, E cmp.Ordered](d, s1, s2 S) int {
	i, j, k := 0, 0, 0
	for ; i < len(s1) && j < len(s2) && k < len(d); k++ {
		v1, v2 := s1[i], s2[j]
		if v1 < v2 {
			d[k] = v1
			i++
		} else {
			d[k] = v2
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	k += copy(d[k:], s2[j:])
	return k
}

// MergeFunc is like [Merge] but uses a comparison function.
func MergeFunc[S ~[]E, E any](d, s1, s2 S, cmp func(E, E) int) int {
	i, j, k := 0, 0, 0
	for ; i < len(s1) && j < len(s2) && k < len(d); k++ {
		v1, v2 := s1[i], s2[j]
		if cmp(v1, v2) < 0 {
			d[k] = v1
			i++
		} else {
			d[k] = v2
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	k += copy(d[k:], s2[j:])
	return k
}

// SymmetricDifference copies all elements that are in s1 or s2 but not both to d.
// It is assumed that s1 and s2 are sorted.
// At most len(s1)+len(s2) elements are copied.
//
// SymmetricDifference does not allocate memory. Rather, the result is limited by len(d),
// and the return value is the number of elements copied to d.
// The time complexity is O(n).
func SymmetricDifference[S ~[]E, E cmp.Ordered](d, s1, s2 S) int {
	i, j, k := 0, 0, 0
	for i < len(s1) && j < len(s2) && k < len(d) {
		v1, v2 := s1[i], s2[j]
		switch {
		case v1 < v2:
			d[k] = v1
			i++
			k++
		case v1 > v2:
			d[k] = v2
			j++
			k++
		default:
			i++
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	k += copy(d[k:], s2[j:])
	return k
}

// SymmetricDifferenceFunc is like [SymmetricDifference] but uses a comparison function.
func SymmetricDifferenceFunc[S ~[]E, E any](d, s1, s2 S, cmp func(E, E) int) int {
	i, j, k := 0, 0, 0
	for i < len(s1) && j < len(s2) && k < len(d) {
		v1, v2 := s1[i], s2[j]
		c := cmp(v1, v2)
		switch {
		case c < 0:
			d[k] = v1
			i++
			k++
		case c > 0:
			d[k] = v2
			j++
			k++
		default:
			i++
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	k += copy(d[k:], s2[j:])
	return k
}

// Union copies all elements in s1 and s2 in sequential order to d.
// It is assumed that s1 and s2 are sorted.
// The difference with Merge is that elements that appear in both slices
// are copied only once.
// At most len(s1)+len(s2) elements are copied.
//
// Union does not allocate memory. Rather, the result is limited by len(d),
// and the return value is the number of elements copied to d.
// The time complexity is O(n).
func Union[S ~[]E, E cmp.Ordered](d, s1, s2 S) int {
	i, j, k := 0, 0, 0
	for ; i < len(s1) && j < len(s2) && k < len(d); k++ {
		v1, v2 := s1[i], s2[j]
		switch {
		case v1 < v2:
			d[k] = v1
			i++
		case v1 > v2:
			d[k] = v2
			j++
		default:
			d[k] = v1
			i++
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	k += copy(d[k:], s2[j:])
	return k
}

// UnionFunc is like [Union] but uses a comparison function.
func UnionFunc[S ~[]E, E any](d, s1, s2 S, cmp func(E, E) int) int {
	i, j, k := 0, 0, 0
	for ; i < len(s1) && j < len(s2) && k < len(d); k++ {
		v1, v2 := s1[i], s2[j]
		c := cmp(v1, v2)
		switch {
		case c < 0:
			d[k] = v1
			i++
		case c > 0:
			d[k] = v2
			j++
		default:
			d[k] = v1
			i++
			j++
		}
	}
	k += copy(d[k:], s1[i:])
	k += copy(d[k:], s2[j:])
	return k
}
