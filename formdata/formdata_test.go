package formdata

import (
	"encoding/base64"
	"errors"
	"math"
	"net/url"
	"strconv"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestDefault(t *testing.T) {
	var testform struct {
		X int
	}
	testform.X = 42
	data := url.Values{}
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 42, testform.X)
}

func TestDefaultEmpty(t *testing.T) {
	var testform struct {
		X int
	}
	testform.X = 42
	data := url.Values{}
	data["X"] = nil
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 0, testform.X)
}

func TestUnexported(t *testing.T) {
	var testform struct {
		x int
	}
	data := url.Values{}
	data.Set("x", "42")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 0, testform.x)
}

func TestSkip(t *testing.T) {
	var testform struct {
		X int `formdata:"-"`
	}
	data := url.Values{}
	data.Set("X", "42")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 0, testform.X)
}

func TestError(t *testing.T) {
	var testform struct {
		X int `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "hello")
	err := Unmarshal(data, &testform)
	var errs Errors
	require.True(t, errors.As(err, &errs))
	require.Equal(t, 1, len(errs))
	require.Equal(t, "x", errs[0].Field)
}

func TestErrorRequired(t *testing.T) {
	var testform struct {
		X int `formdata:"x,required"`
	}
	data := url.Values{}
	err := Unmarshal(data, &testform)
	var errs Errors
	require.True(t, errors.As(err, &errs))
	require.Equal(t, 1, len(errs))
	require.Equal(t, "x", errs[0].Field)
}

func TestErrorUnsupportedType(t *testing.T) {
	var testform struct {
		X complex128 `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "1+2i")
	err := Unmarshal(data, &testform)
	require.True(t, err != nil)
	require.Equal(t, `error in field "x": cannot unmarshal into unsupported type "complex128"`, err.Error())
}

func TestErrorNilPointer(t *testing.T) {
	data := url.Values{}
	err := Unmarshal(data, nil)
	require.Equal(t, "formdata: expected pointer to struct", err.Error())
}

func TestErrorNotAPointerToStruct(t *testing.T) {
	var x []string
	data := url.Values{}
	err := Unmarshal(data, &x)
	require.Equal(t, "formdata: expected pointer to struct", err.Error())
}

func TestInt(t *testing.T) {
	var testform struct {
		X int `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "-42")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, -42, testform.X)
}

func TestUint(t *testing.T) {
	var testform struct {
		X uint `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "18446744073709551615")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, math.MaxUint64, testform.X)
}

func TestFloat(t *testing.T) {
	var testform struct {
		X float64 `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "3.1415")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 3.1415, testform.X)
}

func TestString(t *testing.T) {
	var testform struct {
		X string `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "hello")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, "hello", testform.X)
}

func TestBool(t *testing.T) {
	var testform struct {
		X bool `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "on")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, true, testform.X)
}

func TestByteSlice(t *testing.T) {
	var testform struct {
		X []byte `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", base64.StdEncoding.EncodeToString([]byte("hello")))
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, []byte("hello"), testform.X)
}

func TestPointer(t *testing.T) {
	var testform struct {
		X *int `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "42")
	require.NoError(t, Unmarshal(data, &testform))
	require.True(t, testform.X != nil)
	require.Equal(t, 42, *testform.X)
}

func TestStringSlice(t *testing.T) {
	var testform struct {
		X []string `formdata:"x"`
	}
	data := url.Values{}
	data["x"] = []string{"a", "b"}
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, []string{"a", "b"}, testform.X)
}

func TestIntSlice(t *testing.T) {
	var testform struct {
		X []int `formdata:"x"`
	}
	data := url.Values{}
	data["x"] = []string{"1", "2"}
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, []int{1, 2}, testform.X)
}

func TestByteSliceSlice(t *testing.T) {
	var testform struct {
		X [][]byte `formdata:"x"`
	}
	data := url.Values{}
	data["x"] = []string{
		base64.StdEncoding.EncodeToString([]byte("hello")),
		base64.StdEncoding.EncodeToString([]byte("world")),
	}
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, [][]byte{[]byte("hello"), []byte("world")}, testform.X)
}

type textUnmarshaler struct {
	Value int
}

func (t *textUnmarshaler) UnmarshalText(x []byte) (err error) {
	t.Value, err = strconv.Atoi(string(x))
	return err
}

func TestTextUnmarshaler(t *testing.T) {
	var testform struct {
		X textUnmarshaler `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "1337")
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 1337, testform.X.Value)
}

func TestTextUnmarshalerError(t *testing.T) {
	var testform struct {
		X textUnmarshaler `formdata:"x"`
	}
	data := url.Values{}
	data.Set("x", "hello")
	err := Unmarshal(data, &testform)
	require.Equal(t, `error in field "x": strconv.Atoi: parsing "hello": invalid syntax`, err.Error())
}

func TestTextUnmarshalerSlice(t *testing.T) {
	var testform struct {
		X []textUnmarshaler `formdata:"x"`
	}
	data := url.Values{}
	data["x"] = []string{"1", "2"}
	require.NoError(t, Unmarshal(data, &testform))
	require.Equal(t, 2, len(testform.X))
	require.Equal(t, 1, testform.X[0].Value)
	require.Equal(t, 2, testform.X[1].Value)
}

func TestErrorsError(t *testing.T) {
	errs := Errors{
		{Field: "1", Err: errors.New("a")},
		{Field: "2", Err: errors.New("b")},
	}
	require.Equal(t, "", errs[:0].Error())
	require.Equal(t, `error in field "1": a`, errs[:1].Error())
	require.Equal(t, "error in field \"1\": a\nerror in field \"2\": b", errs[:2].Error())
}
