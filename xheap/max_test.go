package xheap

import (
	"math"
	"math/rand"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestMaxPushPop(t *testing.T) {
	const count = 100
	var h Max[int, int]

	rnd := rand.New(rand.NewSource(0))
	for i, prio := range rnd.Perm(count) {
		h.Push(i, prio)
	}

	require.Equal(t, count, h.Len())

	prio := math.MaxInt
	for !h.Empty() {
		_, peekprio := h.Peek()
		require.True(t, peekprio < prio)
		_, prio = h.Pop()
	}
}

func TestMaxFix(t *testing.T) {
	var h Max[int, int]
	h.Push(1, 1)
	h.Push(2, 3)
	h.Push(3, 2)
	top, _ := h.Peek()
	require.Equal(t, 2, top)
	h[0].Priority = 0
	h.Fix(0)
	top, _ = h.Peek()
	require.Equal(t, 3, top)
}

func TestMaxInit(t *testing.T) {
	h := Max[int, int]{
		{1, 1},
		{2, 2},
		{3, 3},
	}

	h.Init()
	require.Equal(t, 3, h.Len())
	require.Equal(t, 3, h[0].Value)
	require.Equal(t, 2, h[1].Value)
	require.Equal(t, 1, h[2].Value)

	h.Remove(1)
	require.Equal(t, 2, h.Len())
	require.Equal(t, 3, h[0].Value)
	require.Equal(t, 1, h[1].Value)

	h.Reset()
	require.Equal(t, 0, h.Len())
}
