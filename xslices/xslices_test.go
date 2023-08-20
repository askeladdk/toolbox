package xslices

import (
	"math/rand"
	"slices"
	"strconv"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func all[S ~[]E, E any](s S, f func(E) bool) bool {
	for i := range s {
		if !f(s[i]) {
			return false
		}
	}
	return true
}

func iseven(i int) bool {
	return i&1 == 0
}

func TestCount(t *testing.T) {
	xs := []int{1, 2, 3, 2, 3, 1, 3, 1, 2, 1}
	require.Equal(t, 4, Count(xs, 1))
	require.Equal(t, 4, CountFunc(xs, func(x int) bool { return x == 1 }))
}

func TestGroup(t *testing.T) {
	xs := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	ys := Group(nil, xs)
	expected := [][]int{{1}, {2, 2}, {3, 3, 3}, {4, 4, 4, 4}}
	require.Equal(t, expected, ys)
	require.True(t, Group(nil, []int{}) == nil)
}

func TestGroupFunc(t *testing.T) {
	eq := func(a, b int) bool { return iseven(a) == iseven(b) }
	xs := []int{1, 3, 5, 2, 4, 6, 7, 9, 8}
	ys := GroupFunc(nil, xs, eq)
	expected := [][]int{{1, 3, 5}, {2, 4, 6}, {7, 9}, {8}}
	require.Equal(t, expected, ys)
	require.True(t, GroupFunc(nil, []int{}, nil) == nil)
}

func TestReorder(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e"}

	for i, tt := range []struct {
		Input []int
	}{
		{Input: []int{0, 1, 2, 3, 4}},
		{Input: []int{4, 3, 2, 1, 0}},
		{Input: []int{4, 3, 2, 0, 1}},
		{Input: []int{1, 2, 3, 4, 0}},
		{Input: []int{4, 0, 1, 2, 3}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := slices.Clone(a)
			Reorder(b, tt.Input)
			for i, x := range b {
				require.Equal(t, a[tt.Input[i]], x, b, x)
			}
		})
	}

	t.Run("panic", func(t *testing.T) {
		var panicked bool
		func() {
			defer func() { panicked = recover() != nil }()
			Reorder(a, []int{0, 1, 2, 3})
		}()
		require.True(t, panicked)
	})
}

func BenchmarkReorder(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	elems := make([]int, b.N)
	perms := make([][]int, b.N)
	for i := range elems {
		elems[i] = i
	}
	for i := range perms {
		perms[i] = rnd.Perm(len(elems))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Reorder(elems, perms[i])
	}
}
