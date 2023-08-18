package xslices_test

import (
	"fmt"

	"github.com/askeladdk/toolbox/xslices"
)

func ExampleDifference() {
	s1 := []int{0, 1, 1, 3, 4, 5, 6, 7}
	s2 := []int{1, 4, 6}
	d := make([]int, len(s1))
	i := xslices.Difference(d, s1, s2)
	fmt.Println(d[:i])
	// Output:
	// [0 1 3 5 7]
}

func ExampleIncludes() {
	superset := []int{0, 1, 1, 3, 4, 5, 6, 7}
	subset := []int{1, 4, 6}
	fmt.Println(xslices.Includes(superset, subset))
	// Output:
	// true
}

func ExampleIntersect() {
	s1 := []int{0, 1, 1, 3, 4, 5, 6, 7}
	s2 := []int{1, 2, 4, 8}
	d := make([]int, min(len(s1), len(s2)))
	i := xslices.Intersect(d, s1, s2)
	fmt.Println(d[:i])
	// Output:
	// [1 4]
}

func ExampleMerge() {
	s1 := []int{0, 1, 1, 3, 4, 5, 6, 7}
	s2 := []int{1, 2, 4, 8}
	d := make([]int, len(s1)+len(s2))
	i := xslices.Merge(d, s1, s2)
	fmt.Println(d[:i])
	// Output:
	// [0 1 1 1 2 3 4 4 5 6 7 8]
}

func ExampleSymmetricDifference() {
	s1 := []int{0, 1, 1, 3, 4, 5, 6, 7}
	s2 := []int{1, 4, 6, 7, 8}
	d := make([]int, len(s1)+len(s2))
	i := xslices.SymmetricDifference(d, s1, s2)
	fmt.Println(d[:i])
	// Output:
	// [0 1 3 5 8]
}

func ExampleUnion() {
	s1 := []int{0, 1, 1, 3, 4, 5, 6, 7}
	s2 := []int{1, 2, 4, 8}
	d := make([]int, len(s1)+len(s2))
	i := xslices.Union(d, s1, s2)
	fmt.Println(d[:i])
	// Output:
	// [0 1 1 2 3 4 5 6 7 8]
}
