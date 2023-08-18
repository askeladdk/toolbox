package xslices

import (
	"fmt"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestDifference(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	for i, tt := range []struct {
		A []int
		B []int
		R []int
	}{
		{
			R: []int{},
		},
		{
			A: []int{1},
			R: []int{1},
		},
		{
			A: []int{1, 3, 3, 3, 4, 6, 6, 7},
			B: []int{2, 2, 3, 3, 5, 5, 7, 7, 8, 8},
			R: []int{1, 3, 4, 6, 6},
		},
		{
			A: []int{1, 3, 4, 5, 6, 6, 6},
			B: []int{2, 4, 6},
			R: []int{1, 3, 5, 6, 6},
		},
		{
			A: []int{1, 2, 3},
			B: []int{4, 5, 6},
			R: []int{1, 2, 3},
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			xs := make([]int, 100)
			require.Equal(t, tt.R, xs[:Difference(xs, tt.A, tt.B)])
			clear(xs)
			require.Equal(t, tt.R, xs[:DifferenceFunc(xs, tt.A, tt.B, cmp)])
		})
	}
}

func TestIncludes(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	for i, tt := range []struct {
		Superset []int
		Subset   []int
		Expected bool
	}{
		{
			Expected: true,
		},
		{
			Superset: []int{1},
			Expected: true,
		},
		{
			Subset:   []int{1},
			Expected: false,
		},
		{
			Superset: []int{1, 2, 3, 4},
			Subset:   []int{1, 2, 3},
			Expected: true,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require.Equal(t, tt.Expected, Includes(tt.Superset, tt.Subset))
			require.Equal(t, tt.Expected, IncludesFunc(tt.Superset, tt.Subset, cmp))
		})
	}
}

func TestIntersect(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	for i, tt := range []struct {
		A []int
		B []int
		R []int
	}{
		{
			R: []int{},
		},
		{
			A: []int{1},
			R: []int{},
		},
		{
			A: []int{1, 3, 3, 3, 4, 6, 6, 7},
			B: []int{2, 2, 3, 3, 5, 5, 7, 7, 8, 8},
			R: []int{3, 3, 7},
		},
		{
			A: []int{1, 3, 4, 5, 6, 6, 6},
			B: []int{2, 4, 6},
			R: []int{4, 6},
		},
		{
			A: []int{1, 2, 3},
			B: []int{4, 5, 6},
			R: []int{},
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			xs := make([]int, 100)
			require.Equal(t, tt.R, xs[:Intersect(xs, tt.A, tt.B)])
			clear(xs)
			require.Equal(t, tt.R, xs[:IntersectFunc(xs, tt.A, tt.B, cmp)])
		})
	}
}

func TestMerge(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	for i, tt := range []struct {
		A []int
		B []int
		R []int
	}{
		{
			R: []int{},
		},
		{
			A: []int{1},
			R: []int{1},
		},
		{
			B: []int{1},
			R: []int{1},
		},
		{
			A: []int{1, 3, 3, 3, 4, 6, 6, 7},
			B: []int{2, 2, 3, 3, 5, 5, 7, 7, 8, 8},
			R: []int{1, 2, 2, 3, 3, 3, 3, 3, 4, 5, 5, 6, 6, 7, 7, 7, 8, 8},
		},
		{
			A: []int{1, 3, 4, 5, 6, 6, 6},
			B: []int{2, 4, 6},
			R: []int{1, 2, 3, 4, 4, 5, 6, 6, 6, 6},
		},
		{
			A: []int{1, 2, 3},
			B: []int{4, 5, 6},
			R: []int{1, 2, 3, 4, 5, 6},
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			xs := make([]int, 100)
			require.Equal(t, tt.R, xs[:Merge(xs, tt.A, tt.B)])
			clear(xs)
			require.Equal(t, tt.R, xs[:MergeFunc(xs, tt.A, tt.B, cmp)])
		})
	}
}

func TestSymmetricDifference(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	for i, tt := range []struct {
		A []int
		B []int
		R []int
	}{
		{
			R: []int{},
		},
		{
			A: []int{1},
			R: []int{1},
		},
		{
			A: []int{1, 3, 3, 3, 4, 6, 6, 7},
			B: []int{2, 2, 3, 3, 5, 5, 7, 7, 8, 8},
			R: []int{1, 2, 2, 3, 4, 5, 5, 6, 6, 7, 8, 8},
		},
		{
			A: []int{1, 3, 4, 5, 6, 6, 6},
			B: []int{2, 4, 6},
			R: []int{1, 2, 3, 5, 6, 6},
		},
		{
			A: []int{1, 2, 3},
			B: []int{4, 5, 6},
			R: []int{1, 2, 3, 4, 5, 6},
		},
		{
			A: []int{3, 4},
			B: []int{1, 2, 3},
			R: []int{1, 2, 4},
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			xs := make([]int, 100)
			require.Equal(t, tt.R, xs[:SymmetricDifference(xs, tt.A, tt.B)])
			clear(xs)
			require.Equal(t, tt.R, xs[:SymmetricDifferenceFunc(xs, tt.A, tt.B, cmp)])
		})
	}
}

func TestUnion(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	for i, tt := range []struct {
		A []int
		B []int
		R []int
	}{
		{
			R: []int{},
		},
		{
			A: []int{1},
			R: []int{1},
		},
		{
			B: []int{1},
			R: []int{1},
		},
		{
			A: []int{1, 3, 3, 3, 4, 6, 6, 7},
			B: []int{2, 2, 3, 3, 5, 5, 7, 7, 8, 8},
			R: []int{1, 2, 2, 3, 3, 3, 4, 5, 5, 6, 6, 7, 7, 8, 8},
		},
		{
			A: []int{1, 3, 4, 5, 6, 6, 6},
			B: []int{2, 4, 6},
			R: []int{1, 2, 3, 4, 5, 6, 6, 6},
		},
		{
			A: []int{1, 2, 3},
			B: []int{4, 5, 6},
			R: []int{1, 2, 3, 4, 5, 6},
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			xs := make([]int, 100)
			require.Equal(t, tt.R, xs[:Union(xs, tt.A, tt.B)])
			clear(xs)
			require.Equal(t, tt.R, xs[:UnionFunc(xs, tt.A, tt.B, cmp)])
		})
	}
}
