package xslices

import (
	"cmp"
	"math/rand"
	"slices"
	"strconv"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestArgSort(t *testing.T) {
	funcs := []func(a, x []int){
		ArgSort[[]int, int],
		func(a, x []int) { ArgSortFunc(a, x, cmp.Less[int]) },
		func(a, x []int) { ArgSortStableFunc(a, x, cmp.Less[int]) },
	}
	for _, fn := range funcs {
		data := []int{3, 8, 0, 2, 9, 1, 5, 7, 6, 4}
		x := make([]int, len(data))
		fn(data, x)
		Reorder(data, x)
		require.True(t, slices.IsSorted(data), "not sorted")
	}
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

func TestSelect(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	xs := r.Perm(1337)

	for i := 0; i < len(xs); i++ {
		{
			ys := slices.Clone(xs)
			require.Equal(t, i, Select(ys, i))
			require.Equal(t, i, ys[i])
		}
		{
			ys := slices.Clone(xs)
			require.Equal(t, i, SelectFunc(ys, i, cmp.Compare))
			require.Equal(t, i, ys[i])
		}
	}
}

func TestMedian3(t *testing.T) {
	require.Equal(t, 0, median3(0, 0, 0))
	require.Equal(t, 0, median3(0, 0, 1))
	require.Equal(t, 0, median3(0, 0, 2))
	require.Equal(t, 0, median3(0, 1, 0))
	require.Equal(t, 1, median3(0, 1, 1))
	require.Equal(t, 1, median3(0, 1, 2))
	require.Equal(t, 0, median3(0, 2, 0))
	require.Equal(t, 1, median3(0, 2, 1))
	require.Equal(t, 2, median3(0, 2, 2))
	require.Equal(t, 0, median3(1, 0, 0))
	require.Equal(t, 1, median3(1, 0, 1))
	require.Equal(t, 1, median3(1, 0, 2))
	require.Equal(t, 1, median3(1, 1, 0))
	require.Equal(t, 1, median3(1, 1, 1))
	require.Equal(t, 1, median3(1, 1, 2))
	require.Equal(t, 1, median3(1, 2, 0))
	require.Equal(t, 1, median3(1, 2, 1))
	require.Equal(t, 2, median3(1, 2, 2))
	require.Equal(t, 0, median3(2, 0, 0))
	require.Equal(t, 1, median3(2, 0, 1))
	require.Equal(t, 2, median3(2, 0, 2))
	require.Equal(t, 1, median3(2, 1, 0))
	require.Equal(t, 1, median3(2, 1, 1))
	require.Equal(t, 2, median3(2, 1, 2))
	require.Equal(t, 2, median3(2, 2, 0))
	require.Equal(t, 2, median3(2, 2, 1))
	require.Equal(t, 2, median3(2, 2, 2))
}

func TestMedian3Func(t *testing.T) {
	require.Equal(t, 0, median3Func(0, 0, 0, cmp.Compare))
	require.Equal(t, 0, median3Func(0, 0, 1, cmp.Compare))
	require.Equal(t, 0, median3Func(0, 0, 2, cmp.Compare))
	require.Equal(t, 0, median3Func(0, 1, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(0, 1, 1, cmp.Compare))
	require.Equal(t, 1, median3Func(0, 1, 2, cmp.Compare))
	require.Equal(t, 0, median3Func(0, 2, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(0, 2, 1, cmp.Compare))
	require.Equal(t, 2, median3Func(0, 2, 2, cmp.Compare))
	require.Equal(t, 0, median3Func(1, 0, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 0, 1, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 0, 2, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 1, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 1, 1, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 1, 2, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 2, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(1, 2, 1, cmp.Compare))
	require.Equal(t, 2, median3Func(1, 2, 2, cmp.Compare))
	require.Equal(t, 0, median3Func(2, 0, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(2, 0, 1, cmp.Compare))
	require.Equal(t, 2, median3Func(2, 0, 2, cmp.Compare))
	require.Equal(t, 1, median3Func(2, 1, 0, cmp.Compare))
	require.Equal(t, 1, median3Func(2, 1, 1, cmp.Compare))
	require.Equal(t, 2, median3Func(2, 1, 2, cmp.Compare))
	require.Equal(t, 2, median3Func(2, 2, 0, cmp.Compare))
	require.Equal(t, 2, median3Func(2, 2, 1, cmp.Compare))
	require.Equal(t, 2, median3Func(2, 2, 2, cmp.Compare))
}

func BenchmarkArgSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := rand.New(rand.NewSource(0))
		xs := r.Perm(10000)
		ys := make([]int, len(xs))
		b.StartTimer()
		ArgSort(xs, ys)
	}
}

func BenchmarkArgSortFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := rand.New(rand.NewSource(0))
		xs := r.Perm(10000)
		ys := make([]int, len(xs))
		b.StartTimer()
		ArgSortFunc(xs, ys, cmp.Less)
	}
}

func BenchmarkSelect(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	xs := r.Perm(10000)
	ys := make([]int, len(xs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(ys, xs)
		j := r.Intn(len(ys))
		b.StartTimer()
		Select(ys, j)
	}
}

func BenchmarkSelectFunc(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	xs := r.Perm(10000)
	ys := make([]int, len(xs))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(ys, xs)
		j := r.Intn(len(ys))
		b.StartTimer()
		SelectFunc(ys, j, cmp.Compare)
	}
}
