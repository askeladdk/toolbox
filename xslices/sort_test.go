package xslices

import (
	"cmp"
	"math/rand"
	"slices"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

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
