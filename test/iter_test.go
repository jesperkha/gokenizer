package test

import (
	"testing"

	"github.com/jesperkha/gokenizer/stringiter"
)

func assertEq(t *testing.T, a, b string) {
	if a != b {
		t.Fatalf("expected '%s', got '%s'", a, b)
	}
}

func TestStringIterator(t *testing.T) {
	s := "Hello, world!"
	iter := stringiter.New(s)

	assertEq(t, "H", iter.Consume())

	iter.PeekN(4)
	assertEq(t, "ello", iter.Consume())

	comma := iter.Peek()
	assertEq(t, ",", string(comma))
	assertEq(t, ",", iter.Consume())

	assertEq(t, " world!", iter.Remainder())
}
