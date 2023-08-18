package xslices_test

import (
	"fmt"

	"github.com/askeladdk/toolbox/xslices"
)

func ExamplePartition() {
	xs := []int{8, 2, 3, 0, 5, 7, 1, 4, 6, 9}
	i := xslices.Partition(xs, 4)
	left, right := xs[:i], xs[i:]
	fmt.Println(left)
	fmt.Println(right)
	// Output:
	// [4 2 3 0 1]
	// [7 5 8 6 9]
}

func ExamplePartitionFunc() {
	iseven := func(a int) bool { return a%2 == 0 }
	xs := []int{8, 2, 3, 0, 5, 7, 1, 4, 6, 9}
	i := xslices.PartitionFunc(xs, iseven)
	left, right := xs[:i], xs[i:]
	fmt.Println(left)
	fmt.Println(right)
	// Output:
	// [8 2 6 0 4]
	// [7 1 5 3 9]
}

func ExampleSelect() {
	xs := []int{8, 2, 3, 0, 5, 7, 1, 4, 6, 9}
	v := xslices.Select(xs, 8)
	fmt.Println(v)
	fmt.Println(xs)
	// Output:
	// 8
	// [6 2 3 0 5 7 1 4 8 9]
}

func ExampleSelectFunc() {
	descending := func(a, b int) int { return b - a }
	xs := []int{8, 2, 3, 0, 5, 7, 1, 4, 6, 9}
	v := xslices.SelectFunc(xs, 1, descending)
	fmt.Println(v)
	fmt.Println(xs)
	// Output:
	// 8
	// [9 8 3 0 5 7 1 4 6 2]
}
