package sparse

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestSet(t *testing.T) {
	s := NewSet[uint](0)

	rnd := rand.New(rand.NewSource(0))
	perms := rnd.Perm(100)

	for i := 0; i < len(perms); i += 2 {
		s.Add(uint(perms[i]))
	}

	for i := 0; i < len(perms); i += 2 {
		require.True(t, s.Has(uint(perms[i])))
	}

	for i := 1; i < len(perms); i += 2 {
		require.True(t, !s.Has(uint(perms[i])))
	}

	for i, m := range s.Keys() {
		require.Equal(t, m, uint(perms[i*2]))
	}

	require.Equal(t, 50, s.Len())

	for i := 0; i < len(perms); i++ {
		s.Del(uint(perms[i]))
	}

	require.Equal(t, 0, s.Len())

	s.Add(1)
	s.Add(1)
	require.Equal(t, 1, s.Len())
	s.Reset()
	require.Equal(t, 0, s.Len())
}

func TestSetSort(t *testing.T) {
	s := NewSet[uint](100)

	rnd := rand.New(rand.NewSource(0))
	perms := rnd.Perm(s.Cap())

	for i := 0; i < len(perms); i += 2 {
		s.Add(uint(perms[i]))
	}

	for i := 0; i < len(perms); i += 2 {
		require.True(t, s.Has(uint(perms[i])))
	}

	sort.Sort(s)
	require.True(t, sort.IsSorted(s))

	for i := 0; i < len(perms); i += 2 {
		x := uint(perms[i])
		require.True(t, s.Has(x), x, s.sparse[x])
		s.Del(x)
		require.True(t, !s.Has(x))
	}

	require.Equal(t, 0, s.Len())
}
