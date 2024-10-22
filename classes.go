package gokenizer

import (
	"strings"

	"github.com/jesperkha/gokenizer/stringiter"
)

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isSymbol(c byte) bool {
	s := "!\"#$&%'()*+,-./:;<=>?@[]\\^_`{}|~¤§£"
	return strings.Contains(s, string(c))
}

var classes = map[string]matcherFunc{
	"lbrace": func(iter *stringiter.StringIter) matchResult {
		if iter.Peek() == '{' {
			return matchResult{
				matched: true,
				lexeme:  iter.Consume(),
			}
		}

		return matchResult{matched: false}
	},

	"rbrace": func(iter *stringiter.StringIter) matchResult {
		if iter.Peek() == '}' {
			return matchResult{
				matched: true,
				lexeme:  iter.Consume(),
			}
		}

		return matchResult{matched: false}
	},

	"word": func(iter *stringiter.StringIter) matchResult {
		word := ""

		for isLetter(iter.Peek()) {
			word += iter.Consume()
		}

		return matchResult{
			lexeme:  word,
			matched: len(word) > 0,
		}
	},

	"number": func(iter *stringiter.StringIter) matchResult {
		number := ""

		for isNumber(iter.Peek()) {
			number += iter.Consume()
		}

		return matchResult{
			lexeme:  number,
			matched: len(number) > 0,
		}
	},

	"symbol": func(iter *stringiter.StringIter) matchResult {
		if isSymbol(iter.Peek()) {
			return matchResult{
				matched: true,
				lexeme:  iter.Consume(),
			}
		}

		return matchResult{matched: false}
	},
}
