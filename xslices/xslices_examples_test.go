package xslices_test

import (
	"fmt"

	"github.com/askeladdk/toolbox/xslices"
)

func ExampleCount() {
	xs := []int{1, 3, 0, 1, 4, 1, 2, 2, 0, 1}
	fmt.Println(xslices.Count(xs, 1))
	// Output:
	// 4
}

func ExampleCountFunc() {
	iseven := func(i int) bool { return i%2 == 0 }
	xs := []int{8, 2, 3, 0, 5, 7, 1, 4, 6, 9}
	fmt.Println(xslices.CountFunc(xs, iseven))
	// Output:
	// 5
}

func ExampleGroup() {
	xs := []int{1, 1, 2, 2, 2, 3}
	ys := xslices.Group(nil, xs)
	fmt.Println(ys)
	// Output:
	// [[1 1] [2 2 2] [3]]
}

func ExampleGroupFunc() {
	eq := func(a, b int) bool { return a == b }
	xs := []int{1, 1, 2, 2, 2, 3}
	ys := xslices.GroupFunc(nil, xs, eq)
	fmt.Println(ys)
	// Output:
	// [[1 1] [2 2 2] [3]]
}

func ExamplePermute() {
	a := []string{"e", "d", "c", "a", "b"}
	p := []int{3, 4, 2, 1, 0}
	xslices.Permute(a, p)
	fmt.Println(a)
	// Output:
	// [a b c d e]
}
