package sparsebits

import (
	"math/rand"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestSparseSetGetFlip(t *testing.T) {
	n := 32000
	b := New(n)

	require.True(t, b.Len()&(b.Len()-1) == 0)

	require.True(t, !b.Get(0))
	b.Set(0, false)

	rnd := rand.New(rand.NewSource(0))
	xs := rnd.Perm(b.Len())

	for _, x := range xs {
		b.Set(x, true)
		require.True(t, b.Get(x))
	}

	for _, x := range xs {
		x = xs[x]
		b.Flip(x)
		require.True(t, !b.Get(x))
	}

	require.Equal(t, 0, b.OnesCount())

	require.True(t, !b.Get(b.Len()))
}

func TestSparseCount(t *testing.T) {
	b := New(1000)

	rnd := rand.New(rand.NewSource(0))
	xs := rnd.Perm(b.Len())

	for _, x := range xs {
		b.Set(x, true)
	}

	for i := 0; i < len(xs); i += 2 {
		b.Set(i, false)
	}

	require.Equal(t, b.Len()/2, b.OnesCount())

	b.Reset()
	require.Equal(t, 0, b.OnesCount())

	b.Set(65, true)
	require.Equal(t, 1, b.OnesCount())
}
