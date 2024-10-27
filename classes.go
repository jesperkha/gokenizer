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

func isBase64(c byte) bool {
	s := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	return strings.Contains(s, string(c))
}

func isHex(c byte) bool {
	s := "ABCDEFabcdef0123456789#"
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

	"base64": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("base64", func(b byte) bool {
			return isBase64(b)
		})(iter)
	},

	"hex": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("hex", func(b byte) bool {
			return isHex(b)
		})(iter)
	},

	"number": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("word", func(b byte) bool {
			return isNumber(b)
		})(iter)
	},

	"float": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("float", func(b byte) bool {
			return isNumber(b) || b == '.'
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

	"char": func(iter *stringiter.StringIter) Token {
		if isLetter(iter.Peek()) {
			return Token{
				Lexeme:  iter.Consume(),
				matched: true,
			}
		}

		return Token{matched: false}
	},
}
