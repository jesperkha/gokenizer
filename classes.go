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
	"lbrace": func(iter *stringiter.StringIter) Token {
		if iter.Peek() == '{' {
			return Token{
				Lexeme:  iter.Consume(),
				matched: true,
			}
		}

		return Token{matched: false}
	},

	"rbrace": func(iter *stringiter.StringIter) Token {
		if iter.Peek() == '}' {
			return Token{
				Lexeme:  iter.Consume(),
				matched: true,
			}
		}

		return Token{matched: false}
	},

	"word": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("word", func(b byte) bool {
			return isLetter(b)
		})(iter)
	},

	"number": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("word", func(b byte) bool {
			return isNumber(b)
		})(iter)
	},

	"symbol": func(iter *stringiter.StringIter) Token {
		if isSymbol(iter.Peek()) {
			return Token{
				Lexeme:  iter.Consume(),
				matched: true,
			}
		}

		return Token{matched: false}
	},

	"line": func(iter *stringiter.StringIter) Token {
		if iter.Seek('\n') {
			line := iter.Consume()
			iter.Consume() // Consume newline to prevent infinite loop

			return Token{
				Lexeme:  line,
				matched: true,
			}
		}

		return Token{matched: false}
	},
}
