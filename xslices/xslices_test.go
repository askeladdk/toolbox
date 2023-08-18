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

func TestPermute(t *testing.T) {
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
			Permute(b, tt.Input)
			for i, x := range b {
				require.Equal(t, a[tt.Input[i]], x, b, x)
			}
		})
	}

	t.Run("panic", func(t *testing.T) {
		var panicked bool
		func() {
			defer func() { panicked = recover() != nil }()
			Permute(a, []int{0, 1, 2, 3})
		}()
		require.True(t, panicked)
	})
}

func TestPartitionFunc(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))

	for i, tt := range []struct {
		Input []int
	}{
		{
			Input: []int{},
		},
		{
			Input: []int{1},
		},
		{
			Input: []int{2},
		},
		{
			Input: []int{2, 4, 8, 6, 3, 9, 1, 7, 5},
		},
		{
			Input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			Input: []int{7, 2, 9, 4, 5, 1, 6, 4, 3},
		},
		{
			Input: []int{2, 4, 8, 6},
		},
		{
			Input: []int{3, 9, 1, 7, 5},
		},
		{
			Input: rnd.Perm(1337),
		},
		{
			Input: rnd.Perm(42069),
		},
		{
			Input: rnd.Perm(100000),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			i := PartitionFunc(tt.Input, iseven)
			require.True(t, all(tt.Input[:i], iseven))
			require.True(t, !slices.ContainsFunc(tt.Input[i:], iseven))
		})
	}
}

func BenchmarkPartition(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	s := rnd.Perm(100000)
	b.ResetTimer()
	half := func(n int) bool { return n < 50000 }

	b.Run("PartitionFunc(iseven)", func(b *testing.B) {
		z := make([]int, len(s))
		b.SetBytes(int64(len(s) * 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			copy(z, s)
			PartitionFunc(z, iseven)
		}
	})

	b.Run("PartitionFunc(half)", func(b *testing.B) {
		z := make([]int, len(s))
		b.SetBytes(int64(len(s) * 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			copy(z, s)
			PartitionFunc(z, half)
		}
	})
}

func BenchmarkPermute(b *testing.B) {
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
		Permute(elems, perms[i])
	}
}
