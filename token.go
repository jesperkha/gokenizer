package gokenizer

type Token struct {
	Pos    int    // Column of first character in token
	Line   int    // Line number of token
	Length int    // Length of token lexeme
	Lexeme string // Token lexeme
	Source string // The line the token is on

	// Maps the class name to the parsed value, in order
	values map[string][]string
}

// Get returns the first instance of what the specified class parsed. If
// there are multiple uses of the class see GetAt().
func (t *Token) Get(className string) string {
	return t.GetAt(className, 0)
}

// GetAt returns the n'th parsed string for the given class.
func (t *Token) GetAt(className string, index int) string {
	if l, ok := t.values[className]; ok && index < len(l) {
		return l[index]
	}

	return ""
}
