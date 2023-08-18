package xslices

import "cmp"

// Partition divides s into two subslices and returns the partition index.
// The lower partition contains all elements in s which are less than or equal to target.
// The higher partition contains all elements which are strictly greater than target.
// The relative order of elements is not preserved.
// Partition modifies the contents of the slice s; it does not create a new slice.
func Partition[S ~[]E, E cmp.Ordered](s S, target E) int {
	for i, j := 0, len(s)-1; ; {
		for i < len(s) && s[i] <= target {
			i++
		}
		for j >= 0 && s[j] > target {
			j--
		}
		if i >= j {
			return i
		}
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
}

// PartitionFunc divides s into two subslices and returns the partition index.
// The lower partition contains all elements e in s for which f(e) is true.
// The higher partition contains all elements for which f(e) is false.
// The relative order of elements is not preserved.
// PartitionFunc modifies the contents of the slice s; it does not create a new slice.
func PartitionFunc[S ~[]E, E any](s S, f func(E) bool) int {
	for i, j := 0, len(s)-1; ; {
		for i < len(s) && f(s[i]) {
			i++
		}
		for j >= 0 && !f(s[j]) {
			j--
		}
		if i >= j {
			return i
		}
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
}

// Select finds the k-th smallest element in an unsorted slice.
// The slice is partially sorted in-place such that
// s[k] contains the k-th smallest element.
// The complexity is O(n) on average and O(n^2) in the worst case.
func Select[S ~[]E, E cmp.Ordered](s S, k int) E {
	_ = s[k] // bounds check

	lo, hi := 0, len(s)-1
loop:
	for lo < hi {
		mi := lo + (hi-lo)/2
		pivot := median3(s[lo], s[mi], s[hi])
		for i, j := lo, hi; ; {
			for i < len(s) && s[i] < pivot {
				i++
			}
			for j >= 0 && s[j] > pivot {
				j--
			}
			if i >= j {
				switch {
				case k < i:
					hi = i - 1
					continue loop
				case k > i:
					lo = i + 1
					continue loop
				default:
					return s[i]
				}
			}
			s[i], s[j] = s[j], s[i]
		}
	}
	return s[lo]
}

// SelectFunc is like [Select] but compares elements using a comparison function.
func SelectFunc[S ~[]E, E any](s S, k int, cmp func(E, E) int) E {
	_ = s[k] // bounds check

	lo, hi := 0, len(s)-1
loop:
	for lo < hi {
		mi := lo + (hi-lo)/2
		pivot := median3Func(s[lo], s[mi], s[hi], cmp)
		for i, j := lo, hi; ; {
			for i < len(s) && cmp(s[i], pivot) < 0 {
				i++
			}
			for j >= 0 && cmp(s[j], pivot) > 0 {
				j--
			}
			if i >= j {
				switch {
				case k < i:
					hi = i - 1
					continue loop
				case k > i:
					lo = i + 1
					continue loop
				default:
					return s[i]
				}
			}
			s[i], s[j] = s[j], s[i]
		}
	}
	return s[lo]
}

func median3[E cmp.Ordered](a, b, c E) E {
	if a > b { // a, b = min(a, b), max(a, b)
		a, b = b, a
	}
	b = min(b, c)
	a = max(a, b)
	return a
}

func median3Func[E any](a, b, c E, cmp func(E, E) int) E {
	if cmp(a, b) > 0 { // a, b = min(a, b), max(a, b)
		a, b = b, a
	}
	if cmp(b, c) > 0 { // b = min(b, c)
		b = c
	}
	if cmp(a, b) < 0 { // a = max(a, b)
		a = b
	}
	return a
}
