package densebits

import (
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestDenseSetTestFlip(t *testing.T) {
	a := New(1024)

	for i := 0; i < a.Len(); i++ {
		a.SetBit(i, true)
		require.True(t, a.TestBit(i))
	}

	require.Equal(t, a.Len(), a.OnesCount())

	for i := 0; i < a.Len(); i += 2 {
		a.SetBit(i, false)
		require.True(t, !a.TestBit(i))
	}

	require.Equal(t, a.Len()/2, a.OnesCount())

	a.Reset()
	require.Equal(t, 0, a.OnesCount())
	a.FlipBit(0)
	require.Equal(t, 1, a.OnesCount())

	require.True(t, !a.Equal(Set{}))
}

func TestDenseSetOps(t *testing.T) {
	a := New(128)
	b := New(128)
	c := New(128)

	a.Fill(0x5555555555555555)
	b.Fill(0xaaaaaaaaaaaaaaaa)
	require.Equal(t, 64, a.OnesCount())
	require.Equal(t, 64, b.OnesCount())

	c.Or(a, b)
	require.Equal(t, 128, c.OnesCount())

	c.And(a, b)
	require.Equal(t, 0, c.OnesCount())

	c.AndNot(a, b)
	require.Equal(t, 64, c.OnesCount())

	c.Xor(a, b)
	require.Equal(t, 128, c.OnesCount())

	c.Not(a)
	require.True(t, c.Equal(b))
	require.True(t, !c.Equal(a))

	d := New(256)
	d.Xor(c, d)
	require.Equal(t, 128, d.Len())
}

func TestDenseShift(t *testing.T) {
	a := Set{0xff00ff00ff00ff00, 0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
	b := Set{0xfe01fe01fe01fe01, 0xfe01fe01fe01fe01, 0xfe01fe01fe01fe00}
	c := Set{0x7f00ff00ff00ff00, 0xff00ff00ff00ff00, 0xff00ff00ff00ff00}

	require.Equal(t, uint64(1), a.ShiftLeft(a, 1))
	require.Equal(t, b, a)
	require.Equal(t, uint64(0), a.ShiftRight(a, 1))
	require.Equal(t, c, a)

	d := Set{0xff00ff00ff00ff00, 0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
	e := Set{0x00ff00ff00ff00ff, 0x00ff00ff00ff00ff, 0x00ff00ff00ff00ff}
	require.Equal(t, uint64(0), d.ShiftRight(d, 8))
	require.Equal(t, e, d)
}

func TestDenseRotate(t *testing.T) {
	a := Set{0xff00ff00ff00ff00, 0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
	b := Set{0xfe01fe01fe01fe01, 0xfe01fe01fe01fe01, 0xfe01fe01fe01fe01}
	c := New(a.Len())

	c.RotateLeft(a, 1)
	require.True(t, c.Equal(b))
	c.RotateRight(c, 1)
	require.True(t, c.Equal(a))
}

func TestDenseSlice(t *testing.T) {
	a := New(256)
	a.Slice(0, 128).Fill(0x5555555555555555)
	a.Slice(128, 256).Fill(0xaaaaaaaaaaaaaaaa)

	e := Set{
		0x5555555555555555, 0x5555555555555555,
		0xaaaaaaaaaaaaaaaa, 0xaaaaaaaaaaaaaaaa,
	}

	require.Equal(t, e, a)
}

func TestDenseAccomodate(t *testing.T) {
	a := New(128)
	a.Fill(0x5555555555555555)
	b := New(256)
	b.Fill(0xaaaaaaaaaaaaaaaa)
	var c Set
	require.Equal(t, 0, c.Len())
	c.Or(a, b)
	require.Equal(t, 128, c.Len())
}

func TestDenseString(t *testing.T) {
	a := Set{0x5555555555555555, 0xaaaaaaaaaaaaaaaa, 0x5555555555555555}
	require.True(t, a.String() != "[]")
	require.Equal(t, "[]", a.Slice(0, 0).String())
}
