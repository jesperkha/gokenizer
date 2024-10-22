package gokenizer

import (
	"github.com/jesperkha/gokenizer/stringiter"
)

type MatchFunc func(Token) error

type Tokenizer struct {
	err        error
	pos        int
	matchFuncs []matcherFunc
	callbacks  []MatchFunc
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
		if err := t.matchNext(&iter); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tokenizer) matchNext(iter *stringiter.StringIter) error {
	callbackIdx := 0
	result := matchResult{}

	for idx, mf := range t.matchFuncs {
		iter.Push()
		if result = mf(iter); result.matched {
			callbackIdx = idx
			break
		}

		iter.Pop()
	}

	if !result.matched {
		iter.Consume() // Next
		return nil
	}

	token := Token{
		// Todo: more info in Token
		Lexeme: result.lexeme,
		values: result.values,
	}

	err := t.callbacks[callbackIdx](token)
	return err
}
