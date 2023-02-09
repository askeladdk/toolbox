package require

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected, got any, args ...any) {
	if !reflect.DeepEqual(expected, got) {
		args = append([]any{expected, "!=", got}, args...)
		t.Fatal(args...)
	}
}

func True(t *testing.T, test bool, args ...any) {
	if !test {
		t.Fatal(args...)
	}
}
