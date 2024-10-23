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
		word := ""

		for isLetter(iter.Peek()) {
			word += iter.Consume()
		}

		return Token{
			Lexeme:  word,
			matched: len(word) > 0,
		}
	},

	"number": func(iter *stringiter.StringIter) Token {
		number := ""

		for isNumber(iter.Peek()) {
			number += iter.Consume()
		}

		return Token{
			Lexeme:  number,
			matched: len(number) > 0,
		}
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
}
