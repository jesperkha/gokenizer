package gokenizer

import (
	"fmt"

	"github.com/jesperkha/gokenizer/stringiter"
)

type MatchFunc func(Token) error

type Tokenizer struct {
	err        error
	pos        int
	matchFuncs []matcherFunc
	callbacks  []MatchFunc
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

// Pattern adds a new pattern to the tokenizer. If a match is found, the
// callback function f is called. The callback may return an error which
// will be properly formatted and returned by Run().
func (t *Tokenizer) Pattern(pattern string, f MatchFunc) {
	mf, err := createMatcherFunc(pattern)
	if err != nil {
		t.err = err
	}
	t.matchFuncs = append(t.matchFuncs, mf)
	t.callbacks = append(t.callbacks, f)
}

// Runs tokenizer on given input string. Returns first error received by a
// pattern callback function. Patterns are matched by the order the are
// defined in.
func (t *Tokenizer) Run(s string) error {
	if t.err != nil {
		return t.err
	}

	iter := stringiter.New(s)

	for !iter.Eof() {
		result := matchResult{matched: false}

		for idx, mf := range t.matchFuncs {
			iter.Push()
			result = mf(&iter)

			if result.matched {
				token := Token{
					Lexeme: result.lexeme,
				}

				if err := t.callbacks[idx](token); err != nil {
					return err
				}

				break
			}

			iter.Pop()
		}

		if !result.matched {
			c := iter.Consume() // Next
			fmt.Printf("unmatched '%s'\n", c)
		}
	}

	return nil
}
