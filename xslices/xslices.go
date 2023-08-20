// Package xslices defines algorithms that operate on slices of any type.
package xslices

// Count reports how many elements in s are equal to target.
func Count[S ~[]E, E comparable](s S, target E) int {
	n := 0
	for i := range s {
		if s[i] == target {
			n++
		}
	}
	return n
}

// CountFunc reports how many elements e in s satisfy f(e).
func CountFunc[S ~[]E, E any](s S, f func(E) bool) int {
	n := 0
	for i := range s {
		if f(s[i]) {
			n++
		}
	}
	return n
}

// Group collects consecutive equal elements into subslices of s
// and appends them to d.
// Every subslice contains at least one element,
// and the sum of the lengths of the subslices is equal to len(s).
// The underlying slice s is not modified but should be sorted beforehand
// to group all equal elements into a single subslice.
func Group[S ~[]E, E comparable](d []S, s S) []S {
	if len(s) != 0 {
		i := 0
		for j := 1; j < len(s); j++ {
			if s[j-1] != s[j] {
				d = append(d, s[i:j:j])
				i = j
			}
		}
		d = append(d, s[i:len(s):len(s)])
	}
	return d
}

// GroupFunc is like [Group] but uses eq to determine equality.
func GroupFunc[S ~[]E, E any](d []S, s S, eq func(E, E) bool) []S {
	if len(s) != 0 {
		i := 0
		for j := 1; j < len(s); j++ {
			if !eq(s[j-1], s[j]) {
				d = append(d, s[i:j:j])
				i = j
			}
		}
		d = append(d, s[i:len(s):len(s)])
	}
	return d
}

// Reorder changes the order of elements in s according to the indices in x,
// which must be a permutation of indices in the half-open interval [0, len(s)),
// Every index must appear exactly once or the result will be incorrect.
// Reorder panics if len(s) != len(x).
//
// The complexity is O(nlogn) on average and O(n^2) in the worse case.
func Reorder[S ~[]E, E any](s S, x []int) {
	if len(s) != len(x) {
		panic("xslices.Reorder: unequal slice lengths")
	}
	for i, j := range x {
		for n := len(x); n != 0 && j < i; n-- {
			j = x[j]
		}
		s[i], s[j] = s[j], s[i]
	}
}
