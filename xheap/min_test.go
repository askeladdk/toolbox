package xheap

import (
	"math/rand"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestMinPushPop(t *testing.T) {
	const count = 100
	var h Min[int, int]

	rnd := rand.New(rand.NewSource(0))
	for i, prio := range rnd.Perm(count) {
		h.Push(i, prio)
	}

	require.Equal(t, count, h.Len())

	prio := -1
	for !h.Empty() {
		_, peekprio := h.Peek()
		require.True(t, peekprio > prio)
		_, prio = h.Pop()
	}
}

func TestMinFix(t *testing.T) {
	var h Min[int, int]
	h.Push(1, 1)
	h.Push(2, 3)
	h.Push(3, 2)
	top, _ := h.Peek()
	require.Equal(t, 1, top)
	h[0].Prio = 10
	h.Fix(0)
	top, _ = h.Peek()
	require.Equal(t, 3, top)
}

func TestMinInit(t *testing.T) {
	h := Min[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	h.Init()
	require.Equal(t, 3, h.Len())
	require.Equal(t, 1, h[0].Elem)
	require.Equal(t, 2, h[1].Elem)
	require.Equal(t, 3, h[2].Elem)

	h.Remove(1)
	require.Equal(t, 2, h.Len())
	require.Equal(t, 1, h[0].Elem)
	require.Equal(t, 3, h[1].Elem)

	h.Reset()
	require.Equal(t, 0, h.Len())
}
