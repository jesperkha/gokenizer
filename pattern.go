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
	values  map[string][]string
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

func parseClass(iter *stringiter.StringIter) (name string, err error) {
	if iter.Consume() != "{" {
		return name, fmt.Errorf("expected { before class name")
	}

	if !iter.Seek('}') {
		return name, fmt.Errorf("expected } after class name")
	}

	name = iter.Consume()
	iter.Consume() // }
	return name, err
}

// Returns a function that matches based on the given pattern.
func createMatcherFunc(pattern string) (mf matcherFunc, err error) {
	funcs := []matcherFunc{} // Both should be same length
	classNames := []string{}

	pIter := stringiter.New(pattern)

	for !pIter.Eof() {
		if pIter.Peek() == '{' {
			// Parse class name if we find a {
			pIter.Restore()
			className, err := parseClass(&pIter)
			if err != nil {
				return mf, err
			}

			f, ok := classes[className]
			if !ok {
				return mf, fmt.Errorf("unknown class name '%s'", className)
			}

			funcs = append(funcs, f)
			classNames = append(classNames, className)
		} else if pIter.Seek('{') {
			// Parse static word if there are characters before a {
			staticWord := pIter.Consume()
			if staticWord == "" {
				return mf, fmt.Errorf("parser error")
			}

			funcs = append(funcs, literalMatcherFunc(staticWord))
			classNames = append(classNames, "static")
		} else {
			// Otherwise the rest of the pattern string is a static word
			funcs = append(funcs, literalMatcherFunc(pIter.Remainder()))
			classNames = append(classNames, "static")
			break
		}
	}

	f := func(iter *stringiter.StringIter) (res matchResult) {
		iter.Push()
		values := make(map[string][]string)

		for idx, mf := range funcs {
			tempResult := mf(iter)
			if !tempResult.matched {
				return res
			}

			if className := classNames[idx]; className != "static" {
				values[className] = append(values[className], tempResult.lexeme)
			}
		}

		length := iter.Pop()
		iter.PeekN(uint(length))
		matchedString := iter.Consume()

		return matchResult{
			matched: true,
			lexeme:  matchedString,
			values:  values,
		}
	}

	return f, err
}
