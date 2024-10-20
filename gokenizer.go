package gokenizer

type MatchFunc func(Token) error

type Tokenizer struct {
	pos int
	fs  []MatchFunc
}

type Token struct {
	Pos    int    // Column of first character in token
	Line   int    // Line number of token
	Length int    // Length of token lexeme
	Lexeme string // Token lexeme
	Source string // The line the token is on
}

func New() Tokenizer {
	return Tokenizer{
		pos: 0,
	}
}

// Pattern adds a new pattern to the tokenizer. If a match is found, the callback function f is called.
// The callback may return an error which will be properly formatted and returned by Run().
func (t *Tokenizer) Pattern(pattern string, f MatchFunc) {

}

// Runs tokenizer on given input string. Returns first error received by a pattern callback function.
// Patterns are matched by the order the are defined in.
func (t *Tokenizer) Run(s string) error {

	return nil
}
