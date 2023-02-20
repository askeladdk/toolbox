package require

import (
	"reflect"
	"testing"
)

func Equal[T any](t testing.TB, expected, got T, args ...any) {
	if !reflect.DeepEqual(expected, got) {
		t.Helper()
		if len(args) == 0 {
			t.Fatal(expected, "is not", got)
		}
		t.Fatal(args...)
	}
}

func True(t testing.TB, test bool, args ...any) {
	if !test {
		t.Helper()
		if len(args) == 0 {
			t.Fatal("is false")
		}
		t.Fatal(args...)
	}
}
