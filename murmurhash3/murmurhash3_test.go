package murmurhash3

import (
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"strconv"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestMurmurHash(t *testing.T) {
	for _, tt := range []struct {
		seed uint64
		h0   uint64
		h1   uint64
		s    string
	}{
		{0x00, 0x0000000000000000, 0x0000000000000000, ""},
		{0x00, 0xcbd8a7b341bd9b02, 0x5b1e906a48ae1d19, "hello"},
		{0x00, 0x342fac623a5ebc8e, 0x4cdcbc079642414d, "hello, world"},
		{0x00, 0xb89e5988b737affc, 0x664fc2950231b2cb, "19 Jan 2038 at 3:14:07 AM"},
		{0x00, 0xcd99481f9ee902c9, 0x695da1a38987b6e7, "The quick brown fox jumps over the lazy dog."},

		{0x01, 0x4610abe56eff5cb5, 0x51622daa78f83583, ""},
		{0x01, 0xa78ddff5adae8d10, 0x128900ef20900135, "hello"},
		{0x01, 0x8b95f808840725c6, 0x1597ed5422bd493b, "hello, world"},
		{0x01, 0x2a929de9c8f97b2f, 0x56a41d99af43a2db, "19 Jan 2038 at 3:14:07 AM"},
		{0x01, 0xfb3325171f9744da, 0xaaf8b92a5f722952, "The quick brown fox jumps over the lazy dog."},

		{0x2a, 0xf02aa77dfa1b8523, 0xd1016610da11cbb9, ""},
		{0x2a, 0xc4b8b3c960af6f08, 0x2334b875b0efbc7a, "hello"},
		{0x2a, 0xb91864d797caa956, 0xd5d139a55afe6150, "hello, world"},
		{0x2a, 0xfd8f19ebdc8c6b6a, 0xd30fdc310fa08ff9, "19 Jan 2038 at 3:14:07 AM"},
		{0x2a, 0x74f33c659cda5af7, 0x4ec7a891caf316f0, "The quick brown fox jumps over the lazy dog."},
	} {
		t.Run(tt.s, func(t *testing.T) {
			var expected [Size]byte
			binary.BigEndian.PutUint64(expected[0:], tt.h0)
			binary.BigEndian.PutUint64(expected[8:], tt.h1)

			h128 := NewWithSeed(tt.seed, tt.seed)
			h64 := New64WithSeed(tt.seed, tt.seed)
			_, _ = h128.Write([]byte(tt.s))
			require.Equal(t, expected[:], h128.Sum(nil))
			require.Equal(t, tt.h0, tt.h0, h64.Sum64())
		})
	}
}

func TestWrite(t *testing.T) {
	d128 := New()
	d64 := New64()
	const str = "The quick brown fox jumps over the lazy dog."
	_, _ = io.WriteString(d128, str[:3])
	_, _ = io.WriteString(d128, str[3:9])
	_, _ = io.WriteString(d128, str[9:])
	_, _ = io.WriteString(d64, str[:3])
	_, _ = io.WriteString(d64, str[3:9])
	_, _ = io.WriteString(d64, str[9:])

	h128 := d128.Sum(nil)
	h64 := d64.Sum(nil)
	expectedh := []byte{
		0xcd, 0x99, 0x48, 0x1f, 0x9e, 0xe9, 0x02, 0xc9,
		0x69, 0x5d, 0xa1, 0xa3, 0x89, 0x87, 0xb6, 0xe7,
	}
	require.Equal(t, expectedh, h128)
	require.Equal(t, expectedh[:8], h64)
}

func TestReset(t *testing.T) {
	d := NewWithSeed(1, 2).(*digest128)
	_, _ = d.Write([]byte("blah"))
	d.Reset()
	require.Equal(t, d.s0, 1)
	require.Equal(t, d.s1, 2)
	require.Equal(t, d.h0, 1)
	require.Equal(t, d.h1, 2)
	require.Equal(t, d.n, 0)
	require.Equal(t, d.head, 0)
	require.Equal(t, d.tail, [16]byte{})
}

func TestObvious(t *testing.T) {
	var d128 digest128
	var d64 digest64
	require.Equal(t, d128.BlockSize(), 1)
	require.Equal(t, d64.BlockSize(), 1)
	require.Equal(t, d128.Size(), 16)
	require.Equal(t, d64.Size(), 8)
}

func TestSum(t *testing.T) {
	s := "0123456789ABCDEF0"

	require.Equal(t, 17, len(s))

	for i := range s {
		t.Run(s[:i], func(t *testing.T) {
			h128 := Sum([]byte(s[:i]))
			h64 := Sum64([]byte(s[:i]))
			require.Equal(t, binary.BigEndian.Uint64(h128[:]), h64)
		})
	}
}

func TestBinaryEncoding(t *testing.T) {
	var d digest
	_, _ = d.Write([]byte("hello world"))
	bin, _ := d.MarshalBinary()
	d2 := d
	d.Reset()
	_ = d.UnmarshalBinary(bin)
	require.Equal(t, d, d2)

	require.Equal(t, errors.New("murmurhash3: invalid hash state size"), d.UnmarshalBinary(nil))
}

func BenchmarkSum(b *testing.B) {
	buf := make([]byte, 8192)
	rnd := rand.New(rand.NewSource(1))
	rnd.Read(buf)
	for length := 32; length <= cap(buf); length *= 2 {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			buf = buf[:length]
			b.SetBytes(int64(length))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Sum(buf)
			}
		})
	}
}
