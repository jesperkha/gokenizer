package stringiter

import (
	"strings"
)

type StringIter struct {
	s        string
	pos      int
	peekPos  int
	posStack []int
}

func New(s string) StringIter {
	return StringIter{
		s:       s,
		pos:     0,
		peekPos: 0,
	}
}

func (iter *StringIter) Eof() bool {
	return iter.pos >= len(iter.s)
}

// Saves pos
func (iter *StringIter) Push() {
	iter.posStack = append(iter.posStack, iter.pos)
}

// Restores to previous saved pos. Returns difference of prev pos and pos.
func (iter *StringIter) Pop() int {
	prev := iter.pos
	iter.pos = iter.posStack[len(iter.posStack)-1]
	iter.posStack = iter.posStack[:len(iter.posStack)-1]
	return prev - iter.pos
}

// Consumes and returns all characters from current pos to peek pos.
// Returns empty string on eof.
func (iter *StringIter) Consume() string {
	if iter.Eof() {
		return ""
	}

	end := iter.peekPos
	if iter.peekPos >= len(iter.s) {
		end = len(iter.s)
	}

	if iter.pos == iter.peekPos {
		end++
	}

	s := iter.s[iter.pos:end]
	iter.pos = end
	iter.peekPos = iter.pos
	return s
}

// Peeks next character and moves peek pointer. Returns null byte on eof.
func (iter *StringIter) Peek() byte {
	if iter.peekPos >= len(iter.s) {
		return 0
	}

	b := iter.s[iter.peekPos]
	iter.peekPos++
	return b
}

// Moves peek pointer by n
func (iter *StringIter) PeekN(n uint) {
	iter.peekPos += int(n)
}

// Moves peek pointer to c. Returns false if c is not found.
func (iter *StringIter) Seek(c byte) bool {
	if i := strings.IndexByte(iter.Remainder(), c); i != -1 {
		iter.peekPos = i
		return true
	}

	return false
}

// Returns remaining string from iter pos. Empty string on eof. Does not move pos.
func (iter *StringIter) Remainder() string {
	if iter.Eof() {
		return ""
	}

	return iter.s[iter.pos:]
}

// Restores peek pos to pos.
func (iter *StringIter) Restore() {
	iter.peekPos = iter.pos
}

// Resets iterator to beginning
func (iter *StringIter) Reset() {
	iter.pos = 0
	iter.peekPos = 0
}
