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
	"any": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("any", func(b byte) bool {
			return true
		})(iter)
	},

	"ws": func(iter *stringiter.StringIter) Token {
		c := ""
		for !iter.Eof() {
			b := iter.Peek()
			if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
				c += iter.Consume()
				continue
			}
			break
		}
		return Token{matched: true, Lexeme: c}
	},

	"text": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("text", func(b byte) bool {
			return !(b == ' ' || b == '\t' || b == '\n' || b == '\r')
		})(iter)
	},

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

	"var": func(iter *stringiter.StringIter) Token {
		return checkFuncToMatchFunc("var", func(b byte) bool {
			return isLetter(b) || b == '$' || b == '_'
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
		return checkFuncToMatchFunc("number", func(b byte) bool {
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

	"string": func(iter *stringiter.StringIter) Token {
		if iter.Peek() == '"' {
			iter.Push()
			iter.Consume()

			if iter.Seek('"') {
				// Consumes string content, then terminating quote
				str := iter.Consume()
				iter.Consume()
				return Token{
					Lexeme:  str,
					matched: true,
				}
			}

			iter.Pop()
		}

		return Token{matched: false}
	},
}
