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

func TestBasic(t *testing.T) {
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

func TestSeek(t *testing.T) {
	s := "Hello, world!"
	iter := stringiter.New(s)

	if !iter.Seek(',') {
		t.Fatal("failed to seek to ','")
	}

	assertEq(t, "Hello", iter.Consume())
	assertEq(t, ",", iter.Consume())
	assertEq(t, " world!", iter.Remainder())
}

func TestStack(t *testing.T) {
	s := "Hello, world!"
	iter := stringiter.New(s)

	iter.Push()

	iter.PeekN(5)
	iter.Consume()

	iter.Push()

	iter.PeekN(99)
	iter.Consume()

	l := iter.Pop()
	if l != 8 {
		t.Fatalf("expected popped len of %d, got %d", 8, l)
	}

	iter.Pop()
	assertEq(t, s, iter.Remainder())
}
