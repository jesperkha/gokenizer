package gokenizer

import (
	"fmt"

	"github.com/jesperkha/gokenizer/stringiter"
)

/*

"{word}={number}"  pattern
{word}             class

Classes

word        any sequence of ascii alphabetical characters
number      any sequence of ascii numerical characters
float       same as number, but accepts periods too
var         variable name, same as word in addition to underscore,
dollar, and numbers after the first character
symbol      any printable ascii character that is not a number or letter
whitespace  space, newline, and tab

*/

// For future additional metadata about the expression
type matchResult struct {
	matched bool
	lexeme  string
}

// Matches with the given string. The implementation is dynamically created in createPattern.
type matcherFunc func(iter *stringiter.StringIter) matchResult

// Returns a function that matches the string literal s.
func literalMatcherFunc(s string) matcherFunc {
	return func(iter *stringiter.StringIter) matchResult {
		iter.PeekN(uint(len(s)))
		lexeme := iter.Consume()

		return matchResult{
			matched: lexeme == s,
			lexeme:  lexeme,
		}
	}
}

// Returns a function that matches based on the given pattern.
func createMatcherFunc(pattern string) (mf matcherFunc, err error) {
	funcs := []matcherFunc{}
	pIter := stringiter.New(pattern)

	for !pIter.Eof() {
		foundClass := pIter.Seek('{')
		if foundClass {
			return mf, fmt.Errorf("classes not implemented yet!")
		}

		funcs = append(funcs, literalMatcherFunc(pIter.Remainder()))
		break
	}

	f := func(iter *stringiter.StringIter) (res matchResult) {
		iter.Push()

		for _, mf := range funcs {
			result := mf(iter)
			if !result.matched {
				return res
			}
		}

		length := iter.Pop()
		iter.PeekN(uint(length))
		matchedString := iter.Consume()

		return matchResult{
			matched: true,
			lexeme:  matchedString,
		}
	}

	return f, err
}
