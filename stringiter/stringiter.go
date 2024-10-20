package stringiter

type StringIter struct {
	s       string
	pos     int
	peekPos int
}

func New(s string) StringIter {
	return StringIter{
		s:       s,
		pos:     0,
		peekPos: 0,
	}
}

// Consumes and returns all characters from current pos to peek pos.
// Returns empty string on eof.
func (iter *StringIter) Consume() string {
	if iter.pos >= len(iter.s) {
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

// Returns remaining string from iter pos. Empty string on eof. Does not move pos.
func (iter *StringIter) Remainder() string {
	if iter.pos >= len(iter.s) {
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
