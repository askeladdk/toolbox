package require

import (
	"reflect"
	"testing"
)

func Equal[T any](tb testing.TB, expected, got T, args ...any) {
	tb.Helper()
	if !reflect.DeepEqual(expected, got) {
		if len(args) == 0 {
			tb.Fatal(expected, "is not", got)
		}
		tb.Fatal(args...)
	}
}

func True(tb testing.TB, test bool, args ...any) {
	tb.Helper()
	if !test {
		if len(args) == 0 {
			tb.Fatal("is false")
		}
		tb.Fatal(args...)
	}
}

func NoError(tb testing.TB, err error, args ...any) {
	tb.Helper()
	if err != nil {
		if len(args) == 0 {
			tb.Fatal("unexpected error:", err)
		}
		tb.Fatal(args...)
	}
}
